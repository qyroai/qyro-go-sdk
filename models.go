package qyro

type Session struct {
	ID string `json:"id"`
}

type Message struct {
	ID      string `json:"id"`
	Role    string `json:"role"`
	Content string `json:"content"`
}