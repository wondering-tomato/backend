package mysql

import (
	"backend/explore"
	"context"
	"fmt"
)

func (s *StoreMySQLImpl) PutDecision(ctx context.Context, actorId int, recipientId int, liked int) (bool, error) {
	// Check if already exists.
	var aa int
	err := s.db.QueryRow("SELECT COUNT(*) FROM decisions WHERE ActorID=? AND RecipientID=? AND Liked=0", actorId, recipientId).Scan(&aa)
	if err != nil {
		return false, err
	}
	if aa > 0 {
		result, err := s.db.ExecContext(ctx, "DELETE FROM decisions WHERE ActorID=? AND RecipientID=? AND Liked=0", actorId, recipientId)
		rows, err := result.RowsAffected()
		if err != nil {
			return false, err
		}
		if rows != 1 {
			return false, fmt.Errorf("row != 1: %d", rows)
		}
	}

	// Check if already exists.
	var bb int
	err = s.db.QueryRow("SELECT COUNT(*) FROM decisions WHERE ActorID=? AND RecipientID=? AND Liked=1", actorId, recipientId).Scan(&bb)
	if err != nil {
		return false, err
	}
	if bb > 0 {
		result, err := s.db.ExecContext(ctx, "DELETE FROM decisions WHERE ActorID=? AND RecipientID=? AND Liked=1", actorId, recipientId)
		rows, err := result.RowsAffected()
		if err != nil {
			return false, err
		}
		if rows != 1 {
			return false, fmt.Errorf("row != 1: %d", rows)
		}
	}

	result, err := s.db.ExecContext(ctx, "INSERT INTO decisions (ActorID, RecipientID, Liked) VALUES (?, ?, ?)", actorId, recipientId, liked)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rows != 1 {
		return false, fmt.Errorf("row != 1: %d", rows)
	}

	// Exit early if the actor's decision is not to like the recipient.
	if liked == 0 {
		return false, nil
	}

	// See if mutual.
	var aCount int
	err = s.db.QueryRow("SELECT COUNT(*) FROM decisions WHERE ActorID=? AND RecipientID=? AND Liked=1", actorId, recipientId).Scan(&aCount)
	if err != nil {
		return false, err
	}
	var bCount int
	err = s.db.QueryRow("SELECT COUNT(*) FROM decisions WHERE ActorID=? AND RecipientID=? AND Liked=1", recipientId, actorId).Scan(&bCount)
	if err != nil {
		return false, err
	}
	return aCount > 0 && bCount > 0, nil
}

// List all users who liked the recipient.
func (s *StoreMySQLImpl) GetAllLiked(ctx context.Context, id int) ([]*explore.ListLikedYouResponse_Liker, error) {

	rows, err := s.db.QueryContext(ctx, "SELECT ActorID FROM decisions WHERE RecipientID=? AND Liked=1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	var rr []*explore.ListLikedYouResponse_Liker
	for rows.Next() {
		var actorId explore.ListLikedYouResponse_Liker
		if err := rows.Scan(&actorId.ActorId); err != nil {
			return nil, fmt.Errorf("actorId %d: %v", id, err)
		}
		rr = append(rr, &actorId)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("actorId %d: %v", id, err)
	}
	return rr, nil
}

func (s *StoreMySQLImpl) GetNewAllLiked(ctx context.Context, id int) ([]*explore.ListLikedYouResponse_Liker, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT ActorID FROM decisions WHERE RecipientID=? AND Liked=1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	aa := make(map[*explore.ListLikedYouResponse_Liker]bool)

	for rows.Next() {
		var actorId explore.ListLikedYouResponse_Liker
		if err := rows.Scan(&actorId.ActorId); err != nil {
			return nil, fmt.Errorf("actorId %d: %v", id, err)
		}
		aa[&actorId] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("actorId %d: %v", id, err)
	}

	//
	rows, err = s.db.QueryContext(ctx, "SELECT RecipientID FROM decisions WHERE ActorID=? AND Liked=1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rr []*explore.ListLikedYouResponse_Liker
	for rows.Next() {
		var recipientID explore.ListLikedYouResponse_Liker
		if err := rows.Scan(&recipientID.ActorId); err != nil {
			return nil, fmt.Errorf("recipientID %d: %v", id, err)
		}
		if _, ok := aa[&recipientID]; !ok {
			rr = append(rr, &recipientID)
			aa[&recipientID] = true
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("actorId %d: %v", id, err)
	}
	return rr, nil
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
			return 0, fmt.Errorf("actorId %d: %v", id, err)
		}
		rr = append(rr, &actorId)
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("actorId %d: %v", id, err)
	}
	return uint64(len(rr)), nil
}
