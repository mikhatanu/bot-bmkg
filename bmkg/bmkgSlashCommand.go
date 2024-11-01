package bmkg

import (
	"bot-bmkg/util"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	Commands = []*discordgo.ApplicationCommand{
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
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "kode-wilayah",
					Description: "Kode wilayah sesuai https://kodewilayah.id/ (e.g. 51.08.05.2001)",
					Required:    true,
				},
			},
		},
		{
			Name:        "get-kode-wilayah",
			Description: "Shows weather forecast",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "nama-wilayah",
					Description: "Nama wilayah terkecil (e.g. Pegadungan). Nama bisa duplikat. Ref: https://kodewilayah.id/",
					Required:    true,
				},
			},
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"about": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "BMKG Bot is a bot that shows bmkg's open data (https://data.bmkg.go.id/) with automatic earthquake retrieval.",
				},
			})
		},
		"get-weather-forecast": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			if len(strings.Split(options[0].StringValue(), ".")) != 4 {
				content := "Kode wilayah invalid."
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: content,
					},
				})
				return
			}

			req, err := GetWeatherForecast(options[0].StringValue())
			if err == nil {
				emb := createWeatherEmbedResponse(req)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: emb,
					},
				})
			} else {
				log.Printf("Error: error when getting weather forecast: %v", err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: err.Error(),
					},
				})
			}

		},
		"get-kode-wilayah": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			// get array of adm4 code
			admCode := util.GetAdmCodeFromLocation(strings.ToLower(options[0].StringValue()))

			// check length of adm4 code. 0 = empty, 1 = desired, > 1 = invalid
			if len(admCode) == 0 {
				log.Printf("Error: empty adm code. Options string value is not adm4 or adm4 not found")
				content := fmt.Sprintf("Administration code not found: %v", options[0].StringValue())

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: content,
					},
				})
				return
			}
			returnContent := "Kode wilayah: \n"
			for _, value := range admCode {
				returnContent += "**" + value + "**\t"
				fullLoc := util.GetFullAdministrationLocationName(value)
				returnContent += strings.Join(fullLoc, ", ") + "\n"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: returnContent,
				},
			})
		},
		// todo get last 15 earthquake

	}
)
