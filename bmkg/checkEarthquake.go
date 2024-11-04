package bmkg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"image/color"
	"image/png"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"

	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
)

type TextMarker struct {
	sm.MapObject
	Position   s2.LatLng
	Text       string
	TextWidth  float64
	TextHeight float64
	TipSize    float64
}

// NewTextMarker creates a new TextMarker
func NewTextMarker(pos s2.LatLng, text string) *TextMarker {
	s := new(TextMarker)
	s.Position = pos
	s.Text = text
	s.TipSize = 7.0

	d := &font.Drawer{
		Face: basicfont.Face7x13,
	}
	s.TextWidth = float64(d.MeasureString(s.Text) >> 6)
	s.TextHeight = 12.0
	return s
}

// ExtraMarginPixels returns the left, top, right, bottom pixel margin of the TextMarker object.
func (s *TextMarker) ExtraMarginPixels() (float64, float64, float64, float64) {
	w := math.Max(4.0+s.TextWidth, 2*s.TipSize)
	h := s.TipSize + s.TextHeight + 4.0
	return w * 0.5, h, w * 0.5, 0.0
}

// Bounds returns the bounding rectangle of the TextMarker object, which is just the tip position.
func (s *TextMarker) Bounds() s2.Rect {
	r := s2.EmptyRect()
	r = r.AddPoint(s.Position)
	return r
}

// Draw draws the object.
func (s *TextMarker) Draw(gc *gg.Context, trans *sm.Transformer) {
	if !sm.CanDisplay(s.Position) {
		return
	}

	w := math.Max(4.0+s.TextWidth, 2*s.TipSize)
	h := s.TextHeight + 4.0
	x, y := trans.LatLngToXY(s.Position)
	gc.ClearPath()
	gc.SetLineWidth(1)
	gc.SetLineCap(gg.LineCapRound)
	gc.SetLineJoin(gg.LineJoinRound)
	gc.LineTo(x, y)
	gc.LineTo(x-s.TipSize, y-s.TipSize)
	gc.LineTo(x-w*0.5, y-s.TipSize)
	gc.LineTo(x-w*0.5, y-s.TipSize-h)
	gc.LineTo(x+w*0.5, y-s.TipSize-h)
	gc.LineTo(x+w*0.5, y-s.TipSize)
	gc.LineTo(x+s.TipSize, y-s.TipSize)
	gc.LineTo(x, y)
	gc.SetColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.FillPreserve()
	gc.SetColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	gc.Stroke()

	gc.SetRGBA(0.0, 0.0, 0.0, 1.0)
	gc.DrawString(s.Text, x-s.TextWidth*0.5, y-s.TipSize-4.0)
}

// Create embed message embed
func CreateAllEarthquakeEmbed(req *ResponseEarthquakeLast15) []*discordgo.MessageEmbed {
	mes := []*discordgo.MessageEmbedField{}

	formatContent := ""

	for i, v := range req.Infogempa.Gempa {
		if v.FeltAt != "" {
			formatContent = fmt.Sprintf("```Tanggal: %v\nWilayah: %v\nDirasakan: %v\nMagnitude: %v\nKedalaman: %v\n```", v.Date, v.LocationInformation, v.FeltAt, v.Magnitude, v.Depth)
		} else if v.Potential != "" {
			formatContent = fmt.Sprintf("```Tanggal: %v\nWilayah: %v\nPotensi Tsunami: %v\nMagnitude: %v\nKedalaman: %v\n```", v.Date, v.LocationInformation, v.Potential, v.Magnitude, v.Depth)
		}

		mes = append(mes, &discordgo.MessageEmbedField{
			Name:  strconv.Itoa(i) + " | " + v.DateTime,
			Value: formatContent,
		})
	}

	return []*discordgo.MessageEmbed{
		{
			Title:       "Last 15 earthquake",
			Description: "",
			URL:         "https://data.bmkg.go.id/DataMKG/TEWS/gempadirasakan.json",
			Footer: &discordgo.MessageEmbedFooter{
				Text:    os.Getenv("peringatanFooter"),
				IconURL: os.Getenv("peringatanFooterURL"),
			},
			Fields: mes,
		},
	}
}

// Create static maps using openstreetmap
func createStaticMaps(req *ResponseEarthquakeLast15) (io.Reader, error) {
	ctx := sm.NewContext()
	ctx.SetSize(1920, 1080)

	var wg sync.WaitGroup
	errorChan := make(chan string)

	var iconTsunami *os.File
	iconTsunami, err := os.Open("./assets/house-tsunami.png")
	if err != nil {
		log.Panicf("Error: Error when opening file: %v", err)
	}
	defer iconTsunami.Close()
	tsunamiImage, err := png.Decode(iconTsunami)
	if err != nil {
		log.Panicf("Error: Error when opening file: %v", err)
	}

	for i, v := range req.Infogempa.Gempa {
		wg.Add(1)

		go func() {
			defer wg.Done()
			latlng := strings.Split(v.Coordinates, ",")
			lat, err := strconv.ParseFloat(latlng[0], 64)
			if err != nil {
				errorChan <- err.Error() + " in loop " + strconv.Itoa(i)
			}
			long, err := strconv.ParseFloat(latlng[1], 64)
			if err != nil {
				errorChan <- err.Error() + " in loop " + strconv.Itoa(i)
			}

			mapObj := NewTextMarker(s2.LatLngFromDegrees(lat, long), strconv.Itoa(i))
			ctx.AddObject(mapObj)

			if v.Potential != "" {
				if !strings.Contains(strings.ToLower(v.Potential), "tidak") {
					ctx.AddObject(sm.NewImageMarker(s2.LatLngFromDegrees(lat, long), tsunamiImage, 8, -8))
				}
			}
		}()

	}
	wg.Wait()
	img, err := ctx.Render()
	if err != nil {
		return nil, errors.New(err.Error())
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	select {
	case msg1 := <-errorChan:
		return nil, errors.New(msg1)
	default:
		return buf, nil
	}
}
