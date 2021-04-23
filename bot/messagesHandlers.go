package bot

import dg "github.com/bwmarrin/discordgo"

func textMessageHandler(s *dg.Session, m *dg.MessageCreate) {

	_, _ = reply(s, m.Message, "Hello!")
}
