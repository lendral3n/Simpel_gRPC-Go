package main

import (
	"l3ngrpc/cmd/config"
	"l3ngrpc/cmd/services"
	productPb "l3ngrpc/pb/product"
	
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	netListen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err.Error())
	}

	db := config.ConnectDatabase()

	grpcServer := grpc.NewServer()
	productService := services.ProductService{DB: db}
	productPb.RegisterProductServiceServer(grpcServer, productService)
	
	log.Printf("Server STarted At %v", netListen.Addr())
	if err := grpcServer.Serve(netListen); err != nil {
		log.Fatalf("Failed to serve %v", err.Error())
	}
}
