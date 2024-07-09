package explore_server

import (
	"backend/explore"
	"backend/pkg/store"
	"context"
	"encoding/hex"
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

// List all users who liked the recipient.
// Page size is currently set to 1 to demonstrate the pagination functionality with
// the current small amount of test data in the database.
// The page size can be increased in the file pkg/store/mysql/mysql_queries.go
func (es *ExploreServerImpl) ListLikedYou(rctx context.Context, in *explore.ListLikedYouRequest) (*explore.ListLikedYouResponse, error) {
	id, err := strconv.Atoi(in.GetRecipientUserId())
	if err != nil {
		return nil, err
	}

	previousToken := in.GetPaginationToken()
	var lastId int
	if previousToken != "" {
		lastId, err = decodeToken(previousToken)
		if err != nil {
			return nil, err
		}
	}
	actorIDs, hasMorePages, err := es.Store.GetAllLiked(rctx, id, lastId)
	if err != nil {
		return nil, err
	}

	var token string
	if hasMorePages > 0 {
		intHasMorePages := strconv.Itoa(hasMorePages)
		token = hex.EncodeToString([]byte(intHasMorePages))
	}
	return &explore.ListLikedYouResponse{
		Likers:              actorIDs,
		NextPaginationToken: &token,
	}, nil
}

// List all users who liked the recipient excluding those who have been liked in return.
// Page size is currently set to 1 to demonstrate the pagination functionality with
// the current small amount of test data in the database.
// The page size can be increased in the file pkg/store/mysql/mysql_queries.go
func (es *ExploreServerImpl) ListNewLikedYou(rctx context.Context, in *explore.ListLikedYouRequest) (*explore.ListLikedYouResponse, error) {
	id, err := strconv.Atoi(in.GetRecipientUserId())
	if err != nil {
		return nil, err
	}

	previousToken := in.GetPaginationToken()
	var lastId int
	if previousToken != "" {
		lastId, err = decodeToken(previousToken)
		if err != nil {
			return nil, err
		}
	}

	actorIDs, hasMorePages, err := es.Store.GetNewAllLiked(rctx, id, lastId)
	if err != nil {
		return nil, err
	}
	var token string
	if hasMorePages > 0 {
		intHasMorePages := strconv.Itoa(hasMorePages)
		token = hex.EncodeToString([]byte(intHasMorePages))
	}
	return &explore.ListLikedYouResponse{
		Likers:              actorIDs,
		NextPaginationToken: &token,
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

// PutDecision adds a new decision between actor and recipient.
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

// decodeToken decodes a hex encoded string into an int.
func decodeToken(token string) (int, error) {
	decodedToken, err := hex.DecodeString(token)
	if err != nil {
		return 0, err
	}
	lastId, err := strconv.Atoi(string(decodedToken))
	if err != nil {
		return 0, err
	}
	return lastId, nil
}
