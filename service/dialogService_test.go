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
	id, err := service.Create(0, "test", "test1", "test2", 0)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidId)
	assert.Equal(t, id, int64(0))

	id, err = service.Create(1, "test", "test1", "test2", 1)

	assert.Equal(t, id, int64(1))
	assert.NoError(t, err)

	dialog, err := service.Get(id)

	assert.NoError(t, err)
	assert.NotNil(t, dialog)
	assert.False(t, dialog.Replied)
	assert.Equal(t, dialog.Reply, "")

	_, err = service.Create(1, "test", "test1", "test2", 1)

	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrDialogAlreadyExists)
}

func TestDialogService_Get(t *testing.T) {
	beforeEach()
	dialog, err := service.Get(0)
	assert.Error(t, err)
	assert.Nil(t, dialog)
	assert.ErrorIs(t, err, repository.ErrDialogNotExists)

	id, _ := service.Create(1, "test", "test1", "test2", 1)
	dialog, err = service.Get(id)
	assert.NoError(t, err)
	assert.NotNil(t, dialog)
	assert.Equal(t, id, dialog.Id)
	assert.Equal(t, dialog.UserName, "test")
}

func TestDialogService_Delete(t *testing.T) {
	beforeEach()
	id, _ := service.Create(1, "test", "test1", "test2", 1)
	err := service.Delete(id)
	assert.NoError(t, err)

	err = service.Delete(id)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrDialogNotExists)

}

func TestDialogService_SetReply(t *testing.T) {
	beforeEach()
	id, _ := service.Create(1, "test", "test1", "test2", 1)
	dialog, _ := service.Get(id)

	reply := "testReply"

	err := service.SetReply(dialog.Id, reply)
	assert.NoError(t, err)

	dialog, _ = service.Get(dialog.Id)
	assert.Equal(t, reply, dialog.Reply)
	assert.True(t, dialog.Replied)

	err = service.SetReply(dialog.Id, reply)
	assert.Error(t, err)
	assert.ErrorIs(t, err, entities.ErrDialogAlreadyHasReply)

	err = service.SetReply(dialog.Id, "")
	assert.Error(t, err)
	assert.ErrorIs(t, err, entities.ErrEmptyMessage)

	err = service.SetReply(0, "")
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrDialogNotExists)

}
