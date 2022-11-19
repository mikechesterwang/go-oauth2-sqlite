package sqlite

import (
	"context"
	"encoding/json"

	"github.com/go-oauth2/oauth2/v4"
)

type TokenInfo = oauth2.TokenInfo

type TokenStore struct {
	repository *TokenRepository
}

func NewSqliteStore(filename string) (oauth2.TokenStore, error) {
	r, err := NewTokenRepository(filename)
	if err != nil {
		return nil, err
	}
	return &TokenStore{repository: r}, nil
}

func (s *TokenStore) Close() {
	s.repository.Close()
}

func (s *TokenStore) Create(ctx context.Context, info TokenInfo) error {
	infoRaw, err := json.Marshal(info)
	if err != nil {
		return err
	}
	item := &TokenItem{
		Data: string(infoRaw),
	}
	if code := info.GetCode(); code != "" {
		item.Code = code
		item.ExpiredAt = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn()).Unix()
	} else {
		item.Access = info.GetAccess()
		item.ExpiredAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).Unix()
	}
	if refresh := info.GetRefresh(); refresh != "" {
		item.Refresh = info.GetRefresh()
		item.ExpiredAt = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Unix()
	}
	return s.repository.Insert(item)
}

func (s *TokenStore) RemoveByCode(ctx context.Context, code string) error {
	return s.repository.RemoveCode(code)
}

// use the access token to delete the token information
func (s *TokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return s.repository.RemoveAccess(access)
}

// use the refresh token to delete the token information
func (s *TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return s.repository.RemoveRefresh(refresh)
}

// use the authorization code for token information data
func (s *TokenStore) GetByCode(ctx context.Context, code string) (TokenInfo, error) {
	return s.repository.GetTokenByCode(code)
}

// use the access token for token information data
func (s *TokenStore) GetByAccess(ctx context.Context, access string) (TokenInfo, error) {
	return s.repository.GetTokenByAccess(access)
}

func (s *TokenStore) GetByRefresh(ctx context.Context, refresh string) (TokenInfo, error) {
	return s.repository.GetTokenByRefresh(refresh)
}
