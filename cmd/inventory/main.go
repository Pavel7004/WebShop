package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/Pavel7004/WebShop/pkg/pb/inventory"
)

type InventoryServer struct {
	pb.UnimplementedInventoryServiceServer

	data map[string]int32
}

func (s *InventoryServer) ReserveItems(
	ctx context.Context,
	in *pb.ReserveItemsRequest,
) (*pb.ReserveItemsResponse, error) {
	for _, el := range in.GetItems() {
		if kol, ok := s.data[el.GetItemId()]; ok {
			if kol >= el.GetQuantity() {
				s.data[el.GetItemId()] -= el.GetQuantity()
				fmt.Printf("Item %s - %d\n", el.GetItemId(), s.data[el.GetItemId()])
			} else {
				return &pb.ReserveItemsResponse{
					Success: false,
				}, errors.New("Not enough items")
			}
		}
	}

	return &pb.ReserveItemsResponse{
		Success: true,
	}, nil
}

func (s *InventoryServer) CancelReserve(
	ctx context.Context,
	in *pb.CancelReserveRequest,
) (*pb.CancelReserveResponse, error) {

	for _, el := range in.GetItems() {
		if _, ok := s.data[el.GetItemId()]; ok {
			s.data[el.GetItemId()] += el.GetQuantity()
			fmt.Printf("Item %s - %d\n", el.GetItemId(), s.data[el.GetItemId()])
		}
	}

	return &pb.CancelReserveResponse{
		Success: true,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	invServer := &InventoryServer{
		data: map[string]int32{
			"1111": 10,
		},
	}

	pb.RegisterInventoryServiceServer(s, invServer)

	log.Println("Server is running on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
