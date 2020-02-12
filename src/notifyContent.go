package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func requestSMSSocket() (string, error) {
	response, err := http.Get("https://b2.51ias.com/api.php?m=Sms&a=getSmsUrl")
	if err != nil {
		fmt.Println("request sms socket err: ", err)
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	type TypeResData struct {
		Ret int    `json:"ret"`
		Msg string `json:"msg"`
		URL string `json:"url"`
	}
	var resStruct TypeResData
	err = json.Unmarshal(body, &resStruct)
	if err != nil {
		fmt.Println("requestSmsSocket Unmarshal err: ", err)
		return "", err
	}
	arr := strings.Split(resStruct.URL, "/")
	//fmt.Println(err, arr, resStruct.URL)
	return arr[2], nil
}

func notifyGameInfoHistory(timestamp, gameName, gid, state string, size uint64) string {
	name, err := os.Hostname()
	fmt.Println(name, err)
	if err != nil {
		fmt.Println("notifyGameInfoHistory Hostname err: ", err)
		return ""
	}
	addrArr, err := net.LookupIP(name)
	// addrs, err := net.InterfaceAddrs()
	// addrArr1, err := net.LookupHost(name)
	fmt.Println(len(addrArr), addrArr)
	if err != nil {
		fmt.Println("notifyGameInfoHistory LookupHost err: ", err)
		return ""
	}
	gdcip := ""
	for _, IP := range addrArr {
		if !IP.IsLoopback() && IP.To4() != nil {
			gdcip = IP.String()
			break
		}
	}
	sizeStr := strconv.FormatUint(size, 10)
	item := map[string]string{
		"gdc_ip":              gdcip,
		"gdc_hostname":        name,
		"steam_app_id":        "0",
		"gloud_game_id":       gid,
		"hotorcold":           "hot",
		"game_dir":            gameName,
		"hdisk":               "n",
		"launcher":            "m",
		"status":              state,
		"bytestodownload":     sizeStr,
		"sizeondisk":          "15892631552",
		"pre_buildid":         timestamp,
		"current_buildid":     timestamp,
		"steam_user_id":       "",
		"steam_user_password": "",
		"priority":            "",
		"update_time":         "0"}
	byteArr, err := json.Marshal(item)
	fmt.Println(string(byteArr))
	if err != nil {
		return ""
	}
	address, err := requestSMSSocket()
	fmt.Println(address, err)
	if err != nil {
		fmt.Println("notifyGameInfoHistory requestSMSSocket err: ", err)
		return ""
	}
	response, err := http.Get(fmt.Sprintf("http://%s/update_hoc_gdc?op_token=gloudhotorcoldtoken&content=%s", address, string(byteArr)))
	if err != nil {
		fmt.Println("notifyGameInfoHistory http.Get err: ", err)
		return ""
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("notifyGameInfoHistory ioutil.ReadAll err: ", err)
		return ""
	}
	return string(body)
}
