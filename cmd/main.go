package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"imman/parser_service/internal/service"
	"imman/parser_service/internal/service/repository"
	"imman/parser_service/internal/service/webapi"
	"imman/parser_service/protos/protos/parser_pb"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"google.golang.org/grpc"
)

var (
	port string
)

func main() {
	flag.StringVar(&port, "port", ":8001", "port")
	flag.Parse()

	config, err := pgx.ParseConfig("postgres://fprotimaru:1@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		panic(err)
	}
	config.PreferSimpleProtocol = true

	sqldb := stdlib.OpenDB(*config)
	db := bun.NewDB(sqldb, pgdialect.New())

	repo := repository.NewPostRepository(db)
	webAPI := webapi.NewPostWebAPI(&http.Client{})

	uc := service.NewPostService(repo, webAPI)

	server := grpc.NewServer()
	parser_pb.RegisterPostParserServiceServer(server, uc)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err = server.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println("post_parser gRPC service is running on", port)
	<-quit
	log.Println("stopping post_parser gRPC service...")
	server.GracefulStop()
	log.Println("stopped post_parser gRPC service")
}
