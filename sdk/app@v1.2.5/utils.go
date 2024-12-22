package app

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func NewHTTPSrvWorker(
	httpSrv *echo.Echo,
	address string,
	shutdownTimeout time.Duration,
) Worker {
	return WorkerFn(func(ctx context.Context) error {
		group, ctxGroup := errgroup.WithContext(ctx)
		// run observer for shutdown http server
		group.Go(func() error {
			<-ctxGroup.Done()
			ctxShutdown, cancelShutdown := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
			defer cancelShutdown()
			err := httpSrv.Shutdown(ctxShutdown)
			if err != nil {
				log.Error(ctx, "shutdown http server: %s", err.Error())
				return err
			}
			return nil
		})
		// run http server
		group.Go(func() error {
			err := httpSrv.Start(address)
			if err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					log.Info(ctx, "http server has gracefully shutdown (%s).", err)
					return nil
				}
				log.Error(ctx, "http server abnormal shutdown (%s).", err)
				return fmt.Errorf("start http server: %w", err)
			}
			return nil
		})
		return group.Wait()
	})
}

// NewGRPCSrvWorker ...
// Deprecated: please use [WithGRPCServer] instead of [WithGRPC] and NewGRPCSrvWorker combo.
func NewGRPCSrvWorker(
	grpcSrv *grpc.Server,
	address string,
) Worker {
	return WorkerFn(func(ctx context.Context) error {
		group, ctxGroup := errgroup.WithContext(ctx)
		// run observer for shutdown grpc server
		group.Go(func() error {
			<-ctxGroup.Done()
			grpcSrv.GracefulStop()
			return nil
		})
		// run grpc server
		group.Go(func() error {
			listener, err := net.Listen("tcp", address)
			if err != nil {
				return fmt.Errorf("open listen %s: %w", address, err)
			}
			err = grpcSrv.Serve(listener)
			if err != nil {
				if errors.Is(err, grpc.ErrServerStopped) {
					log.Info(ctx, "grpc server has gracefully shutdown (return: %s).", err)
					return nil
				}
				log.Error(ctx, "shutdown grpc server: %s", err.Error())
			}
			log.Info(ctx, "grpc server has gracefully shutdown.")
			return err
		})
		return group.Wait()
	},
	)
}
