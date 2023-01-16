package service

import (
	"aboutMeStoreService/configuration"
	"aboutMeStoreService/domain/repository"
	"aboutMeStoreService/domain/repository/migrations"
	"aboutMeStoreService/entities"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cfg = configuration.DbTestConnectionConfiguration("../domain/repository/migrations")

var service *DialogService

func TestMain(m *testing.M) {

	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func withTestRepository(db *sql.DB) DialogServiceConfiguration {
	repo := repository.MakeDialogRepository(db)
	return func(ds *DialogService) error {
		ds.repository = repo
		return nil
	}
}

func run(m *testing.M) (code int, err error) {
	migrator := migrations.New(
		cfg.DriverName,
		cfg.DataSourceName,
		cfg.MigrationsPath)

	err = migrator.UpToLastVersion()
	if err != nil {
		log.Fatal(err)
	}

	service, err = NewDialogService(withTestRepository(migrator.Db))

	if err != nil {
		log.Fatal(err)
	}
	defer after()
	return m.Run(), nil
}

func after() {
	service.EndSession()
	service = nil
}

func beforeEach() {
	_, err := service.repository.GetDb().ExecContext(context.TODO(), "DELETE FROM dialogs")
	if err != nil {
		_ = fmt.Errorf(err.Error())
		panic(err)
	}
}

func TestDialogService_Create(t *testing.T) {
	beforeEach()
	id, err := service.Create(context.TODO(), &CreateDialog{
		UserName:  "test",
		FirstName: "test1",
		LastName:  "test2",
		ChatId:    0,
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidId)
	assert.Equal(t, id.GetId(), int64(0))

	id, err = service.Create(context.TODO(), &CreateDialog{
		Id:        1,
		UserName:  "test",
		FirstName: "test1",
		LastName:  "test2",
		ChatId:    1,
	})

	assert.Equal(t, id.GetId(), int64(1))
	assert.NoError(t, err)

	dialog, err := service.Get(context.TODO(), &GetDialog{
		UserId: id.GetId(),
	})

	assert.NoError(t, err)
	assert.NotNil(t, dialog)
	assert.False(t, dialog.GetReplied())
	assert.Equal(t, dialog.GetReply(), "")

	_, err = service.Create(context.TODO(), &CreateDialog{
		Id:        1,
		UserName:  "test",
		FirstName: "test1",
		LastName:  "test2",
		ChatId:    1,
	})

	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrDialogAlreadyExists)
}

func TestDialogService_Get(t *testing.T) {
	beforeEach()
	dialog, err := service.Get(context.TODO(), &GetDialog{UserId: 0})
	assert.Error(t, err)
	assert.Nil(t, dialog)
	assert.ErrorIs(t, err, repository.ErrDialogNotExists)

	id, _ := service.Create(context.TODO(), &CreateDialog{
		Id:        1,
		UserName:  "test",
		FirstName: "test1",
		LastName:  "test2",
		ChatId:    1,
	})
	dialog, err = service.Get(context.TODO(), &GetDialog{UserId: id.GetId()})
	assert.NoError(t, err)
	assert.NotNil(t, dialog)
	assert.Equal(t, dialog.GetId(), id.GetId())
	assert.Equal(t, dialog.GetUserName(), "test")
}

func TestDialogService_Delete(t *testing.T) {
	beforeEach()
	id, _ := service.Create(context.TODO(), &CreateDialog{
		Id:        1,
		UserName:  "test",
		FirstName: "test1",
		LastName:  "test2",
		ChatId:    1,
	})
	res, err := service.Delete(context.TODO(), &DialogId{Id: id.GetId()})
	assert.NoError(t, err)
	assert.True(t, res.GetSuccess())

	res, err = service.Delete(context.TODO(), &DialogId{Id: id.GetId()})
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrDialogNotExists)
	assert.False(t, res.GetSuccess())

}

func TestDialogService_SetReply(t *testing.T) {
	beforeEach()
	id, _ := service.Create(context.TODO(), &CreateDialog{
		Id:        1,
		UserName:  "test",
		FirstName: "test1",
		LastName:  "test2",
		ChatId:    1,
	})

	reply := "testReply"

	res, err := service.SetReply(context.TODO(), &UserReply{
		UserId: id.GetId(),
		Text:   reply,
	})
	assert.NoError(t, err)
	assert.True(t, res.GetSuccess())

	dialog, _ := service.Get(context.TODO(), &GetDialog{
		UserId: id.GetId(),
	})
	assert.Equal(t, reply, dialog.GetReply())
	assert.True(t, dialog.GetReplied())

	res, err = service.SetReply(context.TODO(), &UserReply{
		UserId: dialog.GetId(),
		Text:   reply,
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, entities.ErrDialogAlreadyHasReply)
	assert.False(t, res.GetSuccess())

	res, err = service.SetReply(context.TODO(), &UserReply{
		UserId: dialog.GetId(),
		Text:   "",
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, entities.ErrEmptyMessage)
	assert.False(t, res.GetSuccess())

	res, err = service.SetReply(context.TODO(), &UserReply{
		UserId: 0,
		Text:   "",
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrDialogNotExists)
	assert.False(t, res.GetSuccess())

}
