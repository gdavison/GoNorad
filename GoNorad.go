package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/cookiejar"
	"os"
)

type DataRequest struct {
	result []string
}

type FileStorage interface {
	GetContents()
	WriteContents()
}

type CookieFile struct{}

func (*CookieFile) GetContents(fileName string) string {
	file, err := ioutil.ReadFile(fileName)
	ErrorExit(err)

	return string(file)
}

func (*CookieFile) WriteContents(fileName string, contents string) {
	err := ioutil.WriteFile(fileName, []byte(contents), 0644)
	ErrorExit(err)
}

func main() {

	var e error
	var np2 = new(Neptune)
	config := GetConfig()
	cookieJar, e := cookiejar.New(nil)

	res, e := np2.Login(config["username"], config["password"], cookieJar)
	ErrorExit(e)

	fmt.Println(res)

	res, e = np2.GetData(config["gameNumber"], cookieJar)
	ErrorExit(e)

	fmt.Println(res)

	payloadBytes := []byte(res)
	var gameData NeptuneResponse
	e = json.Unmarshal(payloadBytes, &gameData)
	ErrorExit(e)

	fmt.Println("You are #", gameData.Report.Player_id)
}

func ErrorExit(e error) {
	if e != nil {
		fmt.Println("Error:", e)
		os.Exit(1)
	}
}
