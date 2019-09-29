package main

import (
	//"errors"
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	//"bytes"
	"unicode/utf16"
)

func utf16ToString(b []byte, bom int) (string){
	if len(b) >= 2 {
		switch n:=uint16(b[0])<<8 | uint16(b[1]); n {
		case 0xfffe:
			fallthrough
		case 0xfeff:
			b = b[2:]
			break
		default:
			b = b[1:]
		}
	}
	utf16Arr := make([]uint16, len(b)/2)
	for i := range utf16Arr {
		utf16Arr[i] = uint16(b[2*i+bom&1])<<8 | uint16(b[2*i+(bom+1)&1])
	}
	return string(utf16.Decode(utf16Arr))
}

func main() {
	f, err := os.Open("log")
	fmt.Println(f, err)
	if err != nil {
		return
	}
	defer f.Close()
	arr, e := f.Readdir(0)
	if e != nil {
		return
	}
	var fileNameArr []string
	fmt.Println(arr, len(arr), e)
	for _, info := range arr {
		fmt.Println(info.Name(), info.IsDir())
		strArr := strings.Split(info.Name(), "-")
		if strArr[0] == "kupdate" {
			fileNameArr = append(fileNameArr, info.Name())
		}
	}
	fmt.Println(fileNameArr)
	var fileName = fileNameArr[len(fileNameArr)-1]
	fd, er := os.Open("log/" + fileName)
	fmt.Println(fd, er)
	if er != nil {
		return
	}
	defer fd.Close()
	br := bufio.NewReader(fd)
	for {
		byteArr, _, ed := br.ReadLine()
		if ed == io.EOF {
			break
		}
		str := utf16ToString(byteArr, 1)
		fmt.Println(str, ed)
		idx := strings.Index(str, "start download")
		if idx > -1 {
			timeStr := str[0:8]
			fmt.Println(timeStr)
		}
		idx = strings.Index(str, "download complete")
		if idx > -1 {
			timeStr := str[:8]
			fmt.Println(timeStr)
		}
	}
	// f, err = os.Open("golang.org/" + arr[0].Name())
	// fmt.Println(f, err)
	// arr, e = f.Readdir(0)
	// fmt.Println(arr, len(arr), e)
	// for _, info := range arr {
	// 	fmt.Println(info.Name(), info.IsDir())
	// }
}
