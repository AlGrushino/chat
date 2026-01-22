package models

type CreateChat struct {
	Title string `json:"title"`
}

type CreateChatReposnse struct {
	Status string `json:"status"`
	Title  string `json:"title"`
}

type CreateMessage struct {
	Text string `json:"text"`
}

type CreateMessageResponse struct {
	Status string `json:"status"`
	Text   string `json:"text"`
}

type GetMessagesResponse struct {
	Status   string   `json:"status"`
	ID       int      `json:"id"`
	Messages []string `json:"messages"`
}
