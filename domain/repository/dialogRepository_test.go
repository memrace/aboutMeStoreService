package repository

import (
	"aboutMeStoreService/domain/repository/migrations"
	"aboutMeStoreService/entities"
	"context"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

var repo IDialogRepository

type dbTestConnection struct {
	d     string
	path  string
	mPath string
}

var dbCon = dbTestConnection{
	"sqlite3", "file:test.db?cache=shared&mode=memory", "migrations",
}

func run(m *testing.M) (code int, err error) {
	migrator := migrations.NewMigrator(dbCon.d, dbCon.path, dbCon.mPath)

	migrator.UpToLastVersion()

	repo = MakeDialogRepository(migrator.Db)

	defer migrator.Close()

	return m.Run(), nil
}

func TestCreateDialog(t *testing.T) {
	beforeEach()

	testDialog := createTestDialog(123)

	id, err := repo.Create(testDialog)

	assert.NoError(t, err)
	assert.Equal(t, id, testDialog.Id)

	testDialog2 := createTestDialog(12)

	_, err = repo.Create(testDialog2)

	assert.Error(t, err)

	testDialog2.ChatID = 2345
	testDialog2.UserName = "23"

	id, err = repo.Create(testDialog2)
	assert.NoError(t, err)
	assert.Equal(t, id, testDialog2.Id)
}

func TestGetDialog(t *testing.T) {
	beforeEach()

	dialogPointer := createTestDialog(12)
	id, err := repo.Create(dialogPointer)

	assert.NoError(t, err)

	dialogReturnPointer, err := repo.Get(id)
	assert.NotNil(t, dialogReturnPointer)
	assert.NoError(t, err)
	dialogReturn := *dialogReturnPointer
	dialog := *dialogPointer
	assert.Exactly(t, dialog, dialogReturn)
}

func TestUpdateDialog(t *testing.T) {
	beforeEach()

	dialogPointer := createTestDialog(12)
	res, err := repo.Update(dialogPointer)

	assert.False(t, res)
	assert.Error(t, err)

	id, err := repo.Create(dialogPointer)

	assert.NoError(t, err)
	assert.Equal(t, id, dialogPointer.Id)

	dialogPointer.Reply = "rep2"
	res, err = repo.Update(dialogPointer)

	assert.True(t, res)
	assert.NoError(t, err)

	dialogReturnPointer, err := repo.Get(dialogPointer.Id)
	dialogReturn := *dialogReturnPointer
	dialog := *dialogPointer
	assert.NoError(t, err)
	assert.Exactly(t, dialog, dialogReturn)
}

func TestDeleteDialog(t *testing.T) {
	beforeEach()

	dialog := createTestDialog(12)

	res, err := repo.Delete(dialog.Id)

	assert.Error(t, err)
	assert.False(t, res)

	_, _ = repo.Create(dialog)

	res, err = repo.Delete(dialog.Id)

	assert.NoError(t, err)
	assert.True(t, res)
}

func createTestDialog(userId int64) *entities.Dialog {
	return &entities.Dialog{Id: userId, UserName: "test", FirstName: "t1", LastName: "t2", ChatID: 1234, Reply: "rep", Replied: true}
}

func beforeEach() {
	_, err := repo.GetDb().ExecContext(context.TODO(), "DELETE FROM dialogs")
	if err != nil {
		_ = fmt.Errorf(err.Error())
		panic(err)
	}
}
