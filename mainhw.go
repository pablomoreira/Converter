package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/nxadm/tail"
	"github.com/schollz/progressbar/v3"
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

	/*for k, v := range config {

		fmt.Printf("%s -> %s\n", k, v)
	}*/
	//time.Sleep(time.Second * 2)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signal := make(chan byte)
	go _backWDir(config["_spath"], config["_dpath"], signal, ctx)

	end := <-signal
	log.Print(end, " Conversion finished")
	cancel()
}

func _backWDir(_spath string, _dpath string, signal chan byte, _ctx context.Context) {
	files, err := ioutil.ReadDir(_spath)
	if err != nil {
		log.Fatal(err)
	}
	count := 0
	for _, file := range files {
		if file.IsDir() == false && len(strings.Split(file.Name(), ".")) == 2 {
			if strings.ToLower(strings.Split(file.Name(), ".")[1]) == "mov" {
				outputfile := _dpath + strings.Split(file.Name(), ".")[0]
				inputfile := _spath + file.Name()

				cmd_probe := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0",
					"-count_packets", "-show_entries", "stream=nb_read_packets", "-of", "csv=p=0", inputfile)

				cmd_probe.Stderr = os.Stderr
				data, err := cmd_probe.Output()
				if err != nil {
					log.Fatalf("failed to call Output(): %v", err)
				}
				size := strings.TrimRight(string(data), "\r\n")
				log.Printf("%s -> %s Frame=%s\n", inputfile, outputfile+".mp4", size)

				cmd := exec.Command("ffmpeg", "-y", "-i", inputfile, "-c:v", "h264_amf", outputfile+".mp4", "-v", "0", "-progress", outputfile+".log")
				count++
				delFileIfExist(outputfile + ".log")
				go status_bar(size, outputfile+".log", _ctx)

				if err := cmd.Run(); err != nil {
					fmt.Println("ffmpeg could not run command: ", err)
				}
				delFileIfExist(outputfile + ".log")
			}
		}
		time.Sleep(time.Second * 1)

	}
	time.Sleep(time.Millisecond * 500)
	signal <- byte(count)
}

func status_bar(size string, fileName string, ctx context.Context) {

	var file *os.File
	var err error

	for loop := true; loop == true; {
		file, err = os.Open(fileName)

		if err != nil {
			//log.Print(err)
		} else {
			loop = false
		}
		time.Sleep(time.Millisecond * 50)
	}
	file.Close()
	// Create a tail
	t, err := tail.TailFile(
		fileName, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		panic(err)
	}
	defer t.Stop()

	// Print the text of each received line
	na, _ := strconv.Atoi(size)
	bar := progressbar.Default(int64(na))
	old_nb := 0

	for line := range t.Lines {
		if a := strings.Split(line.Text, "=")[0]; a == "frame" {
			b := strings.Split(line.Text, "=")[1]
			nb, err := strconv.Atoi(b)
			if err != nil {
				log.Panicln(err)
			}
			//fmt.Println(nb*100/na, na)
			barp := nb - old_nb
			bar.Add(int(barp))
			//log.Print(barp)
			old_nb = nb

		}
	}
}
func delFileIfExist(filename string) {

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {

	} else {
		os.Remove(filename)
	}
	f, _ := OpenFile(filename)
	f.Close()
}

func OpenFile(name string) (file *os.File, err error) {
	return os.OpenFile(name, os.O_RDONLY, 0)
}
