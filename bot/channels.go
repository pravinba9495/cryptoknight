package bot

var (
	// ChatID of the admin on telegram
	ChatID string

	// Last confirmation from the admin
	IsWaitingConfirmation bool = false

	// OutboundChannel represents the channel for the outgoing messages from the bot
	OutboundChannel = make(chan string)

	// ErrorChannel represents the channel for the error messages to be sent to the admin
	ErrorChannel = make(chan error)

	// ErrorChannel represents the channel for the confirmation messages sent from the admin
	ConfirmationChannel = make(chan bool)
)
