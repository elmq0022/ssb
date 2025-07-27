package dto

type ArticleUpdateDTO struct {
	Title  *string `json:"title,omitempty"`
	Body   *string `json:"body,omitempty"`
	Author *string `json:"author:omitempty"`
}
