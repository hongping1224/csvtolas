package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/jblindsay/lidario"

	"github.com/hongping1224/csvtolas/lidarpal"
)

var numOFCPU int

func main() {
	numOFCPU = runtime.NumCPU()
	flag.IntVar(&numOFCPU, "cpuCount", numOFCPU, "Cpu use for compute")
	dir := "./"
	//dir := "W:\\data_gravel_v3\\helios_output"
	flag.StringVar(&dir, "dir", dir, "directory to process")
	flag.Parse()
	//check directory exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatal(err)
		return
	}
	//find all las file
	fmt.Println(dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		fullpath := filepath.Join(dir, f.Name())
		if fi, _ := os.Stat(fullpath); !fi.Mode().IsDir() {
			continue
		}
		xyzs := findFile(fullpath, ".xyz")
		if len(xyzs) == 0 {
			continue
		}
		convert(filepath.Join(fullpath, f.Name()+".las"), xyzs)
	}

}

func convert(laspath string, xyzs []string) error {
	headerExample, err := lidario.NewLasFile("./headersample.las", "rh")
	if err != nil {
		return err
	}
	fmt.Println(laspath)
	// open las file
	las, err := lidario.InitializeUsingFile(laspath, headerExample)
	las.Header.PointFormatID = 0
	headerExample.Close()
	if err != nil {
		return err
	}
	writechan := make(chan lidario.LasPointer, numOFCPU*4)
	writer := lidarpal.NewWriter(writechan)
	writer.Serve(las)
	readers := make([]*lidarpal.Reader, len(xyzs))
	files := make([]*os.File, len(xyzs))
	var wg sync.WaitGroup
	for i, xyz := range xyzs {
		files[i], err = os.Open(xyz)
		if err != nil {
			return err
		}
		//setup reader for each file
		scanner := bufio.NewScanner(files[i])
		wg.Add(1)
		readers[i] = lidarpal.NewReader(scanner, &wg)
		readers[i].Serve(writechan)
	}
	//wait all reader done
	wg.Wait()
	//wait writer done
	writer.Close()
	//clean up
	for _, file := range files {
		file.Close()
	}
	return nil
}

func findFile(root string, match string) (file []string) {
	fmt.Println("Finding " + match + " File in :")
	fmt.Println(root)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			//fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return nil
		}

		if strings.Contains(strings.ToLower(info.Name()), match) {
			file = append(file, path)
			return nil
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Found", len(file), match+" file")
	return file
}
