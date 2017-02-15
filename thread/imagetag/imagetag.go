package main

import (
	"os"
	"fmt"
	"runtime"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path"
)

var workers = runtime.NumCPU()

func main() {
	runtime.GOMAXPROCS(workers)
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s [<imagePath>...]\n", os.Args[0])
		os.Exit(1)
	}
	fileInfoChan := make(chan fileInfo)
	go processInput(fileInfoChan, os.Args[1:])
	processOutput(fileInfoChan)
}

type fileInfo struct {
	imgName	string
	width	string
	height	string
}

func processInput(fileInfoChan chan <- fileInfo, paths []string) {
	for _, filePath := range paths {
		go process(fileInfoChan, filePath)

	}
	close(fileInfoChan)
}

func process(fileInfoChan chan <- fileInfo, filePath string)  {
	fileinfo, err := os.Stat(filePath)
	if err != nil || (fileinfo.Mode() & os.ModeType) != 0 {
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	conf, _, err := image.DecodeConfig(file)
	if err != nil {
		return
	}
	fi := fileInfo{
		imgName:path.Base(filePath),
		width: fmt.Sprintf("%d", conf.Width),
		height:fmt.Sprintf("%d", conf.Height),
	}
	fileInfoChan <- fi
}

func processOutput(fileInfoChan chan fileInfo) {
	for fileinfo := range fileInfoChan {
		fmt.Printf(`<img src="%s" width="%s" height="%s" />`, fileinfo.imgName, fileinfo.width, fileinfo.height)
	}
}