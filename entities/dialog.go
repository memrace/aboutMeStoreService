package entities

import "errors"

var ErrDialogAlreadyHasReply = errors.New("ответ уже дан")

var ErrEmptyMessage = errors.New("пустой ответ")

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
		return ErrEmptyMessage
	}

	if dialog.Replied {
		return ErrDialogAlreadyHasReply
	}

	dialog.Reply = message
	dialog.Replied = true
	return nil
}
