package openweather

type CoordinatesResponse struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	}
}

type Weather struct {
	Temp float64 `json:"temp"`
}
