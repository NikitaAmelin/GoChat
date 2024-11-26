package domain

type Message struct {
	ID             string   `json:"ID"`
	Author         string   `json:"author"`
	Text           string   `json:"text"`
	TimeOfSending  string   `json:"timeOfSending"`
	Viewed         []string `json:"viewed"`
	NameMessagesDB string   `json:"nameMessegesDB"`
}
