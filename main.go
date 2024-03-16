package main

import (
	"fmt"
	"log"
	pb "l3n/grpc/pb"

	"google.golang.org/protobuf/proto"
)

func main() {

	products := &pb.Products{
		Data: []*pb.Product{
			{
				Id: 1,
				Name: "Lendra",
				Price: 10.000,
				Stock: 20,
				Category: &pb.Category{
					Id: 1,
					Name: "Manusia",
				},
			},
			{
				Id: 2,
				Name: "Oppo",
				Price: 1000.000,
				Stock: 20,
				Category: &pb.Category{
					Id: 2,
					Name: "Handphone",
				},
			},
		},
	}

	encodeData, err := proto.Marshal(products)
	if err != nil {
		log.Fatal("Error Marshal", err)
	}

	// compact binary wire format
	fmt.Println(encodeData)

	resultProduct := &pb.Products{}
	if err = proto.Unmarshal(encodeData, resultProduct); err != nil {
		log.Fatal("Error Unmarshal", err)
	}

	for _, product := range resultProduct.GetData() {
		fmt.Println(product.GetName()) // get name product
		fmt.Println(product.Category.GetName()) // get name category product
	}
}