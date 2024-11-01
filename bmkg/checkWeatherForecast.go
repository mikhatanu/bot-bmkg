package bmkg

import (
	"bot-bmkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ResponseWeather struct {
	Lokasi Lokasi `json:"lokasi"`
	Data   []Data `json:"data"`
}

type Lokasi struct {
	// Provinsi  string      `json:"provinsi"`
	// Kota      string      `json:"kota"`
	// Kecamatan string      `json:"kecamatan"`
	// Desa      string      `json:"desa"`
	AdministrationCode string      `json:"adm4"`
	Longitude          json.Number `json:"lon"`
	Latitude           json.Number `json:"lat"`
	Timezone           string      `json:"timezone"`
}

// first array is for day, second array is for time.
// To get current weather forecast to the next 3 hour is, [0][1]. To get today's weather forecast completely: first array is [0], then loop through the second array
type Data struct {
	Cuaca [][]Cuaca
}

// tcc             BadExpr
// tp              BadExpr
// Weather         BadExpr
// wd_deg          BadExpr
// wd              "SE"
// wd_to           "NW"
// vs_text     "> 10 km"
// utc_datetime   "2024-10-31 00:00:00"
type Cuaca struct {
	Datetime      string  `json:"datetime"`
	Temperature   int     `json:"t"`
	WeatherDesc   string  `json:"weather_desc"`
	WeatherDescEn string  `json:"weather_desc_en"`
	Windspeed     float64 `json:"ws"`
	Humidity      int     `json:"hu"`
	Visibility    int     `json:"vs"`
	TimeIndex     string  `json:"time_index"`
	AnalysisDate  string  `json:"analysis_date"`
	Image         string  `json:"image"`
	LocalDatetime string  `json:"local_datetime"`
}

var (
	scheme           = "https"
	apiHost          = "api.bmkg.go.id"
	apiPath          = "/publik/prakiraan-cuaca"
	administrationId = "adm4"
)

func getWeatherForecastAPIURL(kodeWilayah string) url.URL {
	q := administrationId + "=" + kodeWilayah
	urlBuilder := url.URL{
		Scheme:   scheme,
		Host:     apiHost,
		Path:     apiPath,
		RawQuery: q,
	}
	return urlBuilder
}

func GetWeatherForecast(KodeWilayah string) (*ResponseWeather, error) {
	urlBuilder := getWeatherForecastAPIURL(KodeWilayah)

	response, err := http.Get(urlBuilder.String())
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 404 {
		return nil, errors.New("404 Not Found")
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	returnData := ResponseWeather{}
	if err := json.Unmarshal(responseData, &returnData); err != nil {
		return nil, err
	}
	return &returnData, nil
}

func createWeatherEmbedResponse(req *ResponseWeather) []*discordgo.MessageEmbed {
	mes := []*discordgo.MessageEmbedField{}
	var timeWeather time.Time
	// set Field based on array
	for i, v := range req.Data[0].Cuaca[0] {
		// In case BMKG gets drunk and send unlimited number of array
		if i > 9 {
			break
		}

		value := fmt.Sprintf(os.Getenv("weatherPlaceholder"), v.WeatherDesc, v.Temperature, v.Humidity, v.Visibility, v.Windspeed, v.AnalysisDate)

		tF, _ := time.Parse("2006-01-02 15:04:05", v.LocalDatetime)
		tFormat := fmt.Sprintf("%02d:%02d", tF.Hour(), tF.Minute())
		if i == 0 {
			timeWeather = tF
		}
		mes = append(mes, &discordgo.MessageEmbedField{
			Name:  tFormat,
			Value: value,
		})
	}

	// Set embed description
	fullAdmName := strings.Join(util.GetFullAdministrationLocationName(req.Lokasi.AdministrationCode), ", ")
	description := "[" + fullAdmName + "](https://kodewilayah.id/" + req.Lokasi.AdministrationCode + ")"

	// Format date for tittle
	timeFormatted := fmt.Sprintf("%v %v %v", timeWeather.Day(), timeWeather.Month(), timeWeather.Year())

	// Get url of api
	url := getWeatherForecastAPIURL(req.Lokasi.AdministrationCode)
	return []*discordgo.MessageEmbed{
		{
			Title:       "Weather Forecast | " + timeFormatted,
			Description: description,
			URL:         url.String(),
			Footer: &discordgo.MessageEmbedFooter{
				Text:    os.Getenv("peringatanFooter"),
				IconURL: os.Getenv("peringatanFooterURL"),
			},
			Fields: mes,
		},
	}
}
