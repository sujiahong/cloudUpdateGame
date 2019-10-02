package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

const IP_PORT = "121.40.111.19:443"

func notifyGameUpdate() string {
	name, err := os.Hostname()
	fmt.Println(name, err)
	if err != nil {
		return ""
	}
	addrArr, err := net.LookupHost(name)
	fmt.Println(addrArr, err)
	if err != nil {
		return ""
	}
	item := map[string]string{
		"gdc_ip":              addrArr[1],
		"gdc_hostname":        name,
		"steam_app_id":        "",
		"gloud_game_id":       "",
		"hotorcold":           "0",
		"game_dir":            "0",
		"hdisk":               "0",
		"cdisk":               "0",
		"status":              "0",
		"bytestodownload":     "0",
		"sizeondisk":          "0",
		"pre_buildid":         "0",
		"current_buildid":     "0",
		"steam_user_id":       "",
		"steam_user_password": "",
		"priority":            "",
		"update_time":         "0"}
	byteArr, err := json.Marshal(item)
	fmt.Println(string(byteArr))
	if err != nil {
		return ""
	}
	response, err := http.Get(fmt.Sprintf("http://%s/update_hoc_gdc?op_token=gloudhotorcoldtoken&content=%s", IP_PORT, string(byteArr)))
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	return string(body)
}
