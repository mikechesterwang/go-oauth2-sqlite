package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	_ "github.com/mattn/go-sqlite3"
)

type TokenRepository struct {
	db     *sql.DB
	ticker *time.Ticker
}

type TokenItem struct {
	Id        int64
	ExpiredAt int64
	Code      string
	Access    string
	Refresh   string
	Data      string
}

func (repo *TokenRepository) Close() {
	repo.ticker.Stop()
	err := repo.db.Close()
	if err != nil {
		repo.error(err)
	}
}

func (repo *TokenRepository) errorf(format string, args ...interface{}) {
	fmt.Printf("[OAUTH2-SQLITE-ERROR]: "+format+"\n", args...)
}

func (repo *TokenRepository) error(err error) {
	repo.errorf("%s", err.Error())
}

func (repo *TokenRepository) clean() error {
	_, err := repo.db.Exec("DELETE FROM oauth2_token WHERE expired_at <= ? OR (code = '' AND access = '' AND refresh = '')", time.Now())
	return err
}

func (repo *TokenRepository) gc() {
	for range repo.ticker.C {
		if err := repo.clean(); err != nil {
			repo.error(err)
		}
	}
}

func NewTokenRepository(filename string) (*TokenRepository, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	repo := &TokenRepository{
		db:     db,
		ticker: time.NewTicker(time.Minute * 10),
	}
	_, err = db.Exec(`BEGIN;
	CREATE TABLE IF NOT EXISTS oauth2_token (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		expired_at INTEGER,
		code       VARCHAR(256),
		access     VARCHAR(256),
		refresh    VARCHAR(256),
		data       VARCHAR(2048)
	);
		CREATE INDEX idx_code ON oauth2_token (code);
		CREATE INDEX idx_access ON oauth2_token (access);
		CREATE INDEX idx_refresh ON oauth2_token (refresh);
		CREATE INDEX idx_expired_at ON oauth2_token (expired_at);
	COMMIT;`)
	if err != nil {
		return nil, err
	}
	go repo.gc()
	return repo, nil
}

func (repo *TokenRepository) Insert(item *TokenItem) error {
	if item == nil {
		return errors.New("token cannot be nil")
	}
	_, err := repo.db.Exec(
		"INSERT INTO oauth2_token (expired_at, code, access, refresh, data) VALUES (?, ?, ?, ?, ?)",
		item.ExpiredAt, item.Code, item.Access, item.Refresh, item.Data,
	)
	return err
}

func (repo *TokenRepository) RemoveCode(code string) error {
	_, err := repo.db.Exec("UPDATE oauth2_token SET code = '' WHERE code = ?", code)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (repo *TokenRepository) RemoveAccess(access string) error {
	_, err := repo.db.Exec("UPDATE oauth2_token SET access = '' WHERE access = ?", access)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (repo *TokenRepository) RemoveRefresh(refresh string) error {
	_, err := repo.db.Exec("UPDATE oauth2_token SET refresh = '' WHERE refresh = ?", refresh)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

func (repo *TokenRepository) unmarshalData(data string) oauth2.TokenInfo {
	var m models.Token
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		repo.error(err)
	}
	return &m
}

func (repo *TokenRepository) getToken(col string, val string) (oauth2.TokenInfo, error) {
	if val == "" {
		return nil, nil
	}
	query := fmt.Sprintf("SELECT data FROM oauth2_token WHERE %s = ?", col)
	rows, err := repo.db.Query(query)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	var data string
	if err := rows.Scan(&data); err != nil {
		return nil, err
	} else {
		return repo.unmarshalData(data), nil
	}
}

func (repo *TokenRepository) GetTokenByCode(code string) (oauth2.TokenInfo, error) {
	return repo.getToken("code", code)
}

func (repo *TokenRepository) GetTokenByAccess(access string) (oauth2.TokenInfo, error) {
	return repo.getToken("access", access)
}

func (repo *TokenRepository) GetTokenByRefresh(refresh string) (oauth2.TokenInfo, error) {
	return repo.getToken("refresh", refresh)
}
