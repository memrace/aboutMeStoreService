package service

import (
	"aboutMeStoreService/configuration"
	"aboutMeStoreService/domain/repository"
	"aboutMeStoreService/entities"
	"database/sql"
	"errors"
	"log"
)

var InvalidId = errors.New("невалидный ид")

type DialogServiceConfiguration func(ds *DialogService) error

type DialogService struct {
	repository repository.IDialogRepository
}

func NewDialogService(cfgs ...DialogServiceConfiguration) (*DialogService, error) {
	service := &DialogService{}

	for _, cfg := range cfgs {
		err := cfg(service)
		if err != nil {
			return nil, err
		}
	}

	return service, nil
}

func WithLocalRepository() DialogServiceConfiguration {
	db, err := sql.Open(
		configuration.DbConnectionConfiguration.DriverName,
		configuration.DbConnectionConfiguration.DataSourceName,
	)

	if err != nil {
		log.Fatal(err)
	}

	repo := repository.MakeDialogRepository(db)
	return func(ds *DialogService) error {
		ds.repository = repo
		return nil
	}
}

func (service *DialogService) Get(userId int64) (*entities.Dialog, error) {
	return service.repository.Get(userId)
}

func (service *DialogService) Create(
	id int64,
	userName string,
	firstName string,
	lastName string,
	chatID int64) (int64, error) {
	if id <= 0 || chatID <= 0 {
		return 0, InvalidId
	}
	newDialog := entities.Dialog{
		Id:        id,
		UserName:  userName,
		FirstName: firstName,
		LastName:  lastName,
		ChatID:    chatID,
		Replied:   false,
	}
	return service.repository.Create(&newDialog)
}

func (service *DialogService) Delete(userId int64) (bool, error) {
	return service.repository.Delete(userId)
}

func (service *DialogService) SetReply(userId int64, message string) error {
	dialog, err := service.repository.Get(userId)

	if err != nil {
		return err
	}

	err = dialog.SetReply(message)

	if err != nil {
		return nil
	}

	return nil
}

func (service *DialogService) EndSession() {
	service.repository.CloseConnection()
}
