package main

import (
	"fmt"
	"log"
	"net"
	"time"

	config "github.com/jeffleon1/email-ms/internal"
	pb "github.com/jeffleon1/email-ms/pkg/mail/domain/proto"
	infra "github.com/jeffleon1/email-ms/pkg/mail/infrastructure"
	mail "github.com/xhit/go-simple-mail/v2"
	"google.golang.org/grpc"
)

func main() {
	config.InitEnvConfigs()
	gRPCServer()
}

func gRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.EnvConfigs.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	SMTPClient, err := InitMail()
	if err != nil {
		log.Fatalf("Failed to create smtp client: %v", err)
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterMailServiceServer(s, &infra.Mailserver{SMTPClient: SMTPClient})
	log.Printf("gRPC server started on port %s", config.EnvConfigs.GRPCPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}

func InitMail() (*mail.SMTPClient, error) {
	server := mail.NewSMTPClient()
	server.Port = config.EnvConfigs.EmailPort
	server.Host = config.EnvConfigs.EmailHost
	server.Password = config.EnvConfigs.EmailPassword
	server.Username = config.EnvConfigs.EmailUserName
	server.Encryption = mail.EncryptionNone
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.KeepAlive = true

	smtpClient, err := server.Connect()
	if err != nil {
		return nil, err
	}
	return smtpClient, nil
}
