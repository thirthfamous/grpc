package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"

	pb "gitlab.com/jonathannobi/go/grpc/transaction"
	"google.golang.org/grpc"
)

type TransactionServer struct {
	pb.UnimplementedTransactionsServer
}

func (d *TransactionServer) CreateTransaction(_ context.Context, trx *pb.Transaction) (*pb.Response, error) {
	log.Printf("Received: %v", trx.GetTitle())

	panic("damn i got panic") // made up error
	// Send Response to client
	return &pb.Response{Body: trx.GetTitle() + trx.GetBody()}, nil
}

// every process will be through this function first
func middlewarePrintHello(ctx context.Context) (context.Context, error) {
	fmt.Println("this middleware called because it is neccesary")
	newCtx := context.WithValue(ctx, "tokenInfo", "123")
	return newCtx, nil
}

// every panic will be recovered here
func recoveryHandler(p interface{}) (err error) {
	fmt.Println(p)
	return status.Errorf(codes.Unknown, "panic triggered: %v", p)
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryHandler),
	}

	grpcServer := grpc.NewServer(grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_auth.StreamServerInterceptor(middlewarePrintHello),
		grpc_recovery.StreamServerInterceptor(opts...),
	)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(middlewarePrintHello),
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
		),
	)

	pb.RegisterTransactionsServer(grpcServer, &TransactionServer{})

	log.Printf("server listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port 9000: %v", err)
	}

}
