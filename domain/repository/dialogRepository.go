package repository

import (
	"aboutMeStoreService/entities"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/mattn/go-sqlite3"
)

type IDialogRepository interface {
	Get(userID int64) (*entities.Dialog, error)
	Update(dialog *entities.Dialog) (bool, error)
	Delete(userID int64) (bool, error)
	Create(dialog *entities.Dialog) (int64, error)
	CloseConnection()
	GetDb() *sql.DB
	ping()
}

func MakeDialogRepository(db *sql.DB) IDialogRepository {
	return &dialogRepository{db}
}

type dialogRepository struct {
	db *sql.DB
}

func (repo *dialogRepository) ping() {
	if err := repo.db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func (repo *dialogRepository) GetDb() *sql.DB {
	repo.ping()
	return repo.db
}

func (repo *dialogRepository) CloseConnection() {
	repo.ping()
	err := repo.db.Close()
	if err != nil {
		_ = fmt.Errorf("could not close db: %w", err)
	}
}

var DialogAlreadyExists = errors.New("диалог уже существует")

func (repo *dialogRepository) Create(dialog *entities.Dialog) (int64, error) {
	repo.ping()
	_, err := repo.db.Exec(
		"insert into dialogs (id, userName, firstName, lastName, chatId, reply, replied) values ($1, $2, $3, $4, $5, $6, $7)",
		dialog.Id, dialog.UserName, dialog.FirstName, dialog.LastName, dialog.ChatID, dialog.Reply, dialog.Replied,
	)

	if err != nil {
		sqliteError := err.(sqlite3.Error)
		println(err)
		if sqliteError.Code == 19 {
			return 0, DialogAlreadyExists
		}
		return 0, sqliteError
	}
	return dialog.Id, nil
}
func (repo *dialogRepository) Get(id int64) (*entities.Dialog, error) {
	repo.ping()
	row := repo.db.QueryRow("select * from dialogs where id = $1", id)
	dialog := entities.Dialog{}
	err := row.Scan(&dialog.Id, &dialog.UserName, &dialog.FirstName, &dialog.LastName, &dialog.ChatID, &dialog.Reply, &dialog.Replied)
	if err != nil {
		return nil, err
	}
	return &dialog, err
}

func (repo *dialogRepository) Update(dialog *entities.Dialog) (bool, error) {
	repo.ping()
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

func (repo *dialogRepository) Delete(id int64) (bool, error) {
	repo.ping()
	result, err := repo.db.Exec("delete from dialogs where id = $1", id)

	if err != nil {
		return false, err
	}

	amount, amountErr := result.RowsAffected()
	if amountErr != nil {
		return false, amountErr
	}
	if amount == 0 {
		return false, sql.ErrNoRows
	}

	return true, nil
}
