package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
	"github.com/nxadm/tail" 

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

	
	signal := make(chan byte)
	size := make(chan string)

	//go _backW1(ctx, singnals)

	//go _backWDir(config["_spath"], config["_dpath"], config["_kw"], config["_args"], data, ctx)
	go _backWDir(config["_spath"], config["_dpath"], signal, size, ctx)
	//w.ShowAndRun()
	
	t, err := tail.TailFile(
		"progress.log", tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		log.Print(err)
	}


	Size :=""
	for _loop := true; _loop == true; {
		select {
		case <-signal:
			log.Print("-")
			_loop = false

		case tmp := <-size:
				log.Print("FRAME > "+ tmp)
				Size = tmp

		default:
			time.Sleep(time.Millisecond * 10)
			
			for line := range t.Lines {
				__line := strings.Split(line.Text,"=")
				if __line[0] == "frame"{
					log.Println(__line[1],Size)
				}
				break
			}
		}
	}

	time.Sleep(time.Second * 2)
	
	//singnals <- 1
	//cancel()
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

func _backWDir(_spath string, _dpath string, signal chan byte, size chan string, _ctx context.Context) {
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

				cmd_probe := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0",
					"-count_packets", "-show_entries", "stream=nb_read_packets", "-of", "csv=p=0", inputfile)

				cmd_probe.Stderr = os.Stderr
				data, err := cmd_probe.Output()
				if err != nil {
					log.Fatalf("failed to call Output(): %v", err)
				}
				Size := string(data)
				size <- Size

				log.Printf("%s -> %s Frame=%s\n", inputfile, ouputfile,Size )

				cmd := exec.Command("ffmpeg", "-y", "-i", inputfile, ouputfile, "-v", "0", "-progress", "progress.log")

				if err := cmd.Run(); err != nil {
					fmt.Println("ffmpeg could not run command: ", err)
				}
			}
		}

	}
	signal <- 1
}
