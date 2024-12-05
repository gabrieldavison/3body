package connections

type HTTPMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

var HTTPMessageChannel = make(chan HTTPMessage)
