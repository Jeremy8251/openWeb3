package models

type UserList struct {
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type UserInfo struct {
	UserName string `json:"username" form:"username" xml:"username"`
	Password string `json:"password" form:"password" xml:"password"`
}

type Article struct {
	Title   string
	Content string
	Score   int
}
