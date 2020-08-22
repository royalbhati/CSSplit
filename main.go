package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/royalbhati/cssplit/extract"
	"github.com/royalbhati/cssplit/walk"
)

func usage() {
	meta :=
		`
NAME: 
	CSSplit - Split your one big CSS files into modules

USAGE: 
	cssplit -dir ~/Desktop/myproject -css ~/Desktop/myproject/main.css -ext ".html,.js,.jsx" -exclude "utils,lib"
`

	fmt.Println(meta)
}

func main() {
	usage()
	dir := flag.String("dir", ".", "Project Directory")
	csspath := flag.String("css", "", "Path to CSS file")
	ext := flag.String("ext", "[.html,.js,.jsx]", "File extensions to parse")
	exclude := flag.String("exclude", "utils", "Files to exclude")

	flag.Parse()

	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		log.Fatal("Invalid Path")
	}

	if *csspath == "" {
		log.Fatal("Path to css file must be provided")
	}

	exts := strings.Split(*ext, ",")
	exc := strings.Split(*exclude, ",")
	dirss, _ := walk.GetDir(*dir, exc, exts)

	var wg sync.WaitGroup
	for _, v := range dirss {
		wg.Add(1)
		go extract.CSS(v, *csspath, &wg)
	}

	wg.Wait()

}
