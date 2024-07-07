package explore_server

import (
	"backend/explore"
	"backend/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListLikedYou(t *testing.T) {
	// Given.
	ListLikedYouResponse := &explore.ListLikedYouRequest{
		RecipientUserId: "2",
		PaginationToken: new(string),
	}
	GetAllLikedResponse := []*explore.ListLikedYouResponse_Liker{
		{
			ActorId:       "1",
			UnixTimestamp: 0,
		},
		{
			ActorId:       "3",
			UnixTimestamp: 0,
		},
	}
	mockStore := mocks.NewStore(t)
	mockStore.On("GetAllLiked", context.TODO(), 2).Return(GetAllLikedResponse, nil).Once()
	// When.
	newExploreServerImpl := New(mockStore)
	res, err := newExploreServerImpl.ListLikedYou(context.TODO(), ListLikedYouResponse)

	// Then.
	require.NoError(t, err)
	assert.Len(t, res.Likers, 2)
	assert.Contains(t, res.Likers, GetAllLikedResponse[0])
	assert.Contains(t, res.Likers, GetAllLikedResponse[1])
}