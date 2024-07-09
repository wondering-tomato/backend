package mysql

import (
	"backend/explore"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllLiked_NoLastId(t *testing.T) {
	// Given.
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := StoreMySQLImpl{db: db}
	rows := sqlmock.NewRows([]string{"ID, ActorID"})
	mock.ExpectQuery("SELECT ID, ActorID FROM decisions WHERE RecipientID=? AND Liked=1 ORDER BY id desc LIMIT ?").WithArgs(1, 1).WillReturnRows(rows)

	// When.
	ctx := context.Background()
	allLiked, lastId, err := client.GetAllLiked(ctx, 1, 0)

	// Then.
	require.NoError(t, err)
	assert.Equal(t, []*explore.ListLikedYouResponse_Liker([]*explore.ListLikedYouResponse_Liker(nil)), allLiked)
	assert.Equal(t, 0, lastId)
}

func TestGetAllLiked_LastId(t *testing.T) {
	// Given.
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := StoreMySQLImpl{db: db}
	rows := sqlmock.NewRows([]string{"ID, ActorID"})
	mock.ExpectQuery("SELECT ID, ActorID FROM decisions WHERE RecipientID=? AND Liked=1 AND ID < ? ORDER BY id desc LIMIT ?").WithArgs(1, 5, 1).WillReturnRows(rows)

	// When.
	ctx := context.Background()
	allLiked, lastId, err := client.GetAllLiked(ctx, 1, 5)

	// Then.
	require.NoError(t, err)
	assert.Equal(t, []*explore.ListLikedYouResponse_Liker([]*explore.ListLikedYouResponse_Liker(nil)), allLiked)
	assert.Equal(t, 0, lastId)
}

func TestGetNewAllLiked_NoLastId(t *testing.T) {
	// Given.
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := StoreMySQLImpl{db: db}
	rows := sqlmock.NewRows([]string{"ID, ActorID"})
	mock.ExpectQuery("SELECT ID, ActorID FROM decisions WHERE RecipientID=? AND Liked=1 AND ActorID NOT IN (SELECT RecipientID FROM decisions WHERE ActorID=? AND Liked=1) ORDER BY ID desc LIMIT ?").WithArgs(1, 1, 1).WillReturnRows(rows)

	// When.
	ctx := context.Background()
	allNewLiked, lastId, err := client.GetNewAllLiked(ctx, 1, 0)

	// Then.
	require.NoError(t, err)
	assert.Equal(t, []*explore.ListLikedYouResponse_Liker([]*explore.ListLikedYouResponse_Liker(nil)), allNewLiked)
	assert.Equal(t, 0, lastId)
}

func TestGetNewAllLiked_LastID(t *testing.T) {
	// Given.
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := StoreMySQLImpl{db: db}
	rows := sqlmock.NewRows([]string{"ID, ActorID"})
	mock.ExpectQuery("SELECT ID, ActorID FROM decisions WHERE RecipientID=? AND Liked=1 AND ActorID NOT IN (SELECT RecipientID FROM decisions WHERE ActorID=? AND Liked=1) AND ID < ? ORDER BY ID desc LIMIT ?").WithArgs(1, 1, 5, 1).WillReturnRows(rows)

	// When.
	ctx := context.Background()
	allNewLiked, lastId, err := client.GetNewAllLiked(ctx, 1, 5)

	// Then.
	require.NoError(t, err)
	assert.Equal(t, []*explore.ListLikedYouResponse_Liker([]*explore.ListLikedYouResponse_Liker(nil)), allNewLiked)
	assert.Equal(t, 0, lastId)
}

func TestCountLikedYou(t *testing.T) {
	// Given.
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := StoreMySQLImpl{db: db}
	rows := sqlmock.NewRows([]string{"ActorID"})
	mock.ExpectQuery("SELECT ActorID FROM decisions WHERE RecipientID=? AND Liked=1").WithArgs(1).WillReturnRows(rows)

	// When.
	ctx := context.Background()
	count, err := client.CountLikedYou(ctx, 1)

	// Then.
	require.NoError(t, err)
	assert.Equal(t, uint64(0), count)
}

func TestPutDecision(t *testing.T) {
	// Given.
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	client := StoreMySQLImpl{db: db}
	rows := sqlmock.NewRows([]string{"COUNT(*)"})
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT(*) FROM decisions WHERE ActorID=? AND RecipientID=? AND Liked=?").WithArgs(1, 2, 0).WillReturnRows(rows.AddRow(1))
	mock.ExpectExec("DELETE FROM decisions WHERE ActorID=? AND RecipientID=? AND Liked=?").WithArgs(1, 2, 0).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT COUNT(*) FROM decisions WHERE ActorID=? AND RecipientID=? AND Liked=?").WithArgs(1, 2, 1).WillReturnRows(rows.AddRow(1))
	mock.ExpectExec("DELETE FROM decisions WHERE ActorID=? AND RecipientID=? AND Liked=?").WithArgs(1, 2, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO decisions (ActorID, RecipientID, Liked) VALUES (?, ?, ?)").WithArgs(1, 2, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT COUNT(*) FROM decisions WHERE ActorID=? AND RecipientID=? AND Liked=?").WithArgs(1, 2, 1).WillReturnRows(rows.AddRow(1))
	mock.ExpectQuery("SELECT COUNT(*) FROM decisions WHERE ActorID=? AND RecipientID=? AND Liked=?").WithArgs(2, 1, 1).WillReturnRows(rows.AddRow(1))
	mock.ExpectCommit()

	// When.
	ctx := context.Background()
	result, err := client.PutDecision(ctx, 1, 2, 1)

	// Then.
	require.NoError(t, err)
	assert.True(t, result)
}
