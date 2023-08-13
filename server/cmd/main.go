package main

import (
    "database/sql"
    "fmt"
    "net"
    "os"
    "path/filepath"
    "sync"
    "time"

    // external packages
    "github.com/getsentry/sentry-go"
    sentryfasthttp "github.com/getsentry/sentry-go/fasthttp"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
    log "github.com/sirupsen/logrus"
    "github.com/valyala/fasthttp"

    // project packages
    . "github.com/ghilbut/finpc/grpc"
    . "github.com/ghilbut/finpc/rest"
)

func init() {
    basepath, err := os.Getwd()
    if err != nil {
        log.Println(err)
    }

    if env := getEnvValue("ENVIRONMENT", ""); env != "" {
        loadEnvs(basepath, ".env."+env+".local")
        loadEnvs(basepath, ".env."+env)
    }
    loadEnvs(basepath, ".env.local")
    loadEnvs(basepath, ".env")
}

func loadEnvs(basepath, name string) {
    envpath := filepath.Join(basepath, name)
    if _, err := os.Stat(envpath); err != nil {
        return
    }
    if err := godotenv.Load(envpath); err != nil {
        log.Fatalf("Error loading %s file", name)
    }
}

func main() {
    log.SetLevel(log.TraceLevel)

    hostname := getEnvValue("HOSTNAME", "")
    if hostname == "" {
        log.Fatal("you should set HOSTNAME for sentry")
    }
    release := getEnvValue("RELEASE", "localhost")
    environment := getEnvValue("ENVIRONMENT", "localhost")

    err := sentry.Init(sentry.ClientOptions{
        Dsn: "https://7caebc6d29f76e2b8f02632b833a0004@o4505689921683456.ingest.sentry.io/4505690170130432",
        //Debug:            true,
        //AttachStacktrace: true,
        SampleRate:       1.0,
        EnableTracing:    true,
        TracesSampleRate: 1.0,
        TracesSampler: sentry.TracesSampler(func(ctx sentry.SamplingContext) float64 {
            log.Trace(ctx)
            return 1.0
        }),
        //BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
        //    if hint.Context != nil {
        //        if ctx, ok := hint.Context.Value(sentry.RequestContextKey).(*fasthttp.RequestCtx); ok {
        //            log.Trace(string(ctx.Request.Host()))
        //        }
        //    }
        //    return event
        //},
        //BeforeSendTransaction: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
        //    log.Trace(event)
        //    log.Trace(hint)
        //    return event
        //},
        ServerName:  hostname,
        Release:     release,
        Environment: environment,
    })
    if err != nil {
        log.Fatalf("failed to initialize sentry: %v", err)
    }
    defer sentry.Flush(2 * time.Second)

    db, err := OpenDatabase()
    if err != nil {
        sentry.CaptureMessage(err.Error())
        log.Fatal(err)
    }
    defer db.Close()

    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        port := 8080
        addr := fmt.Sprintf(":%d", port)

        rest := NewRestServer()
        sentryHandler := sentryfasthttp.New(sentryfasthttp.Options{})
        fastHTTPHandler := sentryHandler.Handle(rest.Handler)

        log.Printf("run RESTful server on port %d", port)

        if err := fasthttp.ListenAndServe(addr, fastHTTPHandler); err != nil {
            sentry.CaptureMessage(err.Error())
            log.Fatalf("failed to run RESTful server: %v", err)
        }
    }()

    go func() {
        port := 9095
        addr := fmt.Sprintf(":%d", port)

        listen, err := net.Listen("tcp4", addr)
        if err != nil {
            sentry.CaptureMessage(err.Error())
            log.Fatalf("failed to listen: %v", err)
        }

        grpc := NewGrpcServer(db)

        log.Printf("run gRPC server on port %d", port)
        if err := grpc.Serve(listen); err != nil {
            sentry.CaptureMessage(err.Error())
            log.Fatalf("failed to run gRPC server: %v", err)
        }
    }()

    wg.Wait()
}

func OpenDatabase() (*sql.DB, error) {

    host := getEnvValue("PG_HOST", "localhost")
    port := getEnvValue("PG_PORT", "5432")
    user := getEnvValue("PG_USER", "postgres")
    pw := getEnvValue("PG_PASSWORD", "postgrespw")
    db := getEnvValue("PG_DATABASE", "postgres")
    ssl := getEnvValue("PG_SSLMODE", "disabled")

    log.Infoln("PG_HOST: ", host)
    log.Infoln("PG_PORT: ", port)
    log.Infoln("PG_USER: ", user)
    log.Infoln("PG_DATABASE: ", db)
    log.Infoln("PG_SSLMODE: ", ssl)

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pw, db, ssl)
    return sql.Open("postgres", dsn)
}

func getEnvValue(name, defaultValue string) string {
    v := os.Getenv(name)
    if v == "" {
        return defaultValue
    }
    return v
}
