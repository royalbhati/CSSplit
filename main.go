package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/royalbhati/cssremove/extract"
	"github.com/royalbhati/cssremove/walk"
)

func usage() {
	u :=
		`
NAME: 
	CSSplit - Split your one big CSS files into modules

USAGE: 
	cssplit -d ~/Desktop/myproject -css ~/Desktop/myproject/main.css -ext ".html,.js,.jsx"

`

	fmt.Println(u)
}

func main() {
	usage()
	dir := flag.String("dir", ".", "Project Directory")
	csspath := flag.String("css", "", "Path to CSS file")
	ext := flag.String("ext", "[.html,.js,.jsx]", "file extensions to parse")
	flag.Parse()

	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		log.Fatal("Invalid Path")
	}

	if *csspath == "" {
		log.Fatal("Path to css file must be provided")
	}

	exts := strings.Split(*ext, ",")

	dirss, _ := walk.GetDir(*dir, exts)

	var wg sync.WaitGroup
	for _, v := range dirss {
		wg.Add(1)
		go extract.CSS(v, *csspath, &wg)
	}

	wg.Wait()

}
