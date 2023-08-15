package grpc

import (
	"context"
	"database/sql"

	// external packages
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	DBSession string = "dbSession"
)

func DBUnaryServerInterceptor(session *sql.DB) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(context.WithValue(ctx, DBSession, session), req)
	}
}

func NewGrpcServer(db *sql.DB) *grpc.Server {

	creds := insecure.NewCredentials()
	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainStreamInterceptor(
			otelgrpc.StreamServerInterceptor(),
		),
		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(),
			//metadataMiddleware(),
			DBUnaryServerInterceptor(db),
		),
	)

	RegisterBoardServer(grpcServer, &Board{})

	return grpcServer
}

func metadataMiddleware() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		const tpKey = "traceparent"

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			log.Info("[gRPC][Metadata] ", tpKey, ": ", md.Get(tpKey))
		}

		return handler(context.WithValue(ctx, tpKey, md.Get(tpKey)), req)
	}
}
