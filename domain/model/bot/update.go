package bot

// Update is a Telegram object that the handler receives every time an user interacts with the bot.
type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}
