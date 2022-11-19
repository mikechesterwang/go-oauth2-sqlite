package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type ClientRepository struct {
	db *sql.DB
}

type ClientItem struct {
	Id     string
	Secret string
	Domain string
	UserId string
}

func (repo *ClientRepository) Close() {
	err := repo.db.Close()
	if err != nil {
		repo.error(err)
	}
}

func (repo *ClientRepository) errorf(format string, args ...interface{}) {
	fmt.Printf("[OAUTH2-SQLITE-ERROR]: "+format+"\n", args...)
}

func (repo *ClientRepository) error(err error) {
	repo.errorf("%s", err.Error())
}

func NewClientRepository(filename string) (*ClientRepository, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	repo := &ClientRepository{
		db: db,
	}
	_, err = db.Exec(`BEGIN;
		CREATE TABLE IF NOT EXISTS oauth2_client (
			id         VARCHAR(512) PRIMARY KEY,
			secret     VARCHAR(512),
			domain     VARCHAR(2048),
			user_id    VARCHAR(512),
		);
		COMMIT;`)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (repo *ClientRepository) Insert(item *ClientItem) error {
	if item == nil {
		return errors.New("client info cannot be nil")
	}
	_, err := repo.db.Exec(
		"INSERT INTO oauth2_client (id, secret, domain, user_id) VALUES (?, ?, ?, ?)",
		item.Id, item.Secret, item.Domain, item.UserId,
	)
	return err
}

func (repo *ClientRepository) GetById(id string) (*ClientItem, error) {
	rows, err := repo.db.Query("SELECT * FROM oauth2_client WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	var item ClientItem
	err = rows.Scan(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *ClientRepository) RemoveById(id string) error {
	_, err := repo.db.Exec("DELETE FROM oauth2_client WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}
