package main

import (
	"context"
	"fmt"
	api "github.com/Azreil-OFD/Avito-test/internal/generate"
	"github.com/Azreil-OFD/Avito-test/internal/handler"
	AuthHandler "github.com/Azreil-OFD/Avito-test/internal/handler/auth"
	BuyItemHandler "github.com/Azreil-OFD/Avito-test/internal/handler/buyItem"
	InfoHandler "github.com/Azreil-OFD/Avito-test/internal/handler/info"
	SendCoinHandler "github.com/Azreil-OFD/Avito-test/internal/handler/sendCoin"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/config"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/database/postgres"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/database/redis"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/logger"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/middleware"
	"github.com/Azreil-OFD/Avito-test/internal/infrastructure/utils"
	AuthRepository "github.com/Azreil-OFD/Avito-test/internal/repository/auth"
	BuyItemRepository "github.com/Azreil-OFD/Avito-test/internal/repository/buyItem"
	InfoRepository "github.com/Azreil-OFD/Avito-test/internal/repository/info"
	SendCoinRepository "github.com/Azreil-OFD/Avito-test/internal/repository/sendCoin"
	AuthService "github.com/Azreil-OFD/Avito-test/internal/usecase/auth"
	BuyItemService "github.com/Azreil-OFD/Avito-test/internal/usecase/buyItem"
	InfoService "github.com/Azreil-OFD/Avito-test/internal/usecase/info"
	SendCoinService "github.com/Azreil-OFD/Avito-test/internal/usecase/sendCoin"
	trmpgxs "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	cfg := config.MustNewConfigWithEnv()
	log := logger.NewLogger(cfg.DevMode(), slog.LevelDebug)
	jwtUtils := utils.NewToken(cfg.JwtSecret())

	pg, err := postgres.New(cfg.PgUrl())
	if err != nil {
		log.WarnContext(ctx, err.Error())
	}
	rdb, err := redis.New(cfg.RdbAddr(), cfg.RdbPassword(), cfg.RdbDb())
	if err != nil {
		log.WarnContext(ctx, err.Error())
	}
	defer rdb.Close()
	defer pg.Close()
	trManager := manager.Must(trmpgxs.NewDefaultFactory(pg.Pool))
	authRepository := AuthRepository.New(pg.Pool)
	infoRepository := InfoRepository.New(pg.Pool)
	buyItemRepository := BuyItemRepository.New(pg.Pool)
	sendCoinRepository := SendCoinRepository.New(pg.Pool)

	authService := AuthService.New(authRepository, jwtUtils, rdb)
	infoService := InfoService.New(infoRepository, rdb)
	buyItemService := BuyItemService.New(buyItemRepository, trManager, rdb)
	sendCoinService := SendCoinService.New(sendCoinRepository, trManager, rdb)

	handlers := &handler.Handlers{
		AuthHandler:     AuthHandler.New(*log, authService),
		BuyItemHandler:  BuyItemHandler.New(*log, buyItemService),
		SendCoinHandler: SendCoinHandler.New(*log, sendCoinService),
		InfoHandler:     InfoHandler.New(*log, infoService),
	}
	middleware := middleware.New(jwtUtils)

	server, err := api.NewServer(handlers, middleware)
	if err != nil {
		log.WarnContext(ctx, err.Error())
		os.Exit(1)
	}

	log.WarnContext(ctx, "server start, listen ::"+cfg.HttpPort())
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.HttpPort()), server); err != nil {
		log.WarnContext(ctx, err.Error())
		os.Exit(1)
	}

}
