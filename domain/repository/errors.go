package repository

import "errors"

var ErrDialogAlreadyExists = errors.New("диалог уже существует")
var ErrDialogNotExists = errors.New("диалог не существует")
