package entities

import "errors"

var DialogAlreadyHasReply = errors.New("ответ уже дан")

var EmptyMessage = errors.New("пустой ответ")

type Dialog struct {
	Id        int64
	UserName  string
	FirstName string
	LastName  string
	ChatID    int64
	Reply     string
	Replied   bool
}

func (dialog *Dialog) SetReply(message string) error {

	if message == "" {
		return EmptyMessage
	}

	if dialog.Replied {
		return DialogAlreadyHasReply
	}

	dialog.Reply = message
	dialog.Replied = true
	return nil
}
