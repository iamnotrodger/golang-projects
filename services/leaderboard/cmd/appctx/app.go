package appctx

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/iamnotrodger/golang-projects/pkg/app"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/processes"
	"github.com/iamnotrodger/golang-projects/services/leaderboard/internal/score"
)

type AppContext struct {
	Engine       *gin.Engine
	scoreService *score.Service
}

func BuildAppProcesses(appCtx *AppContext) map[string]app.Runnable {
	return map[string]app.Runnable{
		"http": processes.NewHttpServer(appCtx.Engine, processes.HttpServerServices{
			ScoreService: appCtx.scoreService,
		}),
	}
}

func NewAppContext(ctx context.Context) *AppContext {
	appCtx := AppContext{}

	appCtx.Engine = gin.New()

	return &appCtx
}

func (a *AppContext) Shutdown(ctx context.Context) error {
	return nil
}
