package bmkg

import (
	"testing"
)

func TestGetNewEarthquake(t *testing.T) {
	t.Run("check earthquake type", func(t *testing.T) {
		got, err := GetEarthquake("autogempa")
		want := got != nil

		if !want || err != nil {
			t.Errorf("got status %v want %v, with error %v", got, want, err)
		}

	})

}
