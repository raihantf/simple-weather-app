package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherDataResponse struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

type weatherDataRequest struct {
	City string `json:"city"`
}

func loadApiConfig(fileName string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return apiConfigData{}, err
	}

	var data apiConfigData

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return apiConfigData{}, err
	}

	return data, nil
}

func queryData(city string) (weatherDataResponse, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return weatherDataResponse{}, err
	}

	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + apiConfig.OpenWeatherMapApiKey)
	if err != nil {
		return weatherDataResponse{}, err
	}

	defer resp.Body.Close()

	var data weatherDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return weatherDataResponse{}, err
	}

	return data, nil
}

func main() {
	app := fiber.New()

	app.Get("/weather/", func(ctx *fiber.Ctx) error {
		params := ctx.Query("city")

		var request weatherDataRequest

		request.City = params

		data, err := queryData(request.City)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		return ctx.JSON(&fiber.Map{
			"message": data,
		})
	})

	app.Listen(":8080")
}
