package bot

import (
	"errors"
	"fmt"
	"github.com/DipandaAser/linker"
	"github.com/DipandaAser/linker-discord/bot/groups"
	dg "github.com/bwmarrin/discordgo"
	"strings"
)

const (
	ErrUseCmdInGroup         = "Please use this command on a group where Linker is in"
	ErrParameterInsufficient = "Insufficient parameters"
	ErrBotNotSupported       = "Sorry but for now other bot can't interact with linker"
	ErrGlobal                = "Something going wrong during the operation.\nPlease retry" +
		"\nIf it is persistent, contact the [project owner](https://www.twitter.com/iamdipanda)"
	ErrInvalidCode  = "Please provide two valid Linker Group Code"
	FOOTER          = "By Dipanda Aser"
	commandPrefix   = "!"
	embedInfoColor  = 1404802
	embedErrorColor = 16711680
)

func formatChannelID(guildId string, channelId string) string {

	return fmt.Sprintf("%s|%s", guildId, channelId)
}

func formatCommandName(commandName string) string {
	return fmt.Sprintf("%s%s", commandPrefix, commandName)
}

func helpHandler(s *dg.Session, m *dg.MessageCreate) {

	commandsField := []*dg.MessageEmbedField{}
	for _, command := range linker.GetCommands() {
		cmd := command
		commandsField = append(commandsField, &dg.MessageEmbedField{
			Name:   formatCommandName(cmd.Text) + " " + cmd.Option,
			Value:  cmd.Description,
			Inline: false,
		})
	}

	mem := &dg.MessageEmbed{
		Title:       "Available commands",
		Description: "[Linker Discord Server](https://discord.gg/dVawHP9gB3)\n[Linker Project](https://github.com/DipandaAser/linker-discord)",
		Color:       embedInfoColor,
		Footer: &dg.MessageEmbedFooter{
			Text: FOOTER,
		},
		Fields: commandsField,
	}

	_, _ = replyWithComplex(s, m.Message, &dg.MessageSend{Embed: mem})
}

func listHandler(s *dg.Session, m *dg.MessageCreate) {

	/*payload, err := getPayload(formatCommandName("list"), 1, m)
	if err != nil {
		msg := &dg.MessageSend{
			Embed:           buildErrorResponse(fmt.Sprintf("%s, %s", ERR_PARAMETER_INSUFFICIENT,
				"provide ")),
		}
		_, _ = replyWithComplex(s, m.Message, msg)
		return
	}*/

}

func configHandler(s *dg.Session, m *dg.MessageCreate) {

	userChannel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		return
	}

	// We use this command only in a group channel
	if m.GuildID == "" {
		_, _ = s.ChannelMessageSendComplex(userChannel.ID, buildErrorResponse(ErrUseCmdInGroup))
		return
	}

	isAdmin, err := isUserAdmin(s, m.Message)
	if err != nil {
		_, _ = s.ChannelMessageSendComplex(userChannel.ID, buildErrorResponse(ErrGlobal))
		return
	}
	if !isAdmin {
		_, _ = s.ChannelMessageSendComplex(userChannel.ID, buildErrorResponse("You don't have the permission to manage this server"))
		return
	}

	group, err := groups.VerifyGroupExistenceAndCreateIfNot(formatChannelID(m.GuildID, m.ChannelID))
	if err != nil {
		_, _ = s.ChannelMessageSendComplex(userChannel.ID, buildErrorResponse(ErrGlobal))
		return
	}

	guildName := ""
	channelName := ""
	if guild, err := s.Guild(m.GuildID); err == nil {
		guildName = guild.Name
	}
	if channel, err := s.Channel(m.ChannelID); err == nil {
		channelName = channel.Name
	}

	chatName := fmt.Sprintf("%s --> %s", guildName, channelName)
	msg := fmt.Sprintf("Hey Dude this is the linker id of the ***%s*** group/channel. \nLinker Group Code: **%s**", chatName, group.ShortCode)

	_, _ = s.ChannelMessageSendComplex(userChannel.ID, buildInfoResponse(msg))
	return
}

func startHandler(s *dg.Session, m *dg.MessageCreate) {

	payload, err := getPayload(formatCommandName("start"), 1, m)
	if err != nil {
		msg := buildErrorResponse(fmt.Sprintf("%s, %s", ErrParameterInsufficient, "please provide a link or diffusion id"))
		_, _ = replyWithComplex(s, m.Message, msg)
		return
	}

	if strings.TrimSpace(payload[0]) == "" {
		return
	}

	lnk, err := linker.GetLinkByID(payload[0])
	if err == nil {
		err := lnk.StartLink()
		if err != nil {
			_, _ = replyWithComplex(s, m.Message, buildErrorResponse(ErrGlobal))
			return
		}
		_, _ = replyWithComplex(s, m.Message, buildInfoResponse("Link is now active"))
		return
	}

	diff, err := linker.GetDiffusionById(payload[0])
	if err == nil {
		err := diff.StartDiffusion()
		if err != nil {
			_, _ = replyWithComplex(s, m.Message, buildErrorResponse(ErrGlobal))
			return
		}
		_, _ = replyWithComplex(s, m.Message, buildInfoResponse("Diffusion is now active"))
		return
	}

	_, _ = replyWithComplex(s, m.Message, buildErrorResponse("Please provide a good link or diffusion id"))
}

