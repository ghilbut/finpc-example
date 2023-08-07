package grpc

import (
	"context"
	"crypto/tls"
	"database/sql"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	// external packages
	log "github.com/sirupsen/logrus"
	"github.com/speps/go-hashids"
	"google.golang.org/grpc"
)

const (
	DBSession string = "dbSession"
	HashID    string = "hashID"
)

func DBUnaryServerInterceptor(session *sql.DB) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(context.WithValue(ctx, DBSession, session), req)
	}
}

func HashUnaryServerInterceptor(h *hashids.HashID) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(context.WithValue(ctx, HashID, h), req)
	}
}

func LoadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair("server-cert.pem", "server-key.pem")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func NewGrpcServer(db *sql.DB) *grpc.Server {

	hd := hashids.NewData()
	hd.Salt = "salt"
	hd.MinLength = 7
	h, err := hashids.NewWithData(hd)
	if err != nil {
		log.Fatal(err)
	}

	// creds, err := LoadTLSCredentials()
	// if err != nil {
	//     log.Fatal(err)
	// }
	creds := insecure.NewCredentials()
	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(
			DBUnaryServerInterceptor(db),
			HashUnaryServerInterceptor(h),
		),
	)

	trading := Trading{}
	RegisterTradingServer(grpcServer, &trading)

	board := Board{}
	RegisterBoardServer(grpcServer, &board)

	return grpcServer
}
