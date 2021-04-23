package bot

import (
	"github.com/DipandaAser/linker-discord/app"
	dg "github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var session *dg.Session
var commandsHandlers = make(map[string]func(s *dg.Session, m *dg.MessageCreate))

func Init() error {
	s, err := dg.New("Bot " + app.DiscordBotToken)
	if err != nil {
		return err
	}
	session = s

	commandsHandlers[formatCommandName("help")] = helpHandler
	commandsHandlers[formatCommandName("list")] = listHandler
	commandsHandlers[formatCommandName("config")] = configHandler
	commandsHandlers[formatCommandName("start")] = startHandler
	commandsHandlers[formatCommandName("stop")] = stopHandler
	commandsHandlers[formatCommandName("link")] = linkHandler
	commandsHandlers[formatCommandName("diffuse")] = diffuseHandler

	session.AddHandler(globalHandlerMessageCreate)

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	session.Identify.Intents = dg.IntentsGuilds | dg.IntentsGuildMessages | dg.IntentsGuildVoiceStates | dg.IntentsDirectMessages | dg.IntentsGuildMembers

	// Open the websocket and begin listening.
	err = session.Open()
	if err != nil {
		return err
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Linker is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	session.Close()

	return nil
}

//globalHandlerMessageCreate handle all incoming message
func globalHandlerMessageCreate(s *dg.Session, m *dg.MessageCreate) {

	// Ignore all messages:
	//  - created by the bot itself,
	//	- from other bot
	//	- from discord system
	if m.Author.ID == s.State.User.ID || m.Author.System || m.Author.Bot {
		return
	}

	// we check if is a command message
	payload := strings.Split(m.Content, " ")
	if len(payload) != 0 {
		commandName := payload[0]
		if handlerFunc, commandExist := commandsHandlers[commandName]; commandExist {
			handlerFunc(s, m)
			return
		}
	}

	if m.Type != dg.MessageTypeDefault {
		return
	}

	textMessageHandler(s, m)

	/*userChannel, _ := s.UserChannelCreate(m.Author.ID)

	if strings.HasPrefix(m.Message.Content, "!") {
		fmt.Println("command called")
		_, err := session.ChannelMessageSend(userChannel.ID, fmt.Sprintf("You run the command %s", m.Message.Content))
		if err != nil {
			fmt.Println("cant send message")
			return
		}
		return
	}
	_, err := session.ChannelMessageSend(userChannel.ID, m.Message.Content)
	if err != nil {
		fmt.Println("cant send message")
		return
	}
	fmt.Println("message sent")*/
}

//isUserAdmin check if the user have the permission to manage server
func isUserAdmin(s *dg.Session, m *dg.Message) (bool, error) {

	//p, err := s.State.UserChannelPermissions(userID, channelID)  //this is deprecate
	p, err := s.State.MessagePermissions(m)
	b := p&dg.PermissionManageServer == dg.PermissionManageServer
	return b, err
}

func reply(s *dg.Session, m *dg.Message, message string) (*dg.Message, error) {
	ref := &dg.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	}
	return s.ChannelMessageSendReply(m.ChannelID, message, ref)
}

func replyWithComplex(s *dg.Session, m *dg.Message, msg *dg.MessageSend) (*dg.Message, error) {

	msg.Reference = &dg.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	}

	return s.ChannelMessageSendComplex(m.ChannelID, msg)
}
