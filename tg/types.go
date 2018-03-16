package tg

type Update struct {
	UpdateID      int     `json:"update_id"`
	Message       Message `json:"message"`
	EditedMessage Message `json:"edited_message"`
}

type Message struct {
	MessageID         int    `json:"message_id"`
	Date              int    `json:"date"`
	Text              string `json:"text"`
	Caption           string `json:"caption"`
	From              User   `json:"from"`
	ForwardFrom       User   `json:"forward_from"`
	Chat              Chat   `json:"chat"`
	ForwardedFromChat Chat   `json:"forwarded_from_chat"`
}

type User struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bit"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Chat struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Username    string `json:"username"`
	Description string `json:"description"`
}

type SendMessageRequest struct {
	ChatID int
	Text   string
}
