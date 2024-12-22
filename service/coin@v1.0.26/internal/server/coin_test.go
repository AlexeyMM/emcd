package server_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	errors "code.emcdtech.com/emcd/sdk/error"

	pb "code.emcdtech.com/emcd/service/coin/protocol/coin"
)

var testHost = "127.0.0.1:8080"

func TestCoin(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	conn, err := grpc.NewClient(
		testHost,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithUnaryInterceptor(errors.ClientUnaryInterceptor),
	)
	if err != nil {
		log.Fatal().Msgf("gRPC connect: %v", err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Error().Msg("can't close connection")
		}
	}()
	client := pb.NewCoinServiceClient(conn)

	t.Run("GetActiveCoins", func(t *testing.T) {
		request := &pb.GetCoinsRequest{}
		resp, err := client.GetCoins(context.Background(), request)
		require.NoError(t, err)
		fmt.Println(resp)
	})
}
