package mysql

import (
	"backend/explore"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
)

var pageSize = 1

// List all users who liked the recipient.
func (s *StoreMySQLImpl) GetAllLiked(ctx context.Context, id int, lastId int) ([]*explore.ListLikedYouResponse_Liker, int, error) {
	var rows *sql.Rows
	var err error
	if lastId > 0 {
		rows, err = s.db.QueryContext(ctx, "SELECT ID, ActorID FROM decisions WHERE RecipientID=? AND Liked=1 AND ID < ? ORDER BY id desc LIMIT ?", id, lastId, pageSize)
		if err != nil {
			return nil, 0, err
		}
	} else {
		rows, err = s.db.QueryContext(ctx, "SELECT ID, ActorID FROM decisions WHERE RecipientID=? AND Liked=1 ORDER BY id desc LIMIT ?", id, pageSize)
		if err != nil {
			return nil, 0, err
		}

	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	var result []*explore.ListLikedYouResponse_Liker
	minID := math.MaxInt32
	for rows.Next() {
		var id int
		var actorId explore.ListLikedYouResponse_Liker
		if err := rows.Scan(&id, &actorId.ActorId); err != nil {
			return nil, 0, err
		}
		minID = id
		result = append(result, &actorId)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// If there are more rows, return the last seen ID for the pagination response.
	if minID != math.MaxInt32 {
		return result, minID, nil
	}
	// If we've reached the end of the result, return the remainder of results and the id of 0.
	return result, 0, nil
}

func (s *StoreMySQLImpl) GetNewAllLiked(ctx context.Context, id int, lastId int) ([]*explore.ListLikedYouResponse_Liker, int, error) {
	var rows *sql.Rows
	var err error
	if lastId > 0 {
		rows, err = s.db.QueryContext(ctx, "SELECT ID, ActorID FROM decisions WHERE RecipientID=? AND Liked=1 AND ActorID NOT IN (SELECT RecipientID FROM decisions WHERE ActorID=? AND Liked=1) AND ID < ? ORDER BY ID desc LIMIT ?;", id, id, lastId, pageSize)
		if err != nil {
			return nil, 0, err
		}

	} else {
		rows, err = s.db.QueryContext(ctx, "SELECT ID, ActorID FROM decisions WHERE RecipientID=? AND Liked=1 AND ActorID NOT IN (SELECT RecipientID FROM decisions WHERE ActorID=? AND Liked=1) ORDER BY ID desc LIMIT ?;", id, id, pageSize)
		if err != nil {
			return nil, 0, err
		}
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	var result []*explore.ListLikedYouResponse_Liker
	minID := math.MaxInt32
	for rows.Next() {
		var id int
		var actorId explore.ListLikedYouResponse_Liker
		if err := rows.Scan(&id, &actorId.ActorId); err != nil {
			return nil, 0, err
		}
		minID = id
		result = append(result, &actorId)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// If there are more rows, return the last seen ID for the pagination response.
	if minID != math.MaxInt32 {
		return result, minID, nil
	}
	// If we've reached the end of the result, return the remainder of results and the id of 0.
	return result, 0, nil
}

func (s *StoreMySQLImpl) CountLikedYou(ctx context.Context, id int) (uint64, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT ActorID FROM decisions WHERE RecipientID=? AND Liked=1", id)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	var rr []*explore.ListLikedYouResponse_Liker
	for rows.Next() {
		var actorId explore.ListLikedYouResponse_Liker
		if err := rows.Scan(&actorId.ActorId); err != nil {
			return 0, err
		}
		rr = append(rr, &actorId)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return uint64(len(rr)), nil
}

func (s *StoreMySQLImpl) PutDecision(ctx context.Context, actorId int, recipientId int, liked int) (bool, error) {
	// Start transaction.
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}

	// Check if already exists by counting the number of skips between actor and recipient.
	// Decision to skip - not like.
	skips, err := countDecisions(tx, actorId, recipientId, 0)
	if err != nil {
		return false, err
	}
	// If skips exist, remove them.
	if skips > 0 {
		if err := deleteDecision(tx, ctx, actorId, recipientId, 0); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
			}
			return false, err
		}
	}

	// Check if already exists by counting the number of likes between actor and recipient.
	// Decision to like.
	likes, err := countDecisions(tx, actorId, recipientId, 1)
	if err != nil {
		return false, err
	}
	// If likes exist, remove them.
	if likes > 0 {
		if err = deleteDecision(tx, ctx, actorId, recipientId, 1); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
			}
			return false, err
		}
	}

	result, err := tx.ExecContext(ctx, "INSERT INTO decisions (ActorID, RecipientID, Liked) VALUES (?, ?, ?)", actorId, recipientId, liked)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("update drivers: unable to rollback: %v", rollbackErr)
		}
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rows != 1 {
		return false, fmt.Errorf("row != 1: %d", rows)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	// Exit early if the actor's decision is not to like the recipient.
	if liked == 0 {
		return false, nil
	}

	// See if mutual - checking both actor liking recipient and vice-versa.
	aCount, err := countDecisions(tx, actorId, recipientId, 1)
	if err != nil {
		return false, err
	}
	bCount, err := countDecisions(tx, recipientId, actorId, 1)
	if err != nil {
		return false, err
	}
	return aCount > 0 && bCount > 0, nil
}

// Count the number of decisions where actor likes recipient with liked.
func countDecisions(db *sql.Tx, actorId, recipientId, liked int) (int, error) {
	var result int
	err := db.QueryRow("SELECT COUNT(*) FROM decisions WHERE ActorID=? AND RecipientID=? AND Liked=?", actorId, recipientId, liked).Scan(&result)
	if err != nil {
		return 0, err
	}
	return result, err
}

// Delete all decisions from decisions where actor likes recipient with liked.
func deleteDecision(db *sql.Tx, ctx context.Context, actorId, recipientId, liked int) error {
	result, err := db.ExecContext(ctx, "DELETE FROM decisions WHERE ActorID=? AND RecipientID=? AND Liked=?", actorId, recipientId, liked)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("row != 1: %d", rows)
	}
	return nil
}
