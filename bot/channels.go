package bot

var (
	// ChatID of the admin on telegram
	ChatID string

	// OutboundChannel represents the channel for the outgoing messages from the bot
	OutboundChannel = make(chan string)

	// ErrorChannel represents the channel for the error messages to be sent to the admin
	ErrorChannel = make(chan error)
)
