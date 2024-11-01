package bmkg

import (
	"testing"
)

func TestGetNewEarthquake(t *testing.T) {
	t.Run("check earthquake type", func(t *testing.T) {
		got, err := GetEarthquake("autogempa.json")
		// log.Printf("%v", got)
		if err != nil {
			t.Errorf("got status %v , with error %v", got, err)
		}

	})

}
