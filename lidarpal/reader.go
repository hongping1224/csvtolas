package lidarpal

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/hongping1224/lidario"
)

//Reader read Pointcloud
type Reader struct {
	scanner *bufio.Scanner
	wg      *sync.WaitGroup
}

//NewReader Create a new Reader
func NewReader(scanner *bufio.Scanner, wg *sync.WaitGroup) *Reader {
	return &Reader{scanner: scanner, wg: wg}
}

//Read point into channel
func (read *Reader) Read(input chan<- lidario.LasPointer) {
	for read.scanner.Scan() {
		data := strings.Split(read.scanner.Text(), " ")
		if len(data) < 4 {
			continue
		}
		x, err := strconv.ParseFloat(data[0], 64)
		if err != nil {
			continue
		}
		y, err := strconv.ParseFloat(data[1], 64)
		if err != nil {
			continue
		}
		z, err := strconv.ParseFloat(data[2], 64)
		if err != nil {
			continue
		}
		source, err := strconv.Atoi(data[3])
		if err != nil {
			continue
		}
		//Source +1 make sure point source start from 1
		p := lidario.PointRecord0{X: x, Y: y, Z: z, PointSourceID: uint16(source + 1)}
		input <- &p
	}
	read.wg.Done()
	fmt.Println("Reader Done")
}

//Serve read concerrently
func (read *Reader) Serve(input chan<- lidario.LasPointer) {
	go read.Read(input)
}
