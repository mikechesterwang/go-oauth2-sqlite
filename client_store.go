package sqlite

import (
	"context"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
)

type ClientStore struct {
	repository *ClientRepository
}

func NewClientStore(filename string) (*ClientStore, error) {
	repo, err := NewClientRepository(filename)
	if err != nil {
		return nil, err
	}
	return &ClientStore{
		repository: repo,
	}, nil
}

func (c *ClientStore) Set(ctx context.Context, info oauth2.ClientInfo) error {
	return c.repository.Insert(&ClientItem{
		Id:     info.GetID(),
		Secret: info.GetSecret(),
		Domain: info.GetDomain(),
		UserId: info.GetUserID(),
	})
}

func (c *ClientStore) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	item, err := c.repository.GetById(id)
	if err != nil {
		return nil, err
	}
	return &models.Client{
		ID:     item.Id,
		Secret: item.Secret,
		Domain: item.Domain,
		UserID: item.UserId,
	}, nil
}

func (c *ClientStore) RemoveByID(ctx context.Context, id string) error {
	return c.repository.RemoveById(id)
}
