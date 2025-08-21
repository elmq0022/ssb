package schemas

type ArticleCreateSchema struct {
	Title    string `db:"title" json:"title"`
	Body     string `db:"body" json:"body"`
	UserName string `db:"user_name" json:"Author"`
}

type ArticleUpdateSchema struct {
	Title    *string `db:"title" json:"title,omitempty"`
	Body     *string `db:"body" json:"body,omitempty"`
	UserName *string `db:"user_name" json:"author:omitempty"`
}

// for GetByID and ListAll
type ArticleWithAuthorSchema struct {
	Title  string `db:"title" json:"title"`
	Body   string `db:"body" json:"body"`
	Author struct {
		UserName  string `db:"user_name" json:"user_name"`
		FirstName string `db:"first_name" json:"first_name"`
		LastName  string `db:"last_name" json:"last_name"`
	}
}
