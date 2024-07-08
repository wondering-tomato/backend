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

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := []*explore.ListLikedYouResponse_Liker{}
	// Pagination.
	pageToken := ""
	for {
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

	// Contact the server and print out its response.
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ListNewLikedYou(ctx, &explore.ListLikedYouRequest{
		RecipientUserId: "2",
		PaginationToken: new(string),
	})
	if err != nil {
		log.Fatalf("Error ListNewLikedYou: %v", err)
	}
	log.Printf("ListNewLikedYou: %s", r.GetLikers())

	// Contact the server and print out its response.
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	aa, err := c.CountLikedYou(ctx, &explore.CountLikedYouRequest{
		RecipientUserId: "2",
	})
	if err != nil {
		log.Fatalf("Error CountLikedYou: %v", err)
	}
	log.Printf("CountLikedYou: %d", aa.GetCount())

	// Contact the server and print out its response.
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rrr, err := c.PutDecision(ctx, &explore.PutDecisionRequest{
		ActorUserId:     "4",
		RecipientUserId: "3",
		LikedRecipient:  true,
	})
	if err != nil {
		log.Fatalf("Error PutDecision: %v", err)
	}
	log.Printf("PutDecision: %v", rrr.GetMutualLikes())
}
