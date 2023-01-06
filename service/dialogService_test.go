package service

import (
	"aboutMeStoreService/configuration"
	"aboutMeStoreService/domain/repository"
	"aboutMeStoreService/domain/repository/migrations"
	"aboutMeStoreService/entities"
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cfg = configuration.DbTestConnectionConfiguration("../domain/repository/migrations")

var service *DialogService

func withTestRepository(db *sql.DB) DialogServiceConfiguration {
	repo := repository.MakeDialogRepository(db)
	return func(ds *DialogService) error {
		ds.repository = repo
		return nil
	}
}

func before() {
	migrator := migrations.New(
		cfg.DriverName,
		cfg.DataSourceName,
		cfg.MigrationsPath)

	err := migrator.UpToLastVersion()
	if err != nil {
		log.Fatal(err)
	}

	service, err = NewDialogService(withTestRepository(migrator.Db))

	if err != nil {
		log.Fatal(err)
	}

}

func after() {
	service.EndSession()
	service = nil
}

func TestDialogService_Create(t *testing.T) {
	before()
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

	defer after()
}

func TestDialogService_Get(t *testing.T) {
	before()

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
	defer after()
}

func TestDialogService_Delete(t *testing.T) {
	before()

	id, _ := service.Create(1, "test", "test1", "test2", 1)
	err := service.Delete(id)
	assert.NoError(t, err)

	err = service.Delete(id)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrDialogNotExists)

	defer after()
}

func TestDialogService_SetReply(t *testing.T) {
	before()
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

	defer after()
}
