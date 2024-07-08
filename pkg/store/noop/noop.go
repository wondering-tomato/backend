package noop

import (
	"backend/explore"
	"context"
	"log"
)

type StoreNoop interface{}
type StoreNoopImpl struct {
}

var _ StoreNoop = &StoreNoopImpl{}

func (s *StoreNoopImpl) GetAllLiked(ctx context.Context, id int, lastId int) ([]*explore.ListLikedYouResponse_Liker, int, error) {
	return nil, 0, nil
}
func (s *StoreNoopImpl) GetNewAllLiked(ctx context.Context, id int, lastId int) ([]*explore.ListLikedYouResponse_Liker, int, error) {
	return nil, 0, nil
}
func (s *StoreNoopImpl) CountLikedYou(ctx context.Context, id int) (uint64, error) { return 0, nil }
func (s *StoreNoopImpl) PutDecision(ctx context.Context, actorId int, recipientId int, liked int) (bool, error) {
	return false, nil
}

func New() *StoreNoopImpl {
	log.Println("using store: noop")
	return &StoreNoopImpl{}
}
