package store

import (
	"backend/explore"
	"backend/pkg/store/mysql"
	"backend/pkg/store/noop"
	"context"
)

type Store interface {
	GetAllLiked(ctx context.Context, id int, lastId int) ([]*explore.ListLikedYouResponse_Liker, int, error)
	GetNewAllLiked(ctx context.Context, id int) ([]*explore.ListLikedYouResponse_Liker, error)
	CountLikedYou(ctx context.Context, id int) (uint64, error)
	PutDecision(ctx context.Context, actorId int, recipientId int, liked int) (bool, error)
}

type StoreImpl struct {
	Store
}

func New() Store {
	db := mysql.New()
	if db == nil {
		// Use a noop db when the mysql connection fails.
		return noop.New()
	}
	// Using the mysql db.
	return db
}

// Downstream interfaces conform to parent Store interface.
var _ Store = &StoreImpl{}
var _ Store = &mysql.StoreMySQLImpl{}
var _ Store = &noop.StoreNoopImpl{}

func (s *StoreImpl) GetAllLiked(ctx context.Context, id int, lastId int) ([]*explore.ListLikedYouResponse_Liker, int, error) {
	return nil, 0, nil
}
func (s *StoreImpl) GetNewAllLiked(ctx context.Context, id int) ([]*explore.ListLikedYouResponse_Liker, error) {
	return nil, nil
}

func (s *StoreImpl) CountLikedYou(ctx context.Context, id int) (uint64, error) { return 0, nil }
func (s *StoreImpl) PutDecision(ctx context.Context, actorId int, recipientId int, liked int) (bool, error) {
	return false, nil
}
