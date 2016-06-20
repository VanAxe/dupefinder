package main

import (
	"bitbucket.org/alexthekone/dupefinder/format"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
)

var (
	app      = kingpin.New("dupefinder", "Find duplicate files with checksums.")
	verbose  = app.Flag("verbose", "Enable verbose mode.").Short('v').Bool()
	debug    = app.Flag("debug", "Enable debug mode.").Short('d').Bool()
	rec      = app.Flag("recursive", "Search recursively.").Short('r').Bool()
	output   = app.Flag("output", "File to which the report will be written.").Short('o').Default("./dupefinder.csv").String()
	target   = app.Arg("target", "Where to look for duplicate files.").Default(".").String()
	abs_path = ""
	file_map = make(map[string][]byte)
)

func main() {
	app.Version("1.0.0")

	kingpin.MustParse(app.Parse(os.Args[1:]))
	abs_path, _ = filepath.Abs(*target)

	report_dest, report_file := filepath.Split(*output)
	report_dest, _ = filepath.Abs(report_dest)

	format.Print_debug(*verbose, debug, rec, report_dest, report_file, abs_path)
	find_dupes()
	output_map := reduce_duplicates(file_map)
	csv := format.Pp_csv(output_map)
	if *verbose {
		fmt.Println(csv)
	}

	data := []byte(csv)
	if len(output_map) > 0 {
		err := ioutil.WriteFile(report_dest+"/"+report_file, data, 0644)
		check(err)
	} else {
		fmt.Println("No duplicates found.")
	}

}

func find_dupes() {
	if *rec == true {
		find_dupes_rec(abs_path)
	} else {
		find_dupes_flat(abs_path)
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
	for hash_key, file_refs := range output_map {
		sort.Strings(file_refs)
		output_map[hash_key] = file_refs
	}
	return output_map
}

func find_dupes_rec(abs_path string) {
	err := filepath.Walk(abs_path, visit)
	check(err)
}

func find_dupes_flat(abs_path string) {
	files, err := ioutil.ReadDir(abs_path)
	check(err)

	for _, file := range files {
		if !file.IsDir() {
			fullpath := abs_path + "/" + file.Name()
			md5_hash, _ := ComputeMd5(fullpath)
			file_map[fullpath] = md5_hash
		}
	}
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
