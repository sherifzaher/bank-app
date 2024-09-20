package main

import (
	"context"
	"database/sql"
	"github.com/rakyll/statik/fs"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/sherifzaher/clone-simplebank/api"
	"github.com/sherifzaher/clone-simplebank/db/sqlc"
	_ "github.com/sherifzaher/clone-simplebank/doc/statik"
	"github.com/sherifzaher/clone-simplebank/gapi"
	"github.com/sherifzaher/clone-simplebank/pb"
	"github.com/sherifzaher/clone-simplebank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load env:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to DB:", err)
		return
	}

	if err = conn.Ping(); err != nil {
		log.Fatal("DB is not live: ", err)
	}

	store := db.NewStore(conn)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik fs:", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatalf("start HTTP gateway server at %s:", listener.Addr().String())
	}

	log.Printf("HTTP gateway server running at: %s", config.HttpServerAddress)

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server")
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatalf("start gRPC server at %s:", listener.Addr().String())
	}

	log.Printf("gRPC server running at: %s", config.GrpcServerAddress)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("error during run the server %v", err)
		return
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		panic(err)
	}
}
