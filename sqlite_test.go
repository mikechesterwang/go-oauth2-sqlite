package sqlite

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var repo *TokenRepository

func TestMain(m *testing.M) {
	_ = os.Mkdir("tmp", 0755)
	dbName := "tmp/test.db"
	_, err := os.Stat(dbName)
	if err == nil {
		if err := os.Remove(dbName); err != nil {
			panic(err)
		}
	}
	repo, err = NewTokenRepository(dbName)
	if err != nil {
		panic(err)
	}
	code := m.Run()
	repo.Close()
	os.Exit(code)
}

func TestInsert(t *testing.T) {
	err := repo.Insert(&TokenItem{
		ExpiredAt: time.Now().Unix(),
		Code:      "code",
		Refresh:   "refresh",
		Access:    "access",
		Data:      "data",
	})
	assert.NoError(t, err)
}
