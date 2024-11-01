package util

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

var admLocationMap = make(map[string][]string)
var admCodeMap = make(map[string]string)

// Load from file
func loadAdmFile() bool {
	// Open and read file
	myFile, err := os.Open("./util/adm.csv")
	if err != nil {
		log.Panicf("Error: Error when opening file: %v", err)
	}
	defer myFile.Close()

	// Read the CSV data
	reader := csv.NewReader(myFile)
	reader.FieldsPerRecord = 0
	data, err := reader.ReadAll()
	if err != nil {
		log.Panicf("Error: Error when reading csv: %v", err)
	}

	// Save csv to global variable
	for _, row := range data {
		admCodeMap[row[0]] = row[1]
		locationLower := strings.ToLower(row[1])
		admLocationMap[locationLower] = append(admLocationMap[locationLower], row[0])
	}
	return true
}

func init() {
	// This command will always run at the working directory of main.go
	loadAdmFile()
}

// get location from adm code
func GetLocationFromAdmCode(kodeWilayah string) string {
	return admCodeMap[kodeWilayah]
}

// get administration code, with array because some place have same name
func GetAdmCodeFromLocation(location string) []string {
	var adm4Code []string
	for _, v := range admLocationMap[location] {
		if len(strings.Split(v, ".")) == 4 {
			adm4Code = append(adm4Code, v)
		}
	}

	return adm4Code
}

// Get full name in array, e.g. ["BALI", "KAB. BULELENG", "Sukasada", "Pegadungan"]
func GetFullAdministrationLocationName(kodeWilayah string) []string {
	splitKode := strings.Split(kodeWilayah, ".")

	namaLengkap := make([]string, len(splitKode))
	var v string
	for i := range splitKode {
		v = strings.Join(splitKode[0:i+1], ".")

		namaLengkap[i] = GetLocationFromAdmCode(v)

	}
	return namaLengkap
}
