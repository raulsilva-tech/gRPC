package main

import (
	"database/sql"
	"net"

	_ "github.com/mattn/go-sqlite3"
	"github.com/raulsilva-tech/gRPC/internal/database"
	"github.com/raulsilva-tech/gRPC/internal/pb"
	"github.com/raulsilva-tech/gRPC/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	categoryDB := database.NewCategory(db)
	categoryService := service.NewCategoryService(*categoryDB)

	gRPCServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(gRPCServer, categoryService)
	reflection.Register(gRPCServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	if err = gRPCServer.Serve(lis); err != nil {
		panic(err)
	}
}
