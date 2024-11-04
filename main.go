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

	// Initialize some var
	cStore := make(map[string]*bmkg.ChannelIDMemoryStore)
	registeredCommands := make(map[string][]*discordgo.ApplicationCommand)
	earthquake := &bmkg.InMemoryStore{
		LatestEarthquake: time.Now(),
	}

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
	// On ready or start
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Info: Bot is ready")
	})
	// On guild join (also triggered every startup because bot will join the guild and trigger this)
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
		log.Printf("Info: added channel ID to guild: %v", cStore[r.ID])

		log.Printf("Info: Adding commands to guild %v", r.Name)
		for _, v := range bmkg.Commands {
			go func() {
				cmd, err := s.ApplicationCommandCreate(s.State.User.ID, r.ID, v)
				if err != nil {
					log.Panicf("Cannot create '%v' command: %v", v.Name, err)
				}
				registeredCommands[r.ID] = append(registeredCommands[r.ID], cmd)
			}()

		}
	})
	// On slash command interaction
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := bmkg.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.GuildDelete) {
		log.Printf("Info: Guild unavailable, or bot was removed from guild: %v. Deleting channel ID %v", i.ID, cStore[i.ID])

		delete(cStore, i.ID)
	})

	// Create a websocket connection to discord
	err = s.Open()
	if err != nil {
		log.Fatalf("Fatal: Cannot open the session: %v", err)
	}
	defer s.Close()

	// Run a function that checks new earthquake every 7 second
	go func() {
		for {
			if cStore != nil && bmkg.NewEarthquakeHandler(earthquake, s, cStore) {
				log.Println("Info: Message send to guild channel successfully")
				log.Printf("Info: Current earthquake time is %v", earthquake.LatestEarthquake)
			}
			time.Sleep(7 * time.Second)
		}
	}()
	// Debugging

	// wait for stop signal (ctrl + c)
	<-stop

	// execute shutdown command
	log.Println("Info: Graceful shutdown")

	log.Println("Info: Removing commands...")
	for g := range registeredCommands {
		for _, v := range registeredCommands[g] {
			go func() {
				err := s.ApplicationCommandDelete(s.State.User.ID, g, v.ID)
				if err != nil {
					log.Printf("Cannot delete '%v' command from guild name %v: %v", v.ID, g, err)
				}
			}()

		}
	}

}
