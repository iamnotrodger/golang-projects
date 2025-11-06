package appctx

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/iamnotrodger/golang-projects/pkg/app"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/config"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/leaderboard"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/processes"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/score"
	"github.com/redis/go-redis/v9"
)

type AppContext struct {
	engine       *gin.Engine
	rdb          *redis.Client
	hub          *leaderboard.Hub
	scoreService *score.Service
}

func BuildAppProcesses(appCtx *AppContext) map[string]app.Runnable {
	httpServices := processes.HttpServerServices{
		ScoreService: appCtx.scoreService,
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

	return &appCtx
}

func (a *AppContext) Shutdown(ctx context.Context) error {
	return nil
}
