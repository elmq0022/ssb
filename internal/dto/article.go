package dto

type ArticleCreateDTO struct {
	Title string `json:"title"`
	Body string `json:"body"`
	Author string `json:"Author"`
}

type ArticleUpdateDTO struct {
	Title  *string `json:"title,omitempty"`
	Body   *string `json:"body,omitempty"`
	Author *string `json:"author:omitempty"`
}
