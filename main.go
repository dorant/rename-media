package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/barsanuphe/goexiftool"
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

func hasExif(extension string) bool {
	switch strings.ToUpper(extension) {
	case ".MTS", ".JPG", ".3GP":
		return true
	}
	return false
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
		fmt.Printf("Error: %s\n", err)
		return
	}

	for _, file := range files {

		var t time.Time

		if hasExif(path.Ext(file)) {
			m, err := goexiftool.NewMediaFile(file)
			if err != nil {
				fmt.Printf("Error while probing file '%s': %v\n", file, err)
				continue
			}

			if *showInfo {
				fmt.Printf("%s: %s\n", file, m)
				continue
			}

			val, err := m.Get("Date/Time Original")
			if err != nil {
				fmt.Printf("Error while getting date from file '%s'. Trying other Exif date..\n", file)
				val, err = m.Get("File Modification Date/Time")
				if err != nil {
					fmt.Printf("Error while getting date from file '%s': %v\n", file, err)
					continue
				}
			}
			t, err = time.Parse("2006:01:02 15:04:05-07:00 MST", val)
			if err != nil {
				t, err = time.Parse("2006:01:02 15:04:05-07:00", val)
				if err != nil {
					t, err = time.Parse("2006:01:02 15:04:05", val)
					if err != nil {
						fmt.Printf("Error while parsing date from file '%s': %v\n", file, err)
						continue
					}
				}
			}

		} else {
			// Try getting metadata for other filetypes

			data, err := ffprobe.GetProbeData(file, 4000*time.Millisecond)
			if err != nil {
				fmt.Printf("Error while probing file '%s': %v\n", file, err)
				continue
			}

			if *showInfo {
				buf, _ := json.MarshalIndent(data, "", "  ")
				fmt.Printf("%s: %s\n", file, string(buf))
				continue
			}

			t, err = time.Parse(time.RFC3339, data.Format.Tags.CreationTime)
			if err != nil {
				fmt.Printf("Error converting date from file '%s': %v\n", file, err)
				continue
			}
		}

		newname := fmt.Sprintf("%d-%02d-%02d_%02d.%02d_%s",
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(),
			filepath.Base(file))
		newname = filepath.Join(filepath.Dir(file), newname)

		if *force {
			fmt.Printf("Renaming: %s to %s\n", file, newname)
			err = os.Rename(file, newname)
			if err != nil {
				fmt.Printf("Error renaming file '%s': %v\n", file, err)
				continue
			}
		} else {
			fmt.Printf("DRYRUN: Renaming: %s to %s\n", file, newname)
		}
	}
}
