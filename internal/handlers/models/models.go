package models

type CreateChat struct {
	Title string `json:"title"`
}

type CreateChatReposnse struct {
	Status string `json:"status"`
	Title  string `json:"title"`
}
