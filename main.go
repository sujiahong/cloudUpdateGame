package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

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
		str, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		fmt.Println(str)
		idx := strings.Index(str, "start download")
		if idx > -1 {

		}
		idx = strings.Index(str, "download complete")
		if idx > -1 {

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
