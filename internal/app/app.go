package app

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/epic55/BankAppNew/internal/middleware"
	"github.com/epic55/BankAppNew/internal/repository"
	"github.com/epic55/BankAppNew/internal/services"
	"github.com/epic55/BankAppNew/internal/transport"
	"github.com/epic55/BankAppNew/pkg/db"

	pb "github.com/epic55/BankAppNew/buyingGRPC"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func Run() {
	cfg := LoadConfig()
	dbConn := db.Connect(cfg.DB.DSN)

	userRepo := repository.NewRepository(dbConn)
	userService := services.NewService(userRepo)
	userHandler := transport.NewHandler(userService)

	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)
	userHandler.RegisterRoutes(router)

	//GRPC
	var wg sync.WaitGroup
	grpcServer := grpc.NewServer()
	pb.RegisterBuyingServer(grpcServer, userService) //userService) //&services.ServiceStruct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", "localhost:8081")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		fmt.Println("GRPC server is running on port 8081...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Server running on :%s", cfg.Server.Port)
		log.Fatal(http.ListenAndServe("localhost:"+cfg.Server.Port, router))
	}()

	//grpcServer.GracefulStop()

	wg.Wait()
	log.Println("Servers gracefully stopped.")
}
