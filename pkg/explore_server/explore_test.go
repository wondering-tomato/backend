package explore_server

import (
	"backend/explore"
	"backend/mocks"
	"context"
	"encoding/hex"
	"strconv"
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
	token := hex.EncodeToString([]byte("3"))
	mockStore := mocks.NewStore(t)
	mockStore.On("GetAllLiked", context.TODO(), 2, 0).Return(GetAllLikedResponse, 3, nil).Once()

	// When.
	newExploreServerImpl := New(mockStore)
	res, err := newExploreServerImpl.ListLikedYou(context.TODO(), ListLikedYouResponse)

	// Then.
	require.NoError(t, err)
	assert.Len(t, res.Likers, 2)
	assert.Contains(t, res.Likers, GetAllLikedResponse[0])
	assert.Contains(t, res.Likers, GetAllLikedResponse[1])

	assert.Equal(t, res.NextPaginationToken, &token)
}

func TestHexEncodeDecode(t *testing.T) {
	// Given.
	hasMorePages := 7
	intHasMorePages := strconv.Itoa(hasMorePages)

	// When.
	encodedIntHasMorePages := hex.EncodeToString([]byte(intHasMorePages))
	decodedIntHasMorePages, err := hex.DecodeString(encodedIntHasMorePages)
	require.NoError(t, err)

	// Then.
	decodedHasMorePages, err := strconv.Atoi(string(decodedIntHasMorePages))
	require.NoError(t, err)
	assert.Equal(t, hasMorePages, decodedHasMorePages)
}
