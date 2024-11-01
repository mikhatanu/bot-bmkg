package bmkg

import (
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type InMemoryStore struct {
	LatestEarthquake time.Time
}

type ChannelIDMemoryStore struct {
	ChannelID   string
	ChannelName string
}

func NewEarthquakeHandler(x *InMemoryStore, s *discordgo.Session, c map[string]*ChannelIDMemoryStore) bool {
	// Get earthquake with file name autogempa.json which gets latest earthquake from bmkg
	curr, err := GetEarthquake("autogempa.json")
	if err != nil {
		log.Printf("Error: %v", err)
		return false
	}

	// convert return datetime to time.time
	convertToDate, err := time.Parse(time.RFC3339, curr.Infogempa.Gempa.DateTime)
	if err != nil {
		log.Printf("Error: Date '%v' from BMKG is not valid date", curr.Infogempa.Gempa.DateTime)
		return false
	}

	// check if converted date is after the latest earthquake, so bot only runs on new earthquake
	if convertToDate.After(x.LatestEarthquake) {
		x.LatestEarthquake = convertToDate
		log.Printf("Info: InMemoryStore updated to %+v", x)
		send := sendDiscordMessageEmbed(curr, s, c)
		// send := true
		if !send {
			return false
		}
	} else {
		return false
	}
	return true
}

func sendDiscordMessageEmbed(r *ResponseEarthquake, s *discordgo.Session, c map[string]*ChannelIDMemoryStore) bool {
	imageUrl := "https://data.bmkg.go.id/DataMKG/TEWS" + r.Infogempa.Gempa.Shakemap

	field := []*discordgo.MessageEmbedField{
		{
			Name:   "Wilayah",
			Value:  r.Infogempa.Gempa.LocationInformation,
			Inline: true,
		},
		{
			Name:   "Potensi",
			Value:  r.Infogempa.Gempa.Potential,
			Inline: true,
		},
		{
			Name:   "Tanggal",
			Value:  r.Infogempa.Gempa.Date,
			Inline: true,
		},
		{
			Name:   "Jam",
			Value:  r.Infogempa.Gempa.Time,
			Inline: true,
		},
		{
			Name:   "Koordinat",
			Value:  r.Infogempa.Gempa.Coordinates,
			Inline: true,
		},
		{
			Name:   "Magnitude",
			Value:  r.Infogempa.Gempa.Magnitude + " sr",
			Inline: true,
		},
		{
			Name:   "Kedalaman",
			Value:  r.Infogempa.Gempa.Depth,
			Inline: true,
		},
		{
			Name:   "Dirasakan di",
			Value:  r.Infogempa.Gempa.FeltAt,
			Inline: true,
		},
	}
	log.Println(c)

	mes := discordgo.MessageEmbed{
		URL:       "https://data.bmkg.go.id/DataMKG/TEWS/autogempa.json",
		Timestamp: r.Infogempa.Gempa.DateTime,
		Image: &discordgo.MessageEmbedImage{
			URL: imageUrl,
		},
		Title: os.Getenv("peringatanHeader"),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    os.Getenv("peringatanFooter"),
			IconURL: os.Getenv("peringatanFooterURL"),
		},
		Description: r.Infogempa.Gempa.LocationInformation,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: os.Getenv("peringatanFooterURL"),
		},
		Fields: field,
	}
	for _, guild := range s.State.Guilds {
		_, err := s.ChannelMessageSendEmbed(c[guild.ID].ChannelID, &mes)
		if err != nil {
			log.Printf("Error: Failed to send embedded message, %v", err)
			return false
		}
		log.Printf("Info: Earthquake message send to guild '%v'", guild.Name)

	}
	return true
}
