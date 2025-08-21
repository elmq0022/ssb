package schemas

type ArticleCreateSchema struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	UserName string `json:"Author"`
}

type ArticleUpdateSchema struct {
	Title    *string `json:"title,omitempty"`
	Body     *string `json:"body,omitempty"`
	UserName *string `json:"author:omitempty"`
}

// for GetByID and ListAll
type ArticleWithAuthorSchema struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author struct {
		UserName  string `json:"user_name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
}
