package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
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
)

func initTcpClient() {
	// fmt.Println("Client started...")
	for {
		connectedSync.Lock()
		alreadyConnected := connected
		connectedSync.Unlock()
		if !alreadyConnected {
			conn, err := net.Dial("tcp", "127.0.0.1:8000")
			if err != nil {
				fmt.Println(err.Error())
				time.Sleep(time.Duration(5) * time.Second)
				continue
			}
			tcpConn = conn
			// fmt.Println(conn.RemoteAddr().String() + ": connected")
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
			// fmt.Println(conn.RemoteAddr().String() + ": disconnected")
			conn.Close()
			connectedSync.Lock()
			connected = false
			connectedSync.Unlock()
			// fmt.Println(conn.RemoteAddr().String() + ": end receiving data")
			return
		}

		var newObj GuiObject
		err = json.Unmarshal([]byte(message), &newObj)
		if err != nil {
			fmt.Println(err.Error())
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
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
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
			case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyCtrlC:
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

	s, e := tcell.NewScreen()
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
	s.Clear()

	go initTcpClient()

	go renderLoop(s, quit)
	eventLoop(s, quit)
}

func printStruct() {
	fmt.Printf("%+v\n", submarine)
	fmt.Printf("%+v\n", artifact)
	fmt.Printf("%+v\n", fish)
	fmt.Printf("\n")

}
