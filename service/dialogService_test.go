package service

import (
	"aboutMeStoreService/configuration"
	"aboutMeStoreService/domain/repository"
	"aboutMeStoreService/domain/repository/migrations"
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
	assert.ErrorIs(t, err, InvalidId)
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
	assert.ErrorIs(t, err, repository.DialogAlreadyExists)

	defer after()
}

func TestDialogService_Get(t *testing.T) {
	before()

	dialog, err := service.Get(0)
	assert.Error(t, err)
	assert.Nil(t, dialog)
	assert.ErrorIs(t, err, DialogNotFound)

	id, _ := service.Create(1, "test", "test1", "test2", 1)
	dialog, err = service.Get(id)
	assert.NoError(t, err)
	assert.NotNil(t, dialog)
	assert.Equal(t, id, dialog.Id)
	assert.Equal(t, dialog.UserName, "test")
	defer after()
}

//
//func TestDialogService_Delete(t *testing.T) {
//	//service := beforeEach()
//}
//
//func TestDialogService_SetReply(t *testing.T) {
//	//service := beforeEach()
//}
