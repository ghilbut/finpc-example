package grpc

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/hex"
	// external packages
	"github.com/getsentry/sentry-go"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	log "github.com/sirupsen/logrus"
	"github.com/speps/go-hashids"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	DBSession string = "dbSession"
	HashID    string = "hashID"
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
			//s.Name = "finpc-server"
			s.Op = "grpc.server"
			s.Description = info.FullMethod

			traceId := metadata.ValueFromIncomingContext(ctx, "traceid")
			if traceId != nil && len(traceId) != 0 {
				_, err := hex.Decode(s.TraceID[:], []byte(traceId[0]))
				if err != nil {
					sentry.CaptureException(err)
				}
			}

			spanId := metadata.ValueFromIncomingContext(ctx, "spanid")
			if spanId != nil && len(spanId) != 0 {
				_, err := hex.Decode(s.SpanID[:], []byte(spanId[0]))
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

			traceId := metadata.ValueFromIncomingContext(ctx, "traceid")
			if traceId != nil && len(traceId) != 0 {
				_, err := hex.Decode(s.TraceID[:], []byte(traceId[0]))
				if err != nil {
					sentry.CaptureException(err)
				}
			}

			spanId := metadata.ValueFromIncomingContext(ctx, "spanid")
			if spanId != nil && len(spanId) != 0 {
				_, err := hex.Decode(s.ParentSpanID[:], []byte(spanId[0]))
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
		sentry.CaptureMessage(err.Error())
		log.Fatal(err)
	}

	// creds, err := LoadTLSCredentials()
	// if err != nil {
	//     log.Fatal(err)
	// }
	creds := insecure.NewCredentials()
	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc.StreamInterceptor(
			SentryStreamInterceptor(),
		),
		grpc.ChainUnaryInterceptor(
			SentryUnaryServerInterceptor(),
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
