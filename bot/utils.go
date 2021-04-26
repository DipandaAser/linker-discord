package bot

import (
	"errors"
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"io"
	"net/http"
	"strings"
)

//parseMessageToMessageSend parse a *dg.Message to *dg.MessageSend
func parseMessageToMessageSend(m *dg.Message) *dg.MessageSend {

	msg := &dg.MessageSend{
		Content:         m.Content,
		TTS:             m.TTS,
		AllowedMentions: &dg.MessageAllowedMentions{},
	}

	// get the first embeds message
	if len(m.Embeds) > 0 {
		msg.Embed = m.Embeds[0]
	}

	for i, attachment := range m.Attachments {
		reader, err := getFileByUrl(attachment.URL)
		if err != nil {
			continue
		}
		msg.Files = append(msg.Files, &dg.File{
			Name:   m.Attachments[i].Filename,
			Reader: reader,
		})
	}

	//mentions
	msg.AllowedMentions.Parse = append(msg.AllowedMentions.Parse, dg.AllowedMentionTypeEveryone)
	msg.AllowedMentions.Parse = append(msg.AllowedMentions.Parse, dg.AllowedMentionTypeUsers)
	msg.AllowedMentions.Parse = append(msg.AllowedMentions.Parse, dg.AllowedMentionTypeRoles)

	return msg
}

func buildErrorResponse(description string) *dg.MessageSend {

	return &dg.MessageSend{
		Embed: &dg.MessageEmbed{
			Title:       "Error",
			Description: description,
			Color:       embedErrorColor,
		},
	}
}

func buildInfoResponse(description string) *dg.MessageSend {

	return &dg.MessageSend{
		Embed: &dg.MessageEmbed{
			Title:       "Info",
			Description: description,
			Color:       embedInfoColor,
		},
	}
}

func getFileByUrl(fileUrl string) (io.Reader, error) {

	response, err := http.Get(fileUrl)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("can't get file")
	}

	return response.Body, nil
}

func formatChannelIdToLinkerGroupId(guildId string, channelId string) string {

	return fmt.Sprintf("%s|%s", guildId, channelId)
}

func formatLinkerGroupIdToDiscordChannelId(groupID string) (guildId string, channelId string) {

	ids := strings.Split(groupID, "|")
	if len(ids) == 2 {
		return ids[0], ids[1]
	}

	return "", ""
}

func getPayload(command string, numberOfParameterWanted int, m *dg.MessageCreate) ([]string, error) {

	payload := strings.Split(strings.TrimPrefix(strings.TrimPrefix(m.Content, command), " "), " ")

	if len(payload) < numberOfParameterWanted {
		return payload, errors.New("insufficient parameters")
	}

	return payload, nil
}

func formatCommandName(commandName string) string {
	return fmt.Sprintf("%s%s", commandPrefix, commandName)
}

func formatStatus(status bool) string {
	if status {
		return ":green_circle:"
	}

	return ":red_circle:"
}
