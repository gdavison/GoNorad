package main

type ReportType struct {
	Player_id int `json:"player_uid"`
	fleets    map[string]FleetType
	Players   map[string]PlayerType `json:"players"`
	Stars     map[string]StarType   `json:"stars"`
}

type NeptuneResponse struct {
	event  string
	order  string
	error  string
	Report ReportType `json:"report"`
}

type FleetType struct {
	ships          int
	owner_id       int
	destination_id int
	name           string
}

type PlayerType struct {
	id       int
	Name     string `json:"alias"`
	economy  int
	industry int
	science  int
	stars    int
	Tech     map[string]TechType `json:"tech"`
}

type TechType struct {
	Value float64 `json:"value"`
	Level int     `json:"level"`
}

type StarType struct {
	Id       int
	Name     string  `json:"n"`
	PlayerId int     `json:"puid"`
	X        float64 `json:"x,string"`
	Y        float64 `json:"y,string"`
	economy  int
	industry int
	science  int
}
