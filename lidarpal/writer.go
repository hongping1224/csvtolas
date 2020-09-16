package lidarpal

import (
	"fmt"
	"time"

	"github.com/hongping1224/lidario"
)

//Writer structer to write point cloud
type Writer struct {
	input chan lidario.LasPointer
	Busy  bool
	las   *lidario.LasFile
}

//NewWriter Create a writer with channel to write
func NewWriter(input chan lidario.LasPointer) *Writer {
	return &Writer{input: input, Busy: false}
}

//Write Point p into buffer
func (writer *Writer) Write(p lidario.LasPointer) {
	writer.input <- p
}

//Serve writer in background
func (writer *Writer) Serve(las *lidario.LasFile) {
	writer.Busy = true
	writer.las = las
	go func() {
		for {
			a, open := <-writer.input

			if open == false {
				writer.Busy = false
				//fmt.Println("Writer Closing")
				break
			}
			las.AddLasPoint(a)
		}
	}()
}

//Close writer
func (writer *Writer) Close() {
	close(writer.input)
	for writer.Busy == true {
		time.Sleep(100 * time.Millisecond)
	}
	err := writer.las.Close()
	fmt.Println(err)
}
