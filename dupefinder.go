package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app				= kingpin.New("dupefinder", "Find duplicate files with checksums.")
	verbose		= app.Flag("verbose", "Enable verbose mode.").Short('v').Bool()
	debug			= app.Flag("debug", "Enable debug mode.").Short('d').Bool()
	rec				= app.Flag("recursive", "Search recursively.").Short('r').Bool()
	output		= app.Flag("output", "File to which the report will be written.").Short('o').Default("./dupefinder.log").String()
	target		= app.Arg("target", "Where to look for duplicate files.").Default(".").String()
	abs_path,_ = filepath.Abs(*target)
)

func main() {
	app.Version("0.0.1")

	kingpin.MustParse(app.Parse(os.Args[1:]))

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
			if *verbose { fmt.Println("Done!") }
		}
	} else {
		files, err := ioutil.ReadDir(abs_path)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if !file.IsDir() {
				if *verbose { fmt.Println(abs_path + "/" + file.Name()) }
			} else {
				if *verbose { fmt.Printf("%s is a directory\n", abs_path + "/" + file.Name()) }
			}
		}
		if *verbose { fmt.Println("Done!") }
	}
}

func isDirectory(path string) (bool, error) {
    fileInfo, err := os.Stat(path)
    return fileInfo.IsDir(), err
}

func visit(path string, f os.FileInfo, err error) error {
	isDir,_ := isDirectory(path)
	if !isDir && *verbose {
	  fmt.Printf("%s\n", path)
	}
	return nil
}
