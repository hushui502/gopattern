package main

import "fmt"

/*
	Delegation e.g1
*/
type Widget struct {
	X, Y int
}

type Label struct {
	Widget
	Text string
}

func (label Label) Paint() {
	fmt.Println("%p: Label.Paint(%q)\n", &label, label.Text)
}

type Painter interface {
	Paint()
}

type Clicker interface {
	Click()
}

type Button struct {
	Label
}

func NewButton(x, y int, text string) Button {
	return Button{Label{Widget{x, y}, text}}
}

func (button Button) Paint() {
	fmt.Println("Button.Paint(%s)\n", button.Text)
}

func (button Button) Click() {
	fmt.Println("Button.Click(%s)\n", button.Text)
}

type ListBox struct {
	Widget
	Texts []string
	Index int
}

func (listBox ListBox) Paint() {
	fmt.Println("ListBox.Paint(%q)\n", listBox.Texts)
}

func (listBox ListBox) Click() {
	fmt.Println("ListBox.Click(%q)\n", listBox.Texts)
}

func main() {
	button1 := Button{Label{Widget{10, 10}, "OK"}}
	button2 := Button{Label{Widget{20, 20}, "Cancel"}}
	listBox := ListBox{Widget{13, 23},
		[]string{"AL", "AZ", "AR", "AK"}, 0}

	for _, painter := range []Painter{listBox, button1, button2} {
		painter.Paint()
	}

	for _, widget := range []interface{}{listBox, button1, button2} {
		if clicker, ok := widget.(Clicker); ok {
			clicker.Click()
		}
	}
}