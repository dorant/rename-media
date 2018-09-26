package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	ffprobe "github.com/vansante/go-ffprobe"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name!")
		return
	}
	file := os.Args[1]

	if _, err := os.Stat(file); err != nil {
		fmt.Printf("File error: %s\n", err)
		return
	}

	data, err := ffprobe.GetProbeData(file, 500*time.Millisecond)
	if err != nil {
		log.Panicf("Error getting data: %v", err)
	}

	t, err := time.Parse(time.RFC3339, data.Format.Tags.CreationTime)
	if err != nil {
		log.Panicf("Error converting date: %v", err)
	}

	newname := fmt.Sprintf("%d-%02d-%02d_%02d.%02d_%s",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), filepath.Base(file))

	err = os.Rename(file, newname)
	if err != nil {
		log.Panicf("Error renaming file: %s, err: %v", file, err)
	}
}
