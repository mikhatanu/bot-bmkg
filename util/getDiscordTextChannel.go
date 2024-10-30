package util

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// Get first text channel from guil did
func GetDiscordTextChannelFromGuildID(s *discordgo.Session, guildID string) (map[string]string, error) {
	channels, _ := s.GuildChannels(guildID)
	returnMap := make(map[string]string)
	flag := 0
	for _, c := range channels {
		if c.Type == discordgo.ChannelTypeGuildText {
			if c.Name == "gempa-alert" {
				returnMap["channelID"] = c.ID
				returnMap["channelName"] = c.Name
				break
			} else if c.Name == "general" && flag == 0 {
				returnMap["channelID"] = c.ID
				returnMap["channelName"] = c.Name
				flag = 1
			}
		}
	}

	if len(returnMap) == 0 {
		return nil, errors.New("No text channel found in Guild ID" + guildID)
	}
	return returnMap, nil
}
