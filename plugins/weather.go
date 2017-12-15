package plugins

import (
	"fmt"
	"github.com/blrn/gbot/bot"
	"github.com/jessfraz/weather/geocode"
	"github.com/jessfraz/weather/forecast"
	"math"
)

const (
	default_gecode_server = "https://geocode.jessfraz.com"
	default_location      = "Auburn, Il"
)

func init() {
	bot.Register("weather", func(gbot *bot.Bot) {
		gbot.SlashCommand(bot.SlashCommandConfig{Command: "/weather"}, weatherOnSlashCommand)
	})
}

func weatherOnSlashCommand(gBot *bot.Bot, req bot.SlashCommandRequest) *bot.SlashCommandResponse {
	fmt.Println("weatherOnSlashCommand")
	location := default_location
	if len(req.Text) != 0 {
		location = req.Text
	}
	resp := bot.NewSlashCommandResponse(req)
	if weatherText, err := getWeather(location); err != nil {
		fmt.Printf("Error getWeather: %+v\n", err)
		return nil
	} else {
		resp.Text = weatherText
	}
	resp.UserName = "Weather"
	return resp

}

func getDirection(degrees float64) string{
	directions := []string{
		"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW",
	}
	return directions[int(math.Mod((degrees + 11.25)/ 22.5, 16))]
}

func getWeather(location string) (string, error) {
	geo, err := geocode.Locate(location, default_gecode_server)
	if err != nil {
		return "", err
	}
	data := forecast.Request{Latitude: geo.Latitude, Longitude: geo.Longitude, Units: "auto", Exclude: []string{"hourly", "minutely"}}
	fc, err := forecast.Get(fmt.Sprintf("%s/forecast", default_gecode_server), data)
	if err != nil {
		return "", err
	}
	currentWeather := fc.Currently
	outStr := fmt.Sprintf("Current weather in %s, %s is %s\n", geo.City, geo.Region, currentWeather.Summary)
	outStr += fmt.Sprintf("The temperature is %.2f°F, but it feels like %.2f°F\n", currentWeather.Temperature, currentWeather.ApparentTemperature)
	outStr += fmt.Sprintf("The humidity is %.2f%%\n", currentWeather.Humidity * 100)
	outStr += fmt.Sprintf("The wind speed is %.2f mph %s\n", currentWeather.WindSpeed, getDirection(currentWeather.WindBearing))
	outStr += fmt.Sprintf("The pressure is %.2f mbar", currentWeather.Pressure)
	return outStr, nil
}
