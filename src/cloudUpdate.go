package main

import (
	//"errors"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	//"bytes"
	"unicode/utf16"
	//"net/http"
)

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

// 存储更新content.
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
	now := time.Now()
	nowYear, nowMonth, nowDay := now.Date()
	formatStr := "%d%d%d"
	if nowMonth < 10 && nowDay < 10 {
		formatStr = "%d0%d0%d"
	} else if nowMonth < 10 {
		formatStr = "%d0%d%d"
	} else if nowDay < 10 {
		formatStr = "%d%d0%d"
	}
	fileName := "kupdate-" + fmt.Sprintf(formatStr, nowYear, nowMonth, nowDay) + ".log"
	fd, err := os.Open("C:\\kserver\\log\\" + fileName)
	fmt.Println(fd, err)
	if err != nil {
		return
	}
	defer fd.Close()
	var year, _ = strconv.Atoi(fileName[8:12])
	var month, _ = strconv.Atoi(fileName[12:14])
	var day, _ = strconv.Atoi(fileName[14:16])
	fmt.Println(time.Now(), year, month, day)
	br := bufio.NewReader(fd)
	for {
		byteArr, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		str := utf16ToString(byteArr, 1)
		idx := strings.Index(str, "start download")
		if idx > -1 {
			strArr, unixTime := getIDNameArrAndTime(str, year, month, day)
			content, ok := idInfoData[strArr[0]]
			if ok {
				if content.stat == "end" && content.endTime < unixTime {
					content.stat = "start"
					content.startTime = unixTime
					//notifyGameUpdate()
				} else {
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
				if content.stat == "start" && content.endTime < unixTime {
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

func isHaveReportSMS(bid, name, stat string) bool {
	type TypeGameStatData struct {
		Name  string
		Stat0 string
		Stat1 string
	}
	var bidStatDataMap map[string]TypeGameStatData
	bidStatDataMap = make(map[string]TypeGameStatData)
	f, err := os.OpenFile("./bidStat.json", os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("isHaveReportSMS os.OpenFile err: ", err)
		return false
	}
	defer f.Close()
	byteArr, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("isHaveReportSMS ioutil.ReadAll err: ", err)
		return false
	}
	json.Unmarshal(byteArr, &bidStatDataMap)
	//fmt.Println("---   ", bidStatDataMap)
	statData, ok := bidStatDataMap[bid]
	if ok && statData.Name == name {
		if stat == "0" && statData.Stat0 == stat {
			return true
		} else if stat == "1" && statData.Stat1 == stat {
			return true
		}else if stat == "1" && statData.Stat1 != stat{
			bidStatDataMap[bid] = TypeGameStatData{name, "0", stat}
			jsonByteArr, err := json.Marshal(bidStatDataMap)
			f.Seek(0, 0)
			w := bufio.NewWriter(f)
			n, err := w.Write(jsonByteArr)
			fmt.Println("++++  ", n, err, jsonByteArr, len(jsonByteArr))
			if err != nil {
				fmt.Println("isHaveReportSMS fd.Write err: ", err)
				return false
			}
			w.Flush()
			return false
		}
		return false
	} else {
		if stat == "0" {
			bidStatDataMap[bid] = TypeGameStatData{name, stat, stat}
		} else {
			bidStatDataMap[bid] = TypeGameStatData{name, "0", stat}
		}
		jsonByteArr, err := json.Marshal(bidStatDataMap)
		f.Seek(0, 0)
		w := bufio.NewWriter(f)
		n, err := w.Write(jsonByteArr)
		fmt.Println("++++  ", n, err, len(jsonByteArr))
		if err != nil {
			fmt.Println("isHaveReportSMS fd.Write err: ", err)
			return false
		}
		w.Flush()
		return false
	}
}

func isGameUpdating(name, gid *string) bool {
	now := time.Now()
	nowYear, nowMonth, nowDay := now.Date()
	formatStr := "%d%d%d"
	if nowMonth < 10 && nowDay < 10 {
		formatStr = "%d0%d0%d"
	} else if nowMonth < 10 {
		formatStr = "%d0%d%d"
	} else if nowDay < 10 {
		formatStr = "%d%d0%d"
	}
	fileName := "kupdate-" + fmt.Sprintf(formatStr, nowYear, nowMonth, nowDay) + ".log"
	fd, err := os.Open("C:\\kserver\\log\\" + fileName)
	if err != nil {
		fmt.Println("isGameUpdating  os.Open  err: ", err)
		return false
	}
	defer fd.Close()
	// var year, _ = strconv.Atoi(fileName[8:12])
	// var month, _ = strconv.Atoi(fileName[12:14])
	// var day, _ = strconv.Atoi(fileName[14:16])
	// fmt.Println("222222222     ", time.Now(), year, month, day)
	br := bufio.NewReader(fd)
	var updateStat = false
	var curBuildId = ""
	for {
		byteArr, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		str := utf16ToString(byteArr, 1)
		idx := strings.Index(str, *name)
		if idx > -1 {
			arr := strings.Split(str, ":")
			l := len(arr)
			fmt.Println("#@#@#   ", arr, len(arr[0]), len(arr[1]), len(arr[2][:2]), len(arr[2]))
			if strings.Index(arr[l-1], "start download") > -1 {
				fmt.Println(strings.Index(arr[l-1], "\r\n"))
				h, _ := strconv.Atoi(arr[0])
				m, _ := strconv.Atoi(arr[1])
				s, _ := strconv.Atoi(arr[2][:2])
				t := time.Date(nowYear, nowMonth, nowDay, h, m, s, 0, time.Local)
				curBuildId = strconv.FormatInt(t.Unix(), 10)
				if !isHaveReportSMS(curBuildId, *name, "0") {
					notifyGameInfoHistory(curBuildId, *name, *gid, "0", 0)
				}
				updateStat = true
			} else if arr[l-2] == "download complete, update total" {
				sizeStr := arr[l-1][0:len(arr[l-1])-2]
				size, err := strconv.ParseFloat(sizeStr, 64)
				if err != nil{
					fmt.Println("isGameUpdating strconv.ParseFloat err:", err)
					return false
				}
				if !isHaveReportSMS(curBuildId, *name, "1") {
					notifyGameInfoHistory(curBuildId, *name, *gid, "1", uint64(size*1024*1024))
				}
				updateStat = false
			}
		}
	}
	return updateStat
}

func main() {
	gameName := flag.String("name", "英雄联盟", "指定查找的游戏！")
	gloudGameID := flag.String("gid", "0", "gloud game id 格来云游戏gameid")
	flag.Parse()
	fmt.Println("dddd   ", *gameName, *gloudGameID, flag.Args())
	if isGameUpdating(gameName, gloudGameID) {
		os.Exit(11)
	}
	os.Exit(0)
}
