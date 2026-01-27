package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/Pavel7004/WebShop/pkg/pb/inventory"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewInventoryServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Вызов ReserveItems
	resReq := &pb.ReserveItemsRequest{
		Items: []*pb.InventoryItem{
			{
				ItemId:   "1111",
				Quantity: 22,
			},
		},
	}
	res, err := c.ReserveItems(ctx, resReq)
	if err != nil {
		log.Fatalf("could not reserve items: %v", err)
	}
	log.Printf("ReserveItems Response: %v", res.Success)

	// Вызов CancelReserve
	canReq := &pb.CancelReserveRequest{
		Items: []*pb.InventoryItem{{
			ItemId:   "1111",
			Quantity: 2,
		}},
	}
	canRes, err := c.CancelReserve(ctx, canReq)
	if err != nil {
		log.Fatalf("could not cancel reservation: %v", err)
	}
	log.Printf("CancelReserve Response: %v", canRes.Success)
}
