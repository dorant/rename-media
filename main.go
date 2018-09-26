package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	ffprobe "github.com/vansante/go-ffprobe"
)

// Find all files within given paths
func getFiles(paths []string) ([]string, error) {
	var files []string

	for _, path := range paths {

		fi, err := os.Stat(path)
		if err != nil {
			return files, err
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
					return err
				}
				if !info.IsDir() {
					files = append(files, path)
				}
				return nil
			})
			if err != nil {
				return files, err
			}

		case mode.IsRegular():
			files = append(files, path)
		}
	}
	return files, nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name!")
		return
	}

	force := flag.Bool("f", false, "Rename files")
	showInfo := flag.Bool("i", false, "Show video metadata information")
	flag.Parse()

	files, err := getFiles(flag.Args())
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	for _, file := range files {

		data, err := ffprobe.GetProbeData(file, 500*time.Millisecond)
		if err != nil {
			fmt.Printf("Error getting data: %v", err)
		}

		if *showInfo {
			buf, _ := json.MarshalIndent(data, "", "  ")
			fmt.Printf("%s: %s\n", file, string(buf))
			break
		}

		t, err := time.Parse(time.RFC3339, data.Format.Tags.CreationTime)
		if err != nil {
			fmt.Printf("Error converting date: %v", err)
			return
		}

		newname := fmt.Sprintf("%d-%02d-%02d_%02d.%02d_%s",
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(),
			filepath.Base(file))
		newname = filepath.Join(filepath.Dir(file), newname)

		if *force {
			fmt.Printf("FORCE: Rename: %s to %s\n", file, newname)
			err = os.Rename(file, newname)
			if err != nil {
				fmt.Printf("Error renaming file: %s, err: %v", file, err)
				return
			}
		} else {
			fmt.Printf("Rename: %s to %s\n", file, newname)
		}
	}
}
