package main

import (
	"bot-bmkg/bmkg"
	"bot-bmkg/util"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	cStore := make(map[string]*bmkg.ChannelIDMemoryStore)

	// Create stop signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Load env from file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Fatal: Error loading .env file")
	}
	BotToken, exist := os.LookupEnv("BotToken")
	if !exist {
		log.Fatal("Fatal: environment variable 'BotToken' does not exist")
	}

	// Create new bot session without connecting
	s, _ := discordgo.New("Bot " + BotToken)

	// add event handler
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Info: Bot is ready")
	})
	s.AddHandler(func(s *discordgo.Session, r *discordgo.GuildCreate) {
		log.Printf("Info: Bot joined new guild: %v", r.Name)

		channelMap, err := util.GetDiscordTextChannelFromGuildID(s, r.ID)
		if err != nil {
			log.Printf("Error: failed to get channel ID, %v", err)
		}

		cStore[r.ID] = &bmkg.ChannelIDMemoryStore{
			ChannelName: channelMap["channelName"],
			ChannelID:   channelMap["channelID"],
		}

	})

	// Create a websocket connection to discord
	err = s.Open()
	if err != nil {
		log.Fatalf("Fatal: Cannot open the session: %v", err)
	}

	defer s.Close()

	// Create inMemoryStore struct
	earthquake := &bmkg.InMemoryStore{
		LatestEarthquake: time.Now(),
	}

	go func() {
		for {
			if bmkg.NewEarthquakeHandler(earthquake, s, cStore) {
				log.Println("Info: Message send to guild channel successfully")
			}
			time.Sleep(7 * time.Second)
		}
	}()
	// Debugging

	// wait for stop signal (ctrl + c)
	<-stop
	log.Println("Info: Graceful shutdown")
}
