package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	pb "gitlab.com/jonathannobi/go/grpc/clienttransaction"
	"google.golang.org/grpc"
)

const (
	address = "localhost:9000"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithCredentialsBundle(insecure.NewBundle()), grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewTransactionsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.CreateTransaction(ctx, &pb.Transaction{
		Title:  "join transaction",
		Body:   "body transaction",
		Amount: 900,
	})

	if err != nil {
		log.Fatalf("could not create transaction: %v", err.Error())
	}

	log.Printf(`Response : %v`, r.GetBody())
}
