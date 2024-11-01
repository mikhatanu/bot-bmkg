package bmkg

import (
	"testing"
)

func TestGetWeatherForecast(t *testing.T) {
	t.Run("checkWeatherForecastAPI", func(t *testing.T) {
		got, err := GetWeatherForecast("31.01.02.1002")
		// log.Printf("%+v", got)

		if err != nil {
			t.Errorf("got status %v , with error %v", got, err)
		}

	})

}
