package main

type WeatherResponse struct {
	Temperature float64 `json:"temp"`
	FeelsLike   float64 `json:"feelslike"`
	Humidity    float64 `json:"humidity"`
	Conditions  string  `json:"conditions"`
	UVIndex     float64 `json:"uvindex"`
	WindSpeed   float64 `json:"windspeed"`
	Sunrise     string  `json:"sunrise"`
	Sunset      string  `json:"sunset"`
}

type APIResponse struct {
	CurrentConditions WeatherResponse `json:"currentConditions"`
}
