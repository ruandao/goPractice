package main

import (
	"runtime"
	"regexp"
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"io/ioutil"
	"image"
)

var workers = runtime.NumCPU()

const (
	widthAttr = "width="
	heightAttr = "height="
)

var (
	imageRx *regexp.Regexp
	srcRx	*regexp.Regexp
)

func init() {
	imageRx = regexp.MustCompile(`<[iI][mM][gG][^>]+>`)
	srcRx = regexp.MustCompile(`src=["']([^"'])+["']`)
}

func main() {
	runtime.GOMAXPROCS(workers)
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <html files>\n", os.Args[0])
		os.Exit(1)
	}

	files := commandLineFiles(os.Args[1:])
	jobs := make(chan string, workers * 16)
	done := make(chan struct{}, workers)
	go addJobs(files, jobs)
	for i := 0; i < workers; i++ {
		go doJobs(done, jobs)
	}
	waitUntil(done)
}

func addJobs(files []string, jobs chan <- string) {
	for _, filename := range files {
		suffix := strings.ToLower(filepath.Ext(filename))
		if suffix == ".html" || suffix == ".htm"{
			jobs <- filename
		}
	}
	close(jobs)
}

func doJobs(done chan <- struct{}, jobs <- chan string)  {
	for job := range jobs {
		sizeImages(job)
	}
	done <- struct {}{}
}

func waitUntil(done <- chan struct{})  {
	for i := 0; i < workers; i++ {
		<- done
	}
}

func sizeImages(filename string) {
	if info, err := os.Stat(filename); err != nil || (info.Mode() & os.ModeType) == 1 {
		fmt.Printf("ignoring:%s\n", filename)
		return
	}
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("failed to read: ", filename)
		return
	}
	html := string(raw)
	fmt.Println("reading file: ", filename)
	newHtml := imageRx.ReplaceAllStringFunc(html, sizer)
	if len(html) != len(newHtml) {
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("couldn't update file: %s %v\n", filename, err)
			return
		}
		defer file.Close()
		if _, err := file.WriteString(newHtml); err != nil {
			fmt.Printf("error when update file: %s %v\n", filename, err)
		}
	}
}

func sizer(originalTag string) string {
	tag := originalTag
	if strings.Index(tag, widthAttr) > -1 &&
		strings.Index(tag, heightAttr) > -1 {
		return tag
	}
	match := srcRx.FindStringSubmatch(tag)
	if match == nil {
		fmt.Printf("can't find <img>'s src attribute %s\n",tag)
		return tag
	}
	file, err := os.Open(match[1])
	if err != nil {
		fmt.Printf("can't open <img>'s src file:%s err %v\n", match[1], err)
		return tag
	}
	defer file.Close()
	conf, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Printf("can't decode <img>'s src file:%s err: %v\n", match[1], err)
		return tag
	}
	tag, end := tagEnd(tag)
	if strings.Index(tag, widthAttr) == -1 {
		tag += fmt.Sprintf(`%s"%d"`, widthAttr, conf.Width)
	}
	if strings.Index(tag, heightAttr) == -1 {
		tag += fmt.Sprintf(`%s"%d"`, heightAttr, conf.Height)
	}
	tag += end
	return tag
}
func tagEnd(originTag string) (tag, end string) {
	end = ">"
	tag = originTag[:len(originTag) - 1]
	if tag[len(tag) - 1] == '/' {
		end = "/>"
		tag = tag[:len(tag) - 1]
	}
	return tag, end
}

func commandLineFiles(files []string) []string {
	if runtime.GOOS == "windows" {
		args := make([]string, 0, len(files))
		for _, name := range files {
			if matches, err := filepath.Glob(name); err != nil {
				args = append(args, name)
			} else if matches != nil {
				args = append(args, matches...)
			}
		}
		return args
	}
	return files
}