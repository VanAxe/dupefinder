package main

import (
	"crypto/md5"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	app      = kingpin.New("dupefinder", "Find duplicate files with checksums.")
	verbose  = app.Flag("verbose", "Enable verbose mode.").Short('v').Bool()
	debug    = app.Flag("debug", "Enable debug mode.").Short('d').Bool()
	rec      = app.Flag("recursive", "Search recursively.").Short('r').Bool()
	output   = app.Flag("output", "File to which the report will be written.").Short('o').Default("./dupefinder.log").String()
	target   = app.Arg("target", "Where to look for duplicate files.").Default(".").String()
	abs_path = ""
)

func main() {
	app.Version("0.0.1")

	kingpin.MustParse(app.Parse(os.Args[1:]))
	abs_path, _ = filepath.Abs(*target)

	if *debug == true {
		conf := `== DEBUG CONF ==
Verbose:		%t
Debug:			%t
Recursive:		%t
Output:			%s
Target:			%s

`
		fmt.Printf(conf, *verbose, *debug, *rec, *output, abs_path)
	}
	find_dupes()
}

func find_dupes() {
	if *rec == true {
		err := filepath.Walk(abs_path, visit)
		if err != nil {
			fmt.Printf("Some error! %v\n", err)
		} else {
			if *verbose {
				fmt.Println("Done!")
			}
		}
	} else {
		files, err := ioutil.ReadDir(abs_path)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if !file.IsDir() {
				fullpath := abs_path + "/" + file.Name()
				md5, _ := ComputeMd5(fullpath)
				if *verbose {
					fmt.Printf("%s [ %x ]\n", fullpath, md5)
				}
			} else {
				if *verbose {
					fmt.Printf("%s is a directory\n", abs_path+"/"+file.Name())
				}
			}
		}
		if *verbose {
			fmt.Println("Done!")
		}
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}

func visit(path string, f os.FileInfo, err error) error {
	isDir, _ := isDirectory(path)
	if !isDir && *verbose {
		md5, _ := ComputeMd5(path)
		fmt.Printf("%s [ %x ]\n", path, md5)
	}
	return nil
}

func ComputeMd5(filePath string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}

	return hash.Sum(result), nil
}
