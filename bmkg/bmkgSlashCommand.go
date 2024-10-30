package bmkg

import "github.com/bwmarrin/discordgo"

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "about",
			Description: "About bot.",
		},
		{
			Name:        "get-earthquake",
			Description: "Shows last 15 earthquake.",
		},
		{
			Name:        "get-weather-forecast",
			Description: "Shows weather forecast",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"about": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "BMKG Bot is a bot that shows bmkg's open data (https://data.bmkg.go.id/) with automatic earthquake retrieval.",
				},
			})
		},
		// todo get weather forecast and last 15 earthquake

	}
)
