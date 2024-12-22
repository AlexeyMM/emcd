package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"code.emcdtech.com/emcd/sdk/app"
	sdkError "code.emcdtech.com/emcd/sdk/error"
	"code.emcdtech.com/emcd/service/whitelabel/internal/config"
	grpcController "code.emcdtech.com/emcd/service/whitelabel/internal/controller/grpc"
	healthchecker "code.emcdtech.com/emcd/service/whitelabel/internal/heakthchecker"
	"code.emcdtech.com/emcd/service/whitelabel/internal/repository"
	s "code.emcdtech.com/emcd/service/whitelabel/internal/service"
	pb "code.emcdtech.com/emcd/service/whitelabel/protocol/whitelabel"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	if err = InitLogger(cfg.Log.Level); err != nil {
		log.Fatal().Err(err).Send()
	}

	db, err := cfg.PGXPool.New(ctx)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	repo := repository.NewWhiteLabel(db)
	service := s.NewWhiteLabel(repo)
	whiteLabelServer := grpcController.NewWhiteLabel(service, cfg.Languages)

	var (
		httpSrv *echo.Echo
		grpcSrv *grpc.Server
	)

	httpSrv = echo.New()
	grpcFactory := func() func(opts ...grpc.ServerOption) *grpc.Server {
		return func(opts ...grpc.ServerOption) *grpc.Server {
			opts = append(opts, grpc.UnaryInterceptor(sdkError.ServerUnaryInterceptor))
			grpcSrv = grpc.NewServer(opts...)
			pb.RegisterWhitelabelServiceServer(grpcSrv, whiteLabelServer)
			reflection.Register(grpcSrv)
			return grpcSrv
		}
	}

	if err := app.New().
		WithPprof(httpSrv).
		WithGRPC(grpcFactory(), healthchecker.NewCommon(db)).
		WithWorker(
			app.WorkerFn(func(ctx context.Context) error {
				var listener net.Listener
				listener, err = net.Listen("tcp", cfg.GRPC.ListenAddr)
				if err != nil {
					return fmt.Errorf("open listen %s: %w", cfg.GRPC.ListenAddr, err)
				}
				go func() {
					<-ctx.Done()
					grpcSrv.GracefulStop()
				}()
				return grpcSrv.Serve(listener)
			}),
			app.WorkerFn(
				func(ctx context.Context) error {
					go func() {
						<-ctx.Done()
						_ = httpSrv.Shutdown(context.Background())
					}()
					return httpSrv.Start(cfg.HTTP.ListenAddr)
				},
			),
		).Run(ctx); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func InitLogger(level string) error {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Caller().Stack().Logger()
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return errors.WithStack(err)
	}
	zerolog.SetGlobalLevel(logLevel)
	return nil
}
