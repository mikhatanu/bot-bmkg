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

type ResponseEarthquakeLast15 struct {
	Infogempa InfogempaDirasakan `json:"infogempa"`
}

type Infogempa struct {
	Gempa GempaTerbaru `json:"gempa"`
}

type InfogempaDirasakan struct {
	Gempa []GempaTerbaru `json:"gempa"`
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

const (
	baseUrl = "https://data.bmkg.go.id/DataMKG/TEWS/"
)

func GetEarthquake() (*ResponseEarthquake, error) {
	// if fileName != "autogempa.json" {
	// 	err := "filename " + fileName + " is in wrong format"
	// 	return nil, errors.New(err)
	// }
	fileName := "autogempa.json"
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

func getAllEarthquake(fileName string) (*ResponseEarthquakeLast15, error) {
	// if fileName != "gempadirasakan.json" {
	// 	err := "filename " + fileName + " is in wrong format"
	// 	return nil, errors.New(err)
	// }
	fileMap := map[string]bool{
		"gempaterkini.json":   true,
		"gempadirasakan.json": true,
	}
	if !fileMap[fileName] {
		return nil, errors.New("wrong filename")
	}
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

	returnData := ResponseEarthquakeLast15{}
	if err := json.Unmarshal(responseData, &returnData); err != nil {
		return nil, err
	}
	return &returnData, nil
}
