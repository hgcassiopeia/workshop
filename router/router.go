package router

import (
	"database/sql"
	"net/http"

	cloud_pockets "github.com/kkgo-software-engineering/workshop/cloud_pocket"

	"github.com/kkgo-software-engineering/workshop/account"
	"github.com/kkgo-software-engineering/workshop/config"
	"github.com/kkgo-software-engineering/workshop/featflag"
	"github.com/kkgo-software-engineering/workshop/healthchk"
	mw "github.com/kkgo-software-engineering/workshop/middleware"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func RegRoute(cfg config.Config, logger *zap.Logger, db *sql.DB) *echo.Echo {
	e := echo.New()
	e.Use(mlog.Middleware(logger))
	e.Use(middleware.BasicAuth(mw.Authenicate()))

	hHealthChk := healthchk.New(db)
	e.GET("/healthz", hHealthChk.Check)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	hAccount := account.New(cfg.FeatureFlag, db)
	e.POST("/accounts", hAccount.Create)

	hFeatFlag := featflag.New(cfg)
	e.GET("/features", hFeatFlag.List)

	cloudPockets := cloud_pockets.New(cfg.FeatureFlag, db)
	e.POST("/cloud-pockets", cloudPockets.CreateCloudPockets)
	e.POST("/cloud-pockets/transfer", cloudPockets.Transfer)
	e.GET("/cloud-pockets", cloudPockets.GetAll)
	e.GET("/cloud-pockets/:id", cloudPockets.GetById)

	return e
}
