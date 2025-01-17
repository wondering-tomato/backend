package main

import (
	"backend/explore"
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "0.0.0.0:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := explore.NewExploreServiceClient(conn)

	// ##
	// Make a paginated call to ListLikedYou.
	// Append the results to res and then log res.
	// ##
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := []*explore.ListLikedYouResponse_Liker{}
	pageToken := ""
	for {
		// Pagination.
		r, err := c.ListLikedYou(ctx, &explore.ListLikedYouRequest{
			RecipientUserId: "2",
			PaginationToken: &pageToken,
		})
		if err != nil {
			log.Fatalf("Error ListLikedYou: %v", err)
		}
		log.Printf("ListLikedYou: %s", r.GetLikers())
		res = append(res, r.GetLikers()...)

		pageToken = r.GetNextPaginationToken()
		if pageToken == "" {
			break
		}
		log.Printf("next page token: %s", pageToken)
	}
	log.Printf("Full ListLikedYou: %s", res)

	// ##
	// Make a paginated call to ListNewLikedYou.
	// Append the results to res and then log res.
	// ##
	ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res = []*explore.ListLikedYouResponse_Liker{}
	pageToken = ""
	for {
		r, err := c.ListNewLikedYou(ctx, &explore.ListLikedYouRequest{
			RecipientUserId: "2",
			PaginationToken: &pageToken,
		})
		if err != nil {
			log.Fatalf("Error ListNewLikedYou: %v", err)
		}
		log.Printf("ListNewLikedYou: %s", r.GetLikers())
		res = append(res, r.GetLikers()...)

		pageToken = r.GetNextPaginationToken()
		if pageToken == "" {
			break
		}
		log.Printf("next page token: %s", pageToken)
	}
	log.Printf("Full ListNewLikedYou: %s", res)

	// ##
	// Make a call to CountLikedYou.
	// ##
	ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	aa, err := c.CountLikedYou(ctx, &explore.CountLikedYouRequest{
		RecipientUserId: "2",
	})
	if err != nil {
		log.Fatalf("Error CountLikedYou: %v", err)
	}
	log.Printf("CountLikedYou: %d", aa.GetCount())

	// ##
	// Make a call to PutDecision.
	// ##
	ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	rrr, err := c.PutDecision(ctx, &explore.PutDecisionRequest{
		ActorUserId:     "4",
		RecipientUserId: "3",
		LikedRecipient:  false,
	})
	if err != nil {
		log.Fatalf("Error PutDecision: %v", err)
	}
	log.Printf("PutDecision: %v", rrr.GetMutualLikes())
}
