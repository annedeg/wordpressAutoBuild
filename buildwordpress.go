package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"unsafe"
)

func main() {
	dir := flag.String("dir", "", "Use for other directory")
	flag.Parse()
	result := int8(0)
	fmt.Print("started downloading prereq\n");
	//install prerequisitions
	out, err := exec.Command("/bin/sh", "test.sh").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(BytesToString(out))
	fmt.Print("done\n");
	if dir != nil {
		result = getWordpress(*dir)
	} else {
		result = getWordpress("")
	}
	if result == 1 {
		fmt.Print("Downloaded wordpress\n")
	}
}

func getWordpress(dir string) int8 {
	if dir == "" {
		curPath, err := os.Getwd()
		if err == nil {
			err := download(curPath+"/latest.zip", "https://wordpress.org/latest.zip")
			if err == nil {
				_, err = Unzip(curPath+"/latest.zip", curPath)
				if err != nil {
					fmt.Print("Unzipped wordpress\n")
				}
			} else {
				fmt.Print(err)
			}
		}
	}
	return -1
}

func download(filepath string, url string) error {
	fmt.Print("Downloading: " + url + "\n")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func Unzip(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}
