package util

import (
	"testing"
)

func TestGetWeatherForecast(t *testing.T) {
	t.Run("check load csv", func(t *testing.T) {
		got := loadAdmFile()
		// log.Printf("%+v", got)
		// for i := 0; i < 10; i++ {
		// 	log.Println(got[i])
		// 	for _, col := range got[i] {
		// 		log.Println(col)
		// 	}
		// }
		if !got {
			t.Errorf("got status %v", got)
		}

	})

}
