package main

import (
	//"errors"
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"strconv"
	//"bytes"
	"unicode/utf16"
	//"net/http"
	//"io/ioutil"
	"encoding/json"
)
const IP_PORT = "121.40.111.19:443"

func utf16ToString(b []byte, bom int) string {
	if len(b) >= 2 {
		switch n := uint16(b[0])<<8 | uint16(b[1]); n {
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

type TContent struct {
	id        string
	name      string
	stat      string
	startTime int64
	endTime   int64
	sizeStr   string
}

func getIDNameArrAndTime(str string, year int, month int, day int) ([]string, int64) {
	timeStr := str[0:8]
	timeStrArr := strings.Split(timeStr, ":")
	var h, _ = strconv.Atoi(timeStrArr[0])
	var m, _ = strconv.Atoi(timeStrArr[1])
	var s, _ = strconv.Atoi(timeStrArr[2])
	var tTime = time.Date(year, time.Month(month), day, h, m, s, 0, time.UTC)
	var unixTime = tTime.Unix()
	endIdx := strings.Index(str, "]")
	idNameStr := str[10:endIdx]
	strArr := strings.Split(idNameStr, ":")
	return strArr, unixTime
}

func parseFile(idInfoData map[string]TContent) {
	f, err := os.Open("log")
	fmt.Println(f, err)
	if err != nil {
		return
	}
	defer f.Close()
	arr, err := f.Readdir(0)
	if err != nil {
		return
	}
	var fileNameArr []string
	fmt.Println(arr, len(arr), err)
	for _, info := range arr {
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
	var year, _ = strconv.Atoi(fileName[8:12])
	var month, _ = strconv.Atoi(fileName[12:14])
	var day, _ = strconv.Atoi(fileName[14:16])
	fmt.Println(time.Now(), year, month, day)
	br := bufio.NewReader(fd)
	for {
		byteArr, _, ed := br.ReadLine()
		if ed == io.EOF {
			break
		}
		str := utf16ToString(byteArr, 1)
		idx := strings.Index(str, "start download")
		if idx > -1 {
			strArr, unixTime := getIDNameArrAndTime(str, year, month, day)
			content, ok := idInfoData[strArr[0]]
			if ok {
				if (content.stat == "end" && content.endTime < unixTime){
					content.stat = "start"
					content.startTime = unixTime
				}else{
					fmt.Println("游戏无新更新 ", strArr)
				}
			} else {
				fmt.Println("22222222")
				idInfoData[strArr[0]] = TContent{strArr[0], strArr[1], "start", unixTime, 0, "0M"}
			}
			
		}
		idx = strings.Index(str, "download complete")
		if idx > -1 {
			strArr, unixTime := getIDNameArrAndTime(str, year, month, day)
			content, ok := idInfoData[strArr[0]]
			if ok {
				if (content.stat == "start" && content.endTime < unixTime){
					content.stat = "end"
					content.endTime = unixTime
					content.sizeStr = str[idx+32:]
				}
			} else {
				fmt.Println("没有start，直接end", strArr)
			}
		}
	}
}


func main(){
	var idInfoData map[string]TContent
	idInfoData = make(map[string]TContent)
	item := map[string]string{
		"gdc_ip": "",
		"gdc_hostname": "",
		"steam_app_id": "",
		"gloud_game_id": "",
		"hotorcold": "0",
		"game_dir": "0",
		"hdisk": "0",
		"cdisk": "0",
		"status": "0",
		"bytestodownload": "0",
		"sizeondisk": "0",
		"pre_buildid": "0",
		"current_buildid": "0",
		"steam_user_id": "",
		"steam_user_password": "",
		"priority": "",
		"update_time": "0"}
	byteArr, err := json.Marshal(item)
	fmt.Println(string(byteArr), err)
	// response, err := http.Get(fmt.Sprintf("http://%s/update_hoc_gdc?op_token=gloudhotorcoldtoken&content=%s", IP_PORT, string(byteArr)))
	// defer response.Body.Close()
	// body, err := ioutil.ReadAll(response.Body)
	// fmt.Println(body)
	for{
		parseFile(idInfoData)
		time.Sleep(10*time.Second)
	}
}