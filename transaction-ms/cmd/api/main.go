package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jeffleon1/transaction-ms/pkg/health"
	grpcClientDomain "github.com/jeffleon1/transaction-ms/pkg/mail/domain"
	pb "github.com/jeffleon1/transaction-ms/pkg/mail/domain/proto"
	grpcClientRepo "github.com/jeffleon1/transaction-ms/pkg/mail/infrastructure"
	"github.com/jeffleon1/transaction-ms/pkg/router"
	"github.com/jeffleon1/transaction-ms/pkg/swagger"
	service "github.com/jeffleon1/transaction-ms/pkg/transactions/application"
	"github.com/jeffleon1/transaction-ms/pkg/transactions/domain"
	"github.com/jeffleon1/transaction-ms/pkg/transactions/infrastructure"
	"github.com/jeffleon1/transaction-ms/pkg/transactions/infrastructure/repository"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, grpcMailRepo, err := gRPCClientRepo()
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	defer conn.Close()

	transactionRoutes := InitTransactionPackage(*grpcMailRepo)
	r := router.NewRouter(router.RoutesGroup{
		Health:      health.NewHealthCheckRoutes(),
		Swagger:     swagger.NewSwaggerDocsRoutes(),
		Transaction: transactionRoutes,
	})

	logrus.Fatal(r.Run(fmt.Sprintf(":%s", "8080")))
}

func InitTransactionPackage(grpcClientRepo grpcClientDomain.GrpcMailRepository) *infrastructure.TransactionRoutes {
	proccesorRepo := repository.NewProccesorRepository()
	accountRepo := InitMongoDB()
	transactionService := service.NewTransactionService(proccesorRepo, accountRepo, grpcClientRepo)
	transactionHandler := infrastructure.NewTransactionHandler(transactionService)
	return infrastructure.NewRoutes(transactionHandler)
}

func gRPCClientRepo() (*grpc.ClientConn, *grpcClientDomain.GrpcMailRepository, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "5001"), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, nil, err
	}
	c := pb.NewMailServiceClient(conn)
	client := grpcClientRepo.NewGrpcMailClient(c)
	return conn, &client, nil
}

func InitMongoDB() domain.AccountRepository {
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf(
		"mongodb://%s:%s@%s:27017",
		"stori",
		"stori",
		"localhost",
	)))
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("stori")
	return repository.NewAccountRepository(db)
}
