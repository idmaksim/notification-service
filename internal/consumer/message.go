package consumer

type Message struct {
	Target  string `json:"target"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}
