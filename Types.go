package main

type ReportType struct {
	Player_id int `json:"player_uid"`
	fleets    map[string]FleetType
	players   map[string]PlayerType
	Stars     map[string]StarType `json:"stars"`
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
	name     string
	economy  int
	industry int
	science  int
	stars    int
}

type StarType struct {
	id       int
	name     string
	economy  int
	industry int
	science  int
}
