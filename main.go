package main

import (
	"aboutMeStoreService/configuration"
	"aboutMeStoreService/domain/repository/migrations"
	s "aboutMeStoreService/service"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

func main() {

	migrator := migrations.New(
		configuration.DbConnConfiguration.DriverName,
		configuration.DbConnConfiguration.DataSourceName,
		"./domain/repository/migrations")

	migrator.UpToLastVersion()

	migrator.Close()

	server := grpc.NewServer()

	service, err := s.NewDialogService(s.WithLocalRepository())

	if err != nil {
		log.Fatal(err)
	}

	s.RegisterDialogServiceServer(server, service)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatal(err)
	}

	if err = server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
