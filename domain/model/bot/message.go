package bot

import "github.com/EduardTruuvaart/ev-chargepoint-tracker/domain/model"

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text     string          `json:"text"`
	Chat     Chat            `json:"chat"`
	Location *model.Location `json:"location"`
}
