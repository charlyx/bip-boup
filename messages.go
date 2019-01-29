package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func handleMessage(session *discordgo.Session, message *discordgo.Message) {
	channel, err := session.State.Channel(message.ChannelID)
	if err != nil {
		return
	}

	isDM := channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM

	var member *discordgo.Member = nil
	guild, err := session.State.Guild(channel.GuildID)
	if err != nil {
		if !isDM {
			return
		}
	} else {
		_member, err := session.GuildMember(guild.ID, message.Author.ID)
		if err != nil && !isDM {
			return
		}
		member = _member
	}

	content := message.Content
	if len(content) <= 0 {
		return
	}

	var responseEmbed *discordgo.MessageEmbed

	prefix := ""
	if strings.HasPrefix(content, BotData.CommandPrefix) {
		prefix = BotData.CommandPrefix
	}

	if prefix != "" {
		Debug.Println("[" + channel.Name + "] " + message.Author.Username + ": " + content)

		commandContent := strings.TrimPrefix(content, prefix)
		command := strings.Split(commandContent, " ")

		responseEmbed = callCommand(command[0], command[1:], &CommandEnvironment{
			Guild: guild, Channel: channel,
			User: message.Author, Member: member,
			Message: message,
		})
	}

	if responseEmbed != nil {
		session.ChannelMessageSendEmbed(message.ChannelID, responseEmbed)
	}
}