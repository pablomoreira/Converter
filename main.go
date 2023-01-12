package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"gopkg.in/yaml.v2"
)

func main() {
	yfile, err := ioutil.ReadFile("config.yaml")

	if err != nil {

		log.Fatal(err)
		fmt.Scanln()
	}
	config := make(map[string]string)
	err = yaml.Unmarshal(yfile, &config)

	if err != nil {

		log.Fatal(err)
		fmt.Scanln()
	}

	for k, v := range config {

		fmt.Printf("%s -> %s\n", k, v)
	}
	fmt.Scanln()
	/* := app.New()
	w := a.NewWindow("Converter")

	content1 := widget.NewLabel("text1")

	str := binding.NewString()
	str.Set("Data binding")
	content3 := widget.NewLabelWithData(str)

	data := binding.BindStringList(
		&[]string{},
	)

	listL := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	viewR := widget.NewLabel("viewr")

	contentView := container.New(layout.NewGridLayoutWithColumns(2), listL, viewR)
	content := container.NewBorder(content1, content3, nil, nil, contentView)
	w.Resize(fyne.Size{Height: 320, Width: 480})
	w.SetContent(content)
	*/
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//go timer(str, ctx)
	signal := make(chan byte)

	//go _backW1(ctx, singnals)

	//go _backWDir(config["_spath"], config["_dpath"], config["_kw"], config["_args"], data, ctx)
	go _backWDir(config["_spath"], config["_dpath"], config["_kw"], config["_args"], signal, ctx)
	//w.ShowAndRun()

	for _loop := true; _loop == true; {
		var file *os.File
		time.Sleep(time.Millisecond * 100)
		select {
		case <-signal:
			log.Print("-")
			_loop = false
			file.Close()
		default:
			log.Print(".")

			if file == nil {
				file, err = os.Open("progress.log") // For read access.
				if err != nil {
					log.Panic(err)
				}
			}
			fileScanner := bufio.NewScanner(file)
			fileScanner.Split(bufio.ScanLines)
			for fileScanner.Scan() {
				out := strings.Split(fileScanner.Text(), "=")
				if out[0] == "speed" {
					fmt.Println(out[1])
				}
			}
		}
	}

	time.Sleep(time.Second * 2)
	//singnals <- 1
	//cancel()
}

func timer(str binding.String, ctx context.Context) {
	for true {
		time.Sleep(time.Second)
		//str.Set(time.Now().String())
	}
}

func _backW1(ctx context.Context, signal chan int) {
	i := 0
	for i == 0 {
		select {
		case <-signal:
			log.Print("-")
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
	}
}

// func _backWDir(_spath string, _dpath string, _kw string, _args string, _data binding.ExternalStringList, _ctx context.Context) {
// 	time.Sleep(time.Millisecond * 300)
// 	files, err := ioutil.ReadDir(_spath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, file := range files {
// 		if file.IsDir() == false && len(strings.Split(file.Name(), ".")) == 2 {
// 			if strings.ToLower(strings.Split(file.Name(), ".")[1]) == "mov" {
// 				ouputfile := strings.ToLower(strings.Split(file.Name(), ".")[0])
// 				err = ffmpeg.Input(_spath+file.Name()).
// 					Output(_dpath+ouputfile+".mp4", ffmpeg.KwArgs{"c:v": "h264_amf", "vf": "scale=1024x720", "r": "30"}).
// 					OverWriteOutput().ErrorToStdOut().Run()
// 				if err != nil {
// 					log.Panic("conveter")
// 				}
// 				//_data.Append(_dpath + ouputfile + ".mp4")
// 			}
// 		}
// 	}
// 	time.Sleep(time.Millisecond * 300)
// }

func _backWDir(_spath string, _dpath string, _kw string, _args string, signal chan byte, _ctx context.Context) {
	//time.Sleep(time.Millisecond * 300)
	files, err := ioutil.ReadDir(_spath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() == false && len(strings.Split(file.Name(), ".")) == 2 {
			if strings.ToLower(strings.Split(file.Name(), ".")[1]) == "mov" {
				ouputfile := strings.Split(file.Name(), ".")[0] + ".mp4"
				inputfile := _spath + file.Name()
				log.Printf("%s -> %s\n", inputfile, ouputfile)

				cmd := exec.Command("ffmpeg", "-y", "-i", inputfile, ouputfile, "-v", "0", "-progress", "progress.log")
				//cmd.Stdout = os.Stdout

				if err := cmd.Run(); err != nil {
					fmt.Println("could not run command: ", err)
				}
			}
		}

	}
	signal <- 1
}
