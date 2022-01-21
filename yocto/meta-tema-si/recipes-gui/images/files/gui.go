package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	tcell "github.com/gdamore/tcell/v2"
	encoding "github.com/gdamore/tcell/v2/encoding"
	"github.com/mattn/go-runewidth"
)

type GuiObject struct {
	Type string `json:"type,omitempty"`
	X    int    `json:"x,omitempty"`
	Y    int    `json:"y,omitempty"`
}

var (
	connected      bool
	connectedSync  sync.Mutex
	tcpConn        net.Conn
	fish           []GuiObject
	submarine      GuiObject
	artifact       GuiObject
	fishStyle      tcell.Style
	artifactStyle  tcell.Style
	submarineStyle tcell.Style
)

var defStyle tcell.Style
var submarineDesign = []string{
	"          __|___",
	"         /      \\",
	" _______/    O   \\_______",
	"<                        \\_____  I",
	" \\   O      O     O            >-=",
	"  \\___________________________/  I"}

const (
	fishDesign      = "><(((*>"
	artifactDesign  = "[*]"
	submarineHeight = 6
	submarineLength = 35
	ncursesTTY      = "/dev/ttyAMA0"
)

func initTcpClient() {
	for {
		connectedSync.Lock()
		alreadyConnected := connected
		connectedSync.Unlock()
		if !alreadyConnected {
			conn, err := net.Dial("tcp", "127.0.0.1:8000")
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				time.Sleep(time.Duration(5) * time.Second)
				continue
			}
			tcpConn = conn
			fmt.Fprintf(os.Stderr, "%v : connected\n", conn.RemoteAddr().String())
			connectedSync.Lock()
			connected = true
			connectedSync.Unlock()
			go receiveData(tcpConn)
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func receiveData(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v : disconnected\n", conn.RemoteAddr().String())
			conn.Close()
			connectedSync.Lock()
			connected = false
			connectedSync.Unlock()
			fmt.Fprintf(os.Stderr, "%v : : end receiving data\n", conn.RemoteAddr().String())
			return
		}

		var newObj GuiObject
		err = json.Unmarshal([]byte(message), &newObj)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}

		switch {
		case newObj.Type == "submarine":
			submarine.X = newObj.X
			submarine.Y = newObj.Y

		case newObj.Type == "artifact":
			artifact = newObj
		case newObj.Type == "fish":
			fish = append(fish, newObj)
		}
	}
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func drawFish(s tcell.Screen, fish []GuiObject) {

	for _, f := range fish {
		emitStr(s, f.X, f.Y, fishStyle, fishDesign)
	}
}

func drawArtifact(s tcell.Screen, artifact GuiObject) {

	emitStr(s, artifact.X, artifact.Y, artifactStyle, artifactDesign)
}

func drawSubmarine(s tcell.Screen, submarine GuiObject) {
	screenLine := submarine.Y

	for _, submarineLine := range submarineDesign {
		emitStr(s, submarine.X, screenLine, submarineStyle, submarineLine)

		screenLine += 1
	}
}

func render(s tcell.Screen) {
	s.Clear()
	drawFish(s, fish)
	drawArtifact(s, artifact)
	drawSubmarine(s, submarine)
	s.Show()
}

func renderLoop(s tcell.Screen, quit <-chan struct{}) {
	t := time.NewTicker(300 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			s.Sync()
			render(s)
		case <-quit:
			return
		}
	}
}

func eventLoop(s tcell.Screen, quit chan struct{}) {
	defer close(quit)

	for {

		switch ev := s.PollEvent().(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC:
				close(quit)
				s.Fini()
				os.Exit(0)
			case tcell.KeyCtrlL:
				s.Sync()
			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

func main() {

	fmt.Fprintf(os.Stdout, "Gui started..:\n")

	os.Setenv("TERM", "linux")

	artifact.X = -1
	artifact.Y = -1
	quit := make(chan struct{})

	fishStyle = tcell.StyleDefault.Foreground(tcell.ColorGreen)
	submarineStyle = tcell.StyleDefault.Foreground(tcell.ColorYellow)
	artifactStyle = tcell.StyleDefault.
		Foreground(tcell.ColorDarkRed).
		Background(tcell.ColorRed)
	defStyle = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)

	encoding.Register()

	tty, err := tcell.NewDevTtyFromDev(ncursesTTY)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create new tty device: %v\n", err)
	}
	s, e := tcell.NewTerminfoScreenFromTty(tty)

	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	s.SetStyle(defStyle)
	s.HideCursor()

	go initTcpClient()

	go eventLoop(s, quit)
	renderLoop(s, quit)
}
