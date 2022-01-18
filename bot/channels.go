package bot

var (
	// Last confirmation from the admin
	IsWaitingConfirmation bool = false

	// OutboundChannel represents the channel for the outgoing messages from the bot
	OutboundChannel = make(chan string)

	// ConfirmationChannel represents the channel for the confirmation messages sent from the admin
	ConfirmationChannel = make(chan bool)
)