func stopHandler(s *dg.Session, m *dg.MessageCreate) {

	payload, err := getPayload(formatCommandName("stop"), 1, m)
	if err != nil {
		msg := buildErrorResponse(fmt.Sprintf("%s, %s", ErrParameterInsufficient, "please provide a link or diffusion id"))
		_, _ = replyWithComplex(s, m.Message, msg)
		return
	}

	if strings.TrimSpace(payload[0]) == "" {
		return
	}

	lnk, err := linker.GetLinkByID(payload[0])
	if err == nil {
		err := lnk.StopLink()
		if err != nil {
			_, _ = replyWithComplex(s, m.Message, buildErrorResponse(ErrGlobal))
			return
		}
		_, _ = replyWithComplex(s, m.Message, buildInfoResponse("Link is now deactivate"))
		return
	}

	diff, err := linker.GetDiffusionById(payload[0])
	if err == nil {
		err := diff.StopDiffusion()
		if err != nil {
			_, _ = replyWithComplex(s, m.Message, buildErrorResponse(ErrGlobal))
			return
		}
		_, _ = replyWithComplex(s, m.Message, buildInfoResponse("Diffusion is now deactivate"))
		return
	}

	_, _ = replyWithComplex(s, m.Message, buildErrorResponse("Please provide a good link or diffusion id"))
}

func linkHandler(s *dg.Session, m *dg.MessageCreate) {

	payload, err := getPayload(formatCommandName("link"), 2, m)
	if err != nil {
		msg := buildErrorResponse(fmt.Sprintf("%s, %s", ErrParameterInsufficient, "please provide two linker group/channel id"))
		_, _ = replyWithComplex(s, m.Message, msg)
		return
	}

	var firstGroup, secondGroup *linker.Group
	if firstGroup, err = linker.GetGroupByShortCode(payload[0]); err != nil {
		_, _ = replyWithComplex(s, m.Message, buildErrorResponse(ErrInvalidCode))
		return
	}

	if secondGroup, err = linker.GetGroupByShortCode(payload[1]); err != nil {
		_, _ = replyWithComplex(s, m.Message, buildErrorResponse(ErrInvalidCode))
		return
	}

	// we check if these groups already have link together
	if lnk, _ := linker.GetLinksByGroupsID([2]string{firstGroup.ID, secondGroup.ID}); lnk != nil {

		// we have a link who match
		_, _ = replyWithComplex(s, m.Message, buildInfoResponse("This link already exist."))
		return
	}

	_, err = linker.CreateLink([2]string{firstGroup.ID, secondGroup.ID})
	if err != nil {
		_, _ = replyWithComplex(s, m.Message, buildErrorResponse(ErrGlobal))
		return
	}

	_, _ = replyWithComplex(s, m.Message, buildInfoResponse("Link successfully created. \nYou can start exchange message between these Groups."))
	return
}

func diffuseHandler(s *dg.Session, m *dg.MessageCreate) {

	payload, err := getPayload(formatCommandName("diffuse"), 2, m)
	if err != nil {
		msg := buildErrorResponse(fmt.Sprintf("%s, %s", ErrParameterInsufficient, "please provide two linker group/channel id"))
		_, _ = replyWithComplex(s, m.Message, msg)
		return
	}

	// we check if these groups exist
	var broadcasterGroup, receiverGroup *linker.Group
	if broadcasterGroup, err = linker.GetGroupByShortCode(payload[0]); err != nil {
		_, _ = replyWithComplex(s, m.Message, buildErrorResponse(ErrInvalidCode))
		return
	}

	if receiverGroup, err = linker.GetGroupByShortCode(payload[1]); err != nil {
		_, _ = replyWithComplex(s, m.Message, buildErrorResponse(ErrInvalidCode))
		return
	}

	// we check if these groups already have diffusion together
	if diff, _ := linker.GetDiffusionsByBroadcasterAndReceiver(broadcasterGroup.ID, receiverGroup.ID); diff != nil {

		// we have a diffusion who match
		_, _ = replyWithComplex(s, m.Message, buildInfoResponse("This diffusion already exist."))
		return
	}

	_, err = linker.CreateDiffusion(broadcasterGroup.ID, receiverGroup.ID)
	if err != nil {
		_, _ = replyWithComplex(s, m.Message, buildErrorResponse(ErrGlobal))
		return
	}

	_, _ = replyWithComplex(s, m.Message, buildInfoResponse("Diffusion successfully created."))
	return
}

func getPayload(command string, numberOfParameterWanted int, m *dg.MessageCreate) ([]string, error) {

	payload := strings.Split(strings.TrimPrefix(strings.TrimPrefix(m.Content, command), " "), " ")

	if len(payload) < numberOfParameterWanted {
		return payload, errors.New("insufficient parameters")
	}

	return payload, nil
}

func buildErrorResponse(description string) *dg.MessageSend {

	return &dg.MessageSend{
		Embed: &dg.MessageEmbed{
			Title:       "An error occur",
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
