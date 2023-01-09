package main

import (
	"context"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	content1 := widget.NewLabel("text1")

	str := binding.NewString()
	str.Set("Data binding")
	content3 := widget.NewLabelWithData(str)

	componentsList := []string{}
	viewL := widget.NewList(
		func() int { return len(componentsList) },
		func() fyne.CanvasObject { return widget.NewLabel("text") },
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(componentsList[lii])
		},
	)
	//componentsList = append(componentsList, "2323")
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		//fmt.Println(file.Name(), file.IsDir())
		if file.IsDir() == false && len(strings.Split(file.Name(), ".")) == 2 {
			if strings.Split(file.Name(), ".")[1] == "mov" {
				//log.Print()

				componentsList = append(componentsList, file.Name())
			}
		}
	}
	viewR := widget.NewLabel("viewr")

	contentView := container.New(layout.NewGridLayoutWithColumns(2), viewL, viewR)
	//contentViewMax := container.New(layout.(3), contentView)
	//content := container.New(layout.NewVBoxLayout(), content1, contentView, content3)
	content := container.NewBorder(content1, content3, nil, nil, contentView)
	w.Resize(fyne.Size{Height: 320, Width: 480})
	w.SetContent(content)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go timer(str, ctx)
	singnals := make(chan int)
	go _backW1(ctx, singnals)

	w.ShowAndRun()

	time.Sleep(time.Second * 2)
	singnals <- 1
	cancel()
}

func timer(str binding.String, ctx context.Context) {
	for true {
		time.Sleep(time.Second)
		str.Set(time.Now().String())
	}
}

func _backW1(ctx context.Context, signal chan int) {
	i := 0
	for i == 0 {
		select {
		case <-signal:
			log.Print("-")
			i = 1
		default:
			//	log.Print(".")
			time.Sleep(time.Millisecond * 250)
		}
	}
}
