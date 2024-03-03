package main

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"service-profile/system"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "service-profile/proto"
)

type server struct {
	pb.UnimplementedProfileServiceServer
	storage system.Storage
}

func main() {
	// Set up the logger
	system.InitLogger()

	// Connect to the database
	storage, err := system.NewStorage()
	if err != nil {
		slog.Error("Error opening database", "db.Connect", err)
		panic(err)
	}
	slog.Info("Database connected")

	// Run migrations
	err = storage.Migrations()
	if err != nil {
		slog.Error("Error running migrations", "db.Migrations", err)
		panic(err)
	}
	slog.Info("Migrations completed")

	// Run the gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", system.GRPC_PORT))
	if err != nil {
		slog.Error("Error listening on gRPC port", "net.Listen", err)
		panic(err)
	}
	var s *grpc.Server
	if system.TLS == "true" {
		certificate, err := tls.LoadX509KeyPair(system.CERT_PATH, system.KEY_PATH)
		if err != nil {
			slog.Error("Error loading TLS certificate", "tls.LoadX509KeyPair", err)
			panic(err)
		}
		s = grpc.NewServer(grpc.Creds(credentials.NewServerTLSFromCert(&certificate)))
	} else {
		s = grpc.NewServer()
	}
	pb.RegisterProfileServiceServer(s, &server{
		storage: storage,
	})
	slog.Info("gRPC server listening on", "port", system.GRPC_PORT)
	err = s.Serve(lis)
	if err != nil {
		slog.Error("Error serving gRPC", "s.Serve", err)
		panic(err)
	}
}
