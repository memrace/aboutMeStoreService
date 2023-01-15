package service

import (
	"aboutMeStoreService/configuration"
	"aboutMeStoreService/domain/repository"
	"aboutMeStoreService/entities"
	context "context"
	"database/sql"
	"log"
)

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
		configuration.DbConnConfiguration.DriverName,
		configuration.DbConnConfiguration.DataSourceName,
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

func (service DialogService) mustEmbedUnimplementedDialogServiceServer() {}

func (service DialogService) Get(ctx context.Context, dto *GetDialog) (*Dialog, error) {
	dialog, err := service.repository.Get(dto.GetUserId())
	if err != nil {
		return nil, err
	}

	return &Dialog{
		Id:        dialog.Id,
		UserName:  dialog.UserName,
		FirstName: dialog.FirstName,
		LastName:  dialog.LastName,
		ChatId:    dialog.ChatID,
		Reply:     dialog.Reply,
		Replied:   dialog.Replied,
	}, nil
}

func (service DialogService) Create(ctx context.Context, dto *CreateDialog) (*DialogId, error) {
	if dto.GetId() <= 0 || dto.GetChatId() <= 0 {
		return &DialogId{
			Id: 0,
		}, ErrInvalidId
	}
	newDialog := entities.Dialog{
		Id:        dto.GetId(),
		UserName:  dto.GetUserName(),
		FirstName: dto.GetFirstName(),
		LastName:  dto.GetLastName(),
		ChatID:    dto.GetChatId(),
		Replied:   false,
	}
	id, err := service.repository.Create(&newDialog)
	return &DialogId{
		Id: id,
	}, err
}

func (service DialogService) Delete(ctx context.Context, id *DialogId) (*Result, error) {
	err := service.repository.Delete(id.GetId())
	if err != nil {
		return &Result{Success: false}, err
	}
	return &Result{Success: true}, nil
}

func (service DialogService) SetReply(ctx context.Context, reply *UserReply) (*Result, error) {
	dialog, err := service.repository.Get(reply.GetUserId())

	if err != nil {
		return &Result{Success: false}, err
	}

	err = dialog.SetReply(reply.GetText())

	if err != nil {
		return &Result{Success: false}, err
	}

	err = service.repository.Update(dialog)

	if err != nil {
		return &Result{Success: false}, err
	}
	return &Result{Success: true}, nil
}

func (service *DialogService) EndSession() {
	service.repository.CloseConnection()
}
