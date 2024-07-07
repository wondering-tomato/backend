package explore_server

import (
	"backend/explore"
	"backend/pkg/store"
	"context"
	"strconv"
)

type ExploreServer interface {
}

func New(s store.Store) *ExploreServerImpl {
	return &ExploreServerImpl{
		Store:                             s,
		UnimplementedExploreServiceServer: explore.UnimplementedExploreServiceServer{},
	}
}

type ExploreServerImpl struct {
	Store store.Store
	explore.UnimplementedExploreServiceServer
}

// List all users who liked the recipient
func (es *ExploreServerImpl) ListLikedYou(rctx context.Context, in *explore.ListLikedYouRequest) (*explore.ListLikedYouResponse, error) {
	id, err := strconv.Atoi(in.GetRecipientUserId())
	if err != nil {
		return nil, err
	}
	actorIDs, err := es.Store.GetAllLiked(rctx, id)
	if err != nil {
		return nil, err
	}
	return &explore.ListLikedYouResponse{
		Likers:              actorIDs,
		NextPaginationToken: new(string),
	}, nil
}

// List all users who liked the recipient
func (es *ExploreServerImpl) PutDecision(rctx context.Context, in *explore.PutDecisionRequest) (*explore.PutDecisionResponse, error) {
	actorUserId, err := strconv.Atoi(in.ActorUserId)
	if err != nil {
		return nil, err
	}
	recipientUserId, err := strconv.Atoi(in.RecipientUserId)
	if err != nil {
		return nil, err
	}
	liked := 0
	if in.LikedRecipient {
		liked = 1
	}

	res, err := es.Store.PutDecision(rctx, actorUserId, recipientUserId, liked)
	if err != nil {
		return nil, err
	}
	return &explore.PutDecisionResponse{
		MutualLikes: res,
	}, nil
}

// List all users who liked the recipient excluding those who have been liked in return
func (es *ExploreServerImpl) ListNewLikedYou(rctx context.Context, in *explore.ListLikedYouRequest) (*explore.ListLikedYouResponse, error) {
	id, err := strconv.Atoi(in.GetRecipientUserId())
	if err != nil {
		return nil, err
	}
	actorIDs, err := es.Store.GetNewAllLiked(rctx, id)
	if err != nil {
		return nil, err
	}
	return &explore.ListLikedYouResponse{
		Likers:              actorIDs,
		NextPaginationToken: new(string),
	}, nil
}

// CountLikedYou counts the number of users who liked the recipient.
func (es *ExploreServerImpl) CountLikedYou(rctx context.Context, in *explore.CountLikedYouRequest) (*explore.CountLikedYouResponse, error) {
	id, err := strconv.Atoi(in.GetRecipientUserId())
	if err != nil {
		return nil, err
	}
	count, err := es.Store.CountLikedYou(rctx, id)
	if err != nil {
		return nil, err
	}
	return &explore.CountLikedYouResponse{
		Count: count,
	}, nil
}
