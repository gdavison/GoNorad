package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http/cookiejar"
	"os"
	"strconv"
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

	//fmt.Println(res)

	var gameData NeptuneResponse
	payloadBytes := []byte(res)
	e = json.Unmarshal(payloadBytes, &gameData)
	ErrorExit(e)

	fmt.Println("You are #", gameData.Report.Player_id)

	var myStars []StarType
	for _, star := range gameData.Report.Stars {
		if star.PlayerId == gameData.Report.Player_id {
			myStars = append(myStars, star)
		}
	}

	fmt.Println("Found ", len(myStars), " of your stars")

	star1 := myStars[0]
	fmt.Println("Star 1: ", star1.Name, " (x:", star1.X, ",y:", star1.Y, ")")
	star2 := myStars[1]
	fmt.Println("Star 2: ", star2.Name, " (x:", star2.X, ",y:", star2.Y, ")")

	d := starDistance(star1, star2)
	fmt.Println("raw distance: ", d)

	lightyears := d * 8
	fmt.Println("light years: ", lightyears)

	myPlayer := gameData.Report.Players[strconv.Itoa(gameData.Report.Player_id)]
	myHyperspaceTech := myPlayer.Tech["propulsion"]
	fmt.Println("Hyperspace tech value: ", myHyperspaceTech.Value)

	hyperspaceLevel := math.Ceil(lightyears - 3)
	rawHyperLevel := math.Ceil((d - 0.375) / 0.125)

	fmt.Println("Tech level needed: ", hyperspaceLevel)
	fmt.Println("Tech level needed (raw): ", rawHyperLevel)
}

func starDistance(star1, star2 StarType) float64 {
	dX := star1.X - star2.X
	dY := star1.Y - star2.Y
	fmt.Println("dX: ", dX, ", dY: ", dY)

	return math.Sqrt(dX*dX + dY*dY)
}

func ErrorExit(e error) {
	if e != nil {
		fmt.Println("Error:", e)
		os.Exit(1)
	}
}
