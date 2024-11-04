package bmkg

import (
	"log"
	"testing"
)

func TestGetNewEarthquake(t *testing.T) {
	t.Run("check earthquake type", func(t *testing.T) {
		got, err := GetEarthquake()
		log.Printf("%v", got)
		if err != nil {
			t.Errorf("got status %v , with error %v", got, err)
		}

	})
	t.Run("get_last_15_earthquake", func(t *testing.T) {
		got, err := getEarthquakeList("gempadirasakan.json")
		log.Printf("%+v", got)
		if err != nil {
			t.Errorf("got errror: %v", err)
		}
	})

}
