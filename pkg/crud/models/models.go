package models

type Book struct {
	Id            int64
	Name          string
	Pages         int
	Removed       bool
	FileName      string
	Category      string
	AmountOfLikes int
	LiveBook      bool
}

type Comments struct {
	Id              int64
	Comment         string
	BookId          int64
	CommentatorName string
	Removed         bool
}

type Users struct {
	id int64
	Name string
	Login string
	Password string
	Role bool
	Ban bool
	Removed bool
}
