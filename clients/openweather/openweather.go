package openweather

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OpenWeatherClient struct {
	apiKey string
}

func New(apiKey string) *OpenWeatherClient {
	return &OpenWeatherClient{apiKey: apiKey}
}

func (o OpenWeatherClient) Coordinates(city string) (Coordinates, error) {
	url := "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s"

	resp, err := http.Get(fmt.Sprintf(url, city, o.apiKey))
	if err != nil {
		log.Printf("[ERROR] Ошибка сетевого запроса при получении координат для города '%s': %v\n", city, err)
		return Coordinates{}, fmt.Errorf("сетевая ошибка при запросе координат для города '%s': %w", city, err)
	}

	if resp.StatusCode != http.StatusOK {
		return Coordinates{}, fmt.Errorf("API вернул ошибку при получении координат для города '%s': статус код %d", city, resp.StatusCode)
	}

	var coordinatesResponse []CoordinatesResponse
	err = json.NewDecoder(resp.Body).Decode(&coordinatesResponse)
	if err != nil {
		log.Printf("[ERROR] Ошибка при парсинге ответа координат для города '%s': %v\n", city, err)
		return Coordinates{}, fmt.Errorf("ошибка парсинга координат для города '%s': %w", city, err)
	}
	if len(coordinatesResponse) == 0 {
		return Coordinates{}, fmt.Errorf("город '%s' не найден в базе данных OpenWeather", city)
	}

	log.Printf("[INFO] Получены координаты для города '%s': lat=%.2f, lon=%.2f\n", city, coordinatesResponse[0].Lat, coordinatesResponse[0].Lon)
	return Coordinates{
		Lat: coordinatesResponse[0].Lat,
		Lon: coordinatesResponse[0].Lon,
	}, nil
}

func (o OpenWeatherClient) Weather(lat, lon float64) (Weather, error) {
	url := "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric"
	resp, err := http.Get(fmt.Sprintf(url, lat, lon, o.apiKey))
	if err != nil {
		log.Printf("[ERROR] Ошибка сетевого запроса при получении погоды для координат (lat: %.2f, lon: %.2f): %v\n", lat, lon, err)
		return Weather{}, fmt.Errorf("сетевая ошибка при запросе погоды для координат (lat: %.2f, lon: %.2f): %w", lat, lon, err)
	}

	if resp.StatusCode != http.StatusOK {
		return Weather{}, fmt.Errorf("API вернул ошибку при получении погоды для координат (lat: %.2f, lon: %.2f): статус код %d", lat, lon, resp.StatusCode)
	}

	var weatherResponse WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		log.Printf("[ERROR] Ошибка при парсинге ответа о погоде для координат (lat: %.2f, lon: %.2f): %v\n", lat, lon, err)
		return Weather{}, fmt.Errorf("ошибка парсинга данных о погоде для координат (lat: %.2f, lon: %.2f): %w", lat, lon, err)
	}

	log.Printf("[INFO] Получена информация о погоде для координат (lat: %.2f, lon: %.2f): температура %.2f°C\n", lat, lon, weatherResponse.Main.Temp)
	return Weather{
		Temp: weatherResponse.Main.Temp,
	}, nil

}
