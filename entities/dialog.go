package entities

type Dialog struct {
	Id        int64
	UserName  string
	FirstName string
	LastName  string
	ChatID    int64
	Reply     string
	Replied   bool
}
