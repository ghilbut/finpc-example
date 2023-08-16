package grpc

import (
	"context"
	"database/sql"
	"encoding/hex"

	// external packages
	"github.com/getsentry/sentry-go"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	DBSession string = "dbSession"
	tpKey     string = "traceparent"
)

func SentryStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
		}

		span := sentry.StartTransaction(ctx, info.FullMethod, func(s *sentry.Span) {
			s.Name = "finpc-server"
			s.Op = "grpc.server"
			s.Description = info.FullMethod

			traceParent := metadata.ValueFromIncomingContext(ctx, tpKey)
			if traceParent != nil && len(traceParent) != 0 {
				log.Info("[gRPC][Metadata] ", tpKey, ": ", traceParent[0])

				_, err := hex.Decode(s.TraceID[:], []byte(traceParent[0]))
				if err != nil {
					sentry.CaptureException(err)
				}
			}
		})

		ctx = span.Context()
		defer span.Finish()

		stream := grpc_middleware.WrapServerStream(ss)
		stream.WrappedContext = ctx

		err := handler(srv, stream)
		if err != nil {
			span.Status = toSentrySpanStatus(err)
		}

		return err
	}
}

func SentryUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
		}

		span := sentry.StartTransaction(ctx, info.FullMethod, func(s *sentry.Span) {
			s.Name = "finpc-server"
			s.Op = "grpc.server"
			s.Description = info.FullMethod

			traceParent := metadata.ValueFromIncomingContext(ctx, tpKey)
			if traceParent != nil && len(traceParent) != 0 {
				log.Info("[gRPC][Metadata] ", tpKey, ": ", traceParent[0])

				_, err := hex.Decode(s.ParentSpanID[:], []byte(traceParent[0]))
				if err != nil {
					sentry.CaptureException(err)
				}
			}
		})

		ctx = span.Context()
		defer span.Finish()

		hub.Scope().SetExtra("requestBody", req)

		resp, err := handler(ctx, req)
		if err != nil {
			log.Error(err)
			span.Status = toSentrySpanStatus(err)
		}

		return resp, err
	}
}

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
			SentryStreamInterceptor(),
		),
		grpc.ChainUnaryInterceptor(
			SentryUnaryServerInterceptor(),
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
