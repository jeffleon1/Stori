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
	gRPCServer()
}

func gRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Config.GrpcPort))
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
	log.Printf("gRPC server started on port %s", config.Config.GrpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}

func InitMail() (*mail.SMTPClient, error) {
	server := mail.NewSMTPClient()
	server.Port = config.Config.EmailPort
	server.Host = config.Config.EmailHost
	server.Password = config.Config.EmailPassword
	server.Username = config.Config.EmailUserName
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
