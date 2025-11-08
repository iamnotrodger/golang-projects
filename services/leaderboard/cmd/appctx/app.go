package appctx

import (
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/iamnotrodger/golang-projects/pkg/app"
	"github.com/iamnotrodger/golang-projects/pkg/health"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/config"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/healthcheck"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/leaderboard"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/processes"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/score"
	"github.com/redis/go-redis/v9"
)

type AppContext struct {
	engine        *gin.Engine
	rdb           *redis.Client
	hub           *leaderboard.Hub
	scoreService  *score.Service
	healthService *health.Service
}

func BuildAppProcesses(appCtx *AppContext) map[string]app.Runnable {
	httpServices := processes.HttpServerServices{
		ScoreService:  appCtx.scoreService,
		HealthService: appCtx.healthService,
	}

	return map[string]app.Runnable{
		"http":                   processes.NewHttpServer(appCtx.engine, appCtx.hub, httpServices),
		"leaderboard-subscriber": processes.NewLeaderSubscriber(appCtx.rdb, appCtx.hub.Broadcast),
	}
}

func NewAppContext(ctx context.Context) *AppContext {
	appCtx := AppContext{}

	appCtx.engine = gin.New()
	appCtx.hub = leaderboard.NewHub()
	appCtx.rdb = redis.NewClient(&redis.Options{
		Addr:     config.Global.RedisAddr,
		Password: config.Global.RedisPassword,
		DB:       config.Global.RedisDb,
	})

	appCtx.scoreService = score.NewService(appCtx.rdb)
	appCtx.healthService = health.NewService(map[string]health.HealthCheck{
		"redis": healthcheck.NewRedisCheck(appCtx.rdb),
	})

	return &appCtx
}

func (a *AppContext) Shutdown(ctx context.Context) error {
	if err := a.rdb.Close(); err != nil {
		slog.Error("error closing redis client", "error", err.Error())
		return err
	}
	return nil
}
