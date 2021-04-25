package bot

import (
	"github.com/DipandaAser/linker"
	"github.com/DipandaAser/linker-discord/app"
	"github.com/DipandaAser/linker-discord/bot/groups"
	dg "github.com/bwmarrin/discordgo"
)

func messageHandler(s *dg.Session, m *dg.MessageCreate) {

	//guildId , channelId
	group, getErr := groups.VerifyGroupExistenceAndCreateIfNot(formatChannelIdToLinkerGroupId(m.GuildID, m.ChannelID))
	if getErr != nil {
		return
	}

	_ = group.IncrementMessage()

	go linkSend(group, s, m.Message)
	go diffusionSend(group, s, m.Message)
}

func linkSend(group *linker.Group, s *dg.Session, m *dg.Message) {

	links, err := linker.GetLinksByGroupID(group.ID)
	if err != nil {
		return
	}

	msg := parseMessageToMessageSend(m)
	for _, lnk := range links {

		link := lnk
		// we skip inactive link
		if !link.Active {
			continue
		}

		var otherGroupID string
		for _, id := range link.GroupsID {
			if id != group.ID {
				otherGroupID = id
			}
		}
		grp, err := linker.GetGroupByID(otherGroupID)
		if err != nil {
			continue
		}

		if grp.Service == app.Config.ServiceName {
			_, channelId := formatLinkerGroupIdToDiscordChannelId(grp.ID)
			_, err = s.ChannelMessageSendComplex(channelId, msg)
			if err == nil {
				_ = link.IncrementMessage()
			}
		}

		// TODO implement send message to other service
	}
}

func diffusionSend(group *linker.Group, s *dg.Session, m *dg.Message) {

	diffusions, err := linker.GetDiffusionsByBroadcaster(group.ID)
	if err != nil {
		return
	}

	msg := parseMessageToMessageSend(m)

	for _, diff := range diffusions {

		diffusion := diff

		// we skip inactive diffusion
		if !diffusion.Active {
			continue
		}

		grp, err := linker.GetGroupByID(diffusion.Receiver)
		if err != nil {
			continue
		}

		if grp.Service == app.Config.ServiceName {
			_, channelId := formatLinkerGroupIdToDiscordChannelId(grp.ID)
			_, err = s.ChannelMessageSendComplex(channelId, msg)
			if err == nil {
				_ = diffusion.IncrementMessage()
			}
		}

		// TODO implement send message to other service
	}
}
