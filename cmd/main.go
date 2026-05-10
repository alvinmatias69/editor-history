package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alvinmatias69/editor-history/internal/config"
	"github.com/alvinmatias69/editor-history/internal/controller"
	"github.com/alvinmatias69/editor-history/internal/handler"
	"github.com/alvinmatias69/editor-history/internal/repository"
	"github.com/alvinmatias69/editor-history/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	baseCfg := config.GetBaseCfg()
	dbCfg := config.GetDBCfg()
	redisCfg, err := config.GetRedisCfg()
	if err != nil {
		log.Fatalf("Error parsing redis configuration: %v\n", err)
	}

	pgpool, err := pgxpool.New(context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DbName))
	if err != nil {
		log.Fatalf("Error initiating db connection: %v\n", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		PoolSize: 5,
	})

	documentRepository := repository.NewDocumentRepository(pgpool)
	documentOperationRepository := repository.NewDocumentOperationRepository(pgpool)
	sessionRepository := repository.NewSessionRepository(redisClient, redisCfg)

	documentController := controller.NewDocumentController(documentRepository, documentOperationRepository, sessionRepository)

	templateHandler := handler.NewTemplateHandler()
	documentHandler := handler.NewDocumentHandler(*baseCfg, documentController)

	server := server.New(documentHandler, templateHandler)

	server.Start(baseCfg.Port)
}
