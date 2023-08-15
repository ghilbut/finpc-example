package grpc

import (
	"context"
	"database/sql"

	// external packages
	"github.com/getsentry/sentry-go"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
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
			DBUnaryServerInterceptor(db),
		),
	)

	RegisterBoardServer(grpcServer, &Board{})

	return grpcServer
}

func toSentrySpanStatus(err error) sentry.SpanStatus {
	code := status.Code(err)

	switch code {
	case codes.Internal:
		return sentry.SpanStatusInternalError
	case codes.InvalidArgument:
		return sentry.SpanStatusInvalidArgument
	case codes.FailedPrecondition:
		return sentry.SpanStatusFailedPrecondition
	default:
		return sentry.SpanStatusUndefined
	}
}
