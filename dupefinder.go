package main

import (
	"crypto/md5"
	"encoding/hex"
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
	file_map = make(map[string][]byte)
)

func main() {
	app.Version("0.0.1")

	kingpin.MustParse(app.Parse(os.Args[1:]))
	abs_path, err := filepath.Abs(*target)
	if err != nil {
		log.Fatal(err)
	}
	// Nice snafu : fp.Abs needs to check something that exists
	// gotta separate report_name from report_dest and pass that down...
	report_dest, err := filepath.Abs(*output)
	if err != nil {
		log.Fatal(err)
	}

	if *debug == true {
		conf := `== DEBUG CONF ==
Verbose:		%t
Debug:			%t
Recursive:		%t
Output:			%s
Target:			%s

`
		fmt.Printf(conf, *verbose, *debug, *rec, output, abs_path)
	}
	find_dupes()
}

func find_dupes() {
	if *rec == true {
		err := filepath.Walk(abs_path, visit)
		if err != nil {
			fmt.Printf("Some error! %v\n", err)
		} else {
			output_map := reduce_duplicates(file_map)
			csv := pp_csv(output_map)
			// print to file with []byte(csv)
			if *verbose {
				fmt.Printf("%s", csv)
				// fmt.Println("Done!")
			}
		}
	} else {
		files, err := ioutil.ReadDir(abs_path)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range files {
			if !file.IsDir() {
				// fullpath := abs_path + "/" + file.Name()
				// md5_hash, _ := ComputeMd5(fullpath)
				// if *verbose {
				// 	fmt.Printf("%s [ %x ]\n", fullpath, md5_hash)
				// }
			} else {
				// if *verbose {
				// 	fmt.Printf("%s is a directory\n", abs_path+"/"+file.Name())
				// }
			}
		}
		// if *verbose {
		// 	fmt.Println("Done!")
		// }
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}

func visit(path string, f os.FileInfo, err error) error {
	isDir, _ := isDirectory(path)
	if !isDir {
		md5_hash, _ := ComputeMd5(path)
		file_map[path] = md5_hash
		// if *verbose {
		// 	fmt.Printf("%s [ %x ]\n", path, md5_hash)
		// }
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

// type file_entry struct {
// 	hash    string
// 	matches []string
// }

func reduce_duplicates(input_map map[string][]byte) map[string][]string {
	output_map := make(map[string][]string)
	for file_key, md5_hash := range input_map {
		md5s := hex.EncodeToString(md5_hash)
		output_map[md5s] = append(output_map[md5s], file_key)
	}
	for md5_key, file_array := range output_map {
		if len(file_array) == 1 {
			delete(output_map, md5_key)
		}
	}
	return output_map
}

func pp_csv(input_map map[string][]string) string {
	pretty_csv := ""
	for entry, files := range input_map {
		string_of_files := ""
		for _, file := range files {
			string_of_files = string_of_files + ", " + file
		}
		pretty_csv = pretty_csv + entry + string_of_files + "\n"
		// fmt.Printf("%s, %s\n ", entry, string_of_files)
	}
	return pretty_csv
}
