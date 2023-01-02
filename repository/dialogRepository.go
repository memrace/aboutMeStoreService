package repository

import (
	"aboutMeStoreService/entities"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type IDialogRepository interface {
	GetDialog(userID int64) (*entities.Dialog, error)
	UpdateDialog(dialog *entities.Dialog) (bool, error)
	RemoveDialog(userID int64) (bool, error)
	CreateDialog(dialog *entities.Dialog) (int64, error)
	CloseConnection()
	GetDb() *sql.DB
}

func MakeDialogRepository(db *sql.DB) IDialogRepository {
	return &dialogRepository{db}
}

type dialogRepository struct {
	db *sql.DB
}

func (repo *dialogRepository) GetDb() *sql.DB {
	return repo.db
}

func (repo *dialogRepository) CloseConnection() {
	err := repo.db.Close()
	if err != nil {
		_ = fmt.Errorf("could not close db: %w", err)
	}
}
func (repo *dialogRepository) CreateDialog(dialog *entities.Dialog) (int64, error) {

	_, err := repo.db.Exec(
		"insert into dialogs (id, userName, firstName, lastName, chatId, reply, replied) values ($1, $2, $3, $4, $5, $6, $7)",
		dialog.Id, dialog.UserName, dialog.FirstName, dialog.LastName, dialog.ChatID, dialog.Reply, dialog.Replied,
	)
	if err != nil {
		println(err)
		return 0, err
	}
	return dialog.Id, err
}
func (repo *dialogRepository) GetDialog(id int64) (*entities.Dialog, error) {

	row := repo.db.QueryRow("select * from dialogs where id = $1", id)
	dialog := entities.Dialog{}
	err := row.Scan(&dialog.Id, &dialog.UserName, &dialog.FirstName, &dialog.LastName, &dialog.ChatID, &dialog.Reply, &dialog.Replied)
	if err != nil {
		return nil, err
	}
	return &dialog, err
}

func (repo *dialogRepository) UpdateDialog(dialog *entities.Dialog) (bool, error) {

	result, err := repo.db.Exec("update dialogs set userName = $1, firstName = $2, lastName = $3, chatId = $4, reply = $5, replied = $6 where id = $7",
		dialog.UserName, dialog.FirstName, dialog.LastName, dialog.ChatID, dialog.Reply, dialog.Replied, dialog.Id,
	)
	if err != nil {
		return false, err
	}
	amount, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if amount == 0 {
		return false, errors.New("нет сущности")
	}

	return true, nil
}

func (repo *dialogRepository) RemoveDialog(id int64) (bool, error) {

	result, err := repo.db.Exec("delete from dialogs where id = $1", id)

	if err != nil {
		return false, err
	}

	amount, amountErr := result.RowsAffected()
	if amountErr != nil {
		return false, amountErr
	}
	if amount == 0 {
		return false, errors.New("нет сущности")
	}

	return true, nil
}
