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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signal := make(chan byte)
	go _backWDir(config["_spath"], config["_dpath"], signal, ctx)

	end := <-signal
	log.Print(end)
	cancel()
}

func _backWDir(_spath string, _dpath string, signal chan byte, _ctx context.Context) {
	files, err := ioutil.ReadDir(_spath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() == false && len(strings.Split(file.Name(), ".")) == 2 {
			if strings.ToLower(strings.Split(file.Name(), ".")[1]) == "mov" {
				ouputfile := strings.Split(file.Name(), ".")[0]
				inputfile := _spath + file.Name()

				cmd_probe := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0",
					"-count_packets", "-show_entries", "stream=nb_read_packets", "-of", "csv=p=0", inputfile)

				cmd_probe.Stderr = os.Stderr
				data, err := cmd_probe.Output()
				if err != nil {
					log.Fatalf("failed to call Output(): %v", err)
				}
				size := string(data)
				log.Printf("%s -> %s Frame=%s\n", inputfile, ouputfile+".mp4", size)

				cmd := exec.Command("ffmpeg", "-y", "-i", inputfile, ouputfile+".mp4", "-v", "0", "-progress", ouputfile+".log")

				go status_bar(size, ouputfile+".log", _ctx)

				if err := cmd.Run(); err != nil {
					fmt.Println("ffmpeg could not run command: ", err)
				}
			}
		}

	}
	signal <- 1
}

func status_bar(size string, fileName string, ctx context.Context) {

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for i := 1; scanner.Scan(); i++ {
		fmt.Printf("%d) \"%s\"\n", i, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
