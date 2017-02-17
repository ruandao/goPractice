package main

import (
	"runtime"
	"os"
	"fmt"
	"io"
	"encoding/binary"
	"unicode/utf16"
	"github.com/juju/errors"
	"bufio"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <srcfile> [dstFile]\n", os.Args[0])
		os.Exit(1)
	}
	output := os.Stdout
	if len(os.Args) == 3 {
		outputfilename := os.Args[2]
		var err error
		output, err = os.Create(outputfilename)
		if err != nil {
			fmt.Printf("create output file %s err %v\n", outputfilename, err)
			os.Exit(1)
		}
	}
	defer output.Close()

	filename := os.Args[1]
	inputfile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("open input file %s err %v\n", filename, err)
		os.Exit(1)
	}
	defer inputfile.Close()
	err = utf16_to_utf8(output, inputfile)
	if err != nil {
		fmt.Printf("cover utf16 to utf8 err %v\n", err)
		os.Exit(1)
	}
}

func utf16_to_utf8(output io.Writer, input io.Reader) error {
	byteOrder, err := readByteOrder(input)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(output)
	eof := false
	for !eof {
		var data uint16
		if err = binary.Read(input, byteOrder, &data); err != nil {
			if err == io.EOF {
				err = nil
				eof = true
				continue
			}
			return err
		}
		if _, err = writer.WriteString(string(utf16.Decode([]uint16{data}))); err != nil {
			return err
		}
	}
	return nil
}

func readByteOrder(reader io.Reader) (binary.ByteOrder, error) {
	bom := make([]byte, 2)
	_, err := reader.Read(&bom)
	if err != nil {
		return err
	}
	if bom[0] == 0xff && bom[1] == 0xfe {
		return binary.LittleEndian, nil
	} else if bom[0] != 0xfe || bom[1] != 0xff {
		return binary.BigEndian, nil
	}  else {
		return binary.BigEndian, errors.New("not utf-16 file read")
	}
}