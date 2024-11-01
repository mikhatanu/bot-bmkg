package bmkg

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type ResponseEarthquake struct {
	Infogempa Infogempa `json:"Infogempa"`
}

type Infogempa struct {
	Gempa GempaTerbaru `json:"gempa"`
}

type GempaTerbaru struct {
	Date                string `json:"Tanggal"`
	Time                string `json:"Jam"`
	DateTime            string `json:"DateTime"`
	Coordinates         string `json:"Coordinates"`
	Latitude            string `json:"Lintang"`
	Langitude           string `json:"Bujur"`
	Magnitude           string `json:"Magnitude"`
	Depth               string `json:"Kedalaman"`
	LocationInformation string `json:"wilayah"`
	Potential           string `json:"Potensi"`
	FeltAt              string `json:"Dirasakan"`
	Shakemap            string `json:"Shakemap"`
}

func GetEarthquake(fileName string) (*ResponseEarthquake, error) {
	bmkgFileName := map[string]bool{
		"autogempa.json":      true,
		"gempaterkini.json":   false,
		"gempadirasakan.json": false,
	}
	if !bmkgFileName[fileName] {
		err := "filename " + fileName + " not found"
		return nil, errors.New(err)
	}

	baseUrl := "https://data.bmkg.go.id/DataMKG/TEWS/"
	requestUrl := baseUrl + fileName

	response, err := http.Get(requestUrl)
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
	returnData := ResponseEarthquake{}
	if err := json.Unmarshal(responseData, &returnData); err != nil {
		return nil, err
	}
	return &returnData, nil
}
