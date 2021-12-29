package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/mattn/go-runewidth"
)

var defStyle tcell.Style
var submarine = []string{
	"          __|___",
	"         /      \\",
	" _______/    O   \\_______",
	"<                        \\_____  I",
	" \\   O      O     O            >-=",
	"  \\___________________________/  I"}

const (
	fish = "><(((*>"
	artifact = "[*]"
	submarineHeight = 6
	submarineLength = 35
)

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

func drawFish(s tcell.Screen, x, y int) {

	style := tcell.StyleDefault.
		Foreground(tcell.ColorGreen)
	emitStr(s, x, y, style, fish)
}

func drawArtifact(s tcell.Screen, x, y int) {

	style := tcell.StyleDefault.
		Foreground(tcell.ColorDarkRed).Background(tcell.ColorRed)
	emitStr(s, x, y, style, artifact)
}


func drawSubmarine(s tcell.Screen, x, y int) {

	style := tcell.StyleDefault.
		Foreground(tcell.ColorYellow)

	screenLine := y
	for _, submarineLine := range submarine {
		emitStr(s, x, screenLine, style, submarineLine)
		screenLine += 1
	}
}

func drawSelect(s tcell.Screen, x1, y1, x2, y2 int, sel bool) {

	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			mainc, combc, style, width := s.GetContent(col, row)
			if style == tcell.StyleDefault {
				style = defStyle
			}
			style = style.Reverse(sel)
			s.SetContent(col, row, mainc, combc, style)
			col += width - 1
		}
	}
}

func main() {

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

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
	defStyle = tcell.StyleDefault.
		Background(tcell.ColorReset).
		Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)
	s.EnableMouse()

	s.EnablePaste()
	s.Clear()

	for {
		
		// w, h = s.Size()
		
		// // always clear any old selection box
		// if ox >= 0 && oy >= 0 && bx >= 0 {
			// 	drawSelect(s, ox, oy, bx, by, false)
			// }
		
		drawFish(s, 10, 15)
		drawArtifact(s, 11, 16)
		drawSubmarine(s, 0, 0)


		switch ev := s.PollEvent().(type) {
		case *tcell.EventResize:
			s.Sync()
			drawFish(s, 10, 15)
			drawArtifact(s, 11, 16)
			drawSubmarine(s, 0, 0)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				os.Exit(0)
			}
		}

		// switch ev := ev.(type) {
		// case *tcell.EventResize:
		// 	s.Sync()
		// 	s.SetContent(w-1, h-1, 'R', nil, st)
		
		// 		default:
		// 	s.SetContent(w-1, h-1, 'X', nil, st)
		// }

		// if ox >= 0 && bx >= 0 {
		// 	drawSelect(s, ox, oy, bx, by, true)
		// }
	}
}
