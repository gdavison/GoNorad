package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
	"github.com/gonum/graph/search"
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

	destination := flag.String("destination", "", "the name of the destination star")
	source := flag.String("source", "", "the name of the source star")
	flag.Parse()

	fmt.Println("looking for ", *destination)
	if *source != "" {
		fmt.Println("starting from ", *source)
	}

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

	var gameData NeptuneResponse
	payloadBytes := []byte(res)
	e = json.Unmarshal(payloadBytes, &gameData)
	ErrorExit(e)

	fmt.Println("You are player #", gameData.Report.Player_id)

	var allStars []StarType
	var myStars []StarType
	for _, star := range gameData.Report.Stars {
		if star.PlayerId == gameData.Report.Player_id {
			myStars = append(myStars, star)
		}
		allStars = append(allStars, star)
	}

	for _, star := range myStars {
		if star.HasGate() {
			fmt.Println("Star ", star.Name, " has a Gate")
		}
	}

	fmt.Println("Found ", len(myStars), " of your stars")
	fmt.Println("Found ", len(allStars), " total stars")

	//	star1 := myStars[0]
	//	fmt.Println("Star 1: ", star1.Name, " (x:", star1.X, ",y:", star1.Y, ")")
	//	star2 := myStars[1]
	//	fmt.Println("Star 2: ", star2.Name, " (x:", star2.X, ",y:", star2.Y, ")")
	//
	//	lightyears := starDistanceInLightyears(star1, star2)
	//	fmt.Println("light years: ", lightyears)

	starDistanceGraph := concrete.NewGraph()
	for i, star1 := range allStars {
		for _, star2 := range allStars[i+1:] {
			distance := starDistanceInLightyears(star1, star2)
			//fmt.Println(star1.Name, "(", star1.Id, ") - ", star2.Name, "(", star2.Id, "): ", distance)
			edge := concrete.Edge{
				H: concrete.Node(star1.Id),
				T: concrete.Node(star2.Id)}
			starDistanceGraph.AddUndirectedEdge(edge, distance)
		}
	}

	myPlayer := gameData.Report.Players[strconv.Itoa(gameData.Report.Player_id)]
	myHyperspaceTech := myPlayer.Tech["propulsion"]
	fmt.Println("My Hyperspace tech value: ", myHyperspaceTech.Value)

	reachableStarsGraph := concrete.NewGraph()
	for _, candidateEdge := range starDistanceGraph.EdgeList() {
		cost := starDistanceGraph.Cost(candidateEdge)
		if distanceIsReachable(cost, myHyperspaceTech.Level) {
			reachableStarsGraph.AddUndirectedEdge(candidateEdge, cost)
		}
	}

	testEdge := reachableStarsGraph.EdgeList()[0]
	fmt.Println("test: from ", gameData.Report.Stars[strconv.Itoa(testEdge.Head().ID())].Name,
		"to ", gameData.Report.Stars[strconv.Itoa(testEdge.Tail().ID())].Name,
		": ", reachableStarsGraph.Cost(testEdge))
	fmt.Println("There are ", len(starDistanceGraph.EdgeList()), " total edges")
	fmt.Println("There are ", len(reachableStarsGraph.EdgeList()), " reachable edges")

	var productionStars []StarType
	for _, star := range myStars {
		if star.Industry > 0 {
			productionStars = append(productionStars, star)
		} else {
			fmt.Println("No production at ", star.Name)
		}
	}

	foundDestinationStar := false
	var destinationStar StarType
	for _, star := range allStars {
		if star.Name == *destination {
			destinationStar = star
			foundDestinationStar = true
			continue
		}
	}

	if !foundDestinationStar {
		fmt.Println("Could not find ", *destination)
		os.Exit(0)
	}

	destinationNode := concrete.Node(destinationStar.Id)
	paths, _ := search.Dijkstra(destinationNode, reachableStarsGraph, func(edge graph.Edge) float64 {
		return reachableStarsGraph.Cost(edge)
	})

	var testStar StarType
	if *source != "" {
		foundSourceStar := false
		var sourceStar StarType
		for _, star := range allStars {
			if star.Name == *source {
				sourceStar = star
				foundSourceStar = true
				continue
			}
		}
		if !foundSourceStar {
			fmt.Println("Could not find ", *source)
			os.Exit(0)
		}
		testStar = sourceStar
	} else {
		testStar = productionStars[0]
	}
	fmt.Println("Path to ", testStar.Name, " (", testStar.Id, ")")
	testPath := paths[testStar.Id]
	for _, starId := range testPath {
		fmt.Println(gameData.Report.Stars[strconv.Itoa(starId.ID())].Name)
	}
}

func starDistanceInLightyears(star1, star2 StarType) float64 {
	dX := star1.X - star2.X
	dY := star1.Y - star2.Y

	return math.Sqrt(dX*dX+dY*dY) * 8
}

func distanceIsReachable(lightyears float64, techLevel int) bool {
	requiredTechLevel := int(math.Max(math.Ceil(lightyears-3), 1))
	//fmt.Println("REQUIRED: ", requiredTechLevel)

	return techLevel >= requiredTechLevel
}

func ErrorExit(e error) {
	if e != nil {
		fmt.Println("Error:", e)
		os.Exit(1)
	}
}
