package main

import (
	"net/http"

	"github.com/cycloss/aj-bell-test/share"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (hr *handlerWrapper) buyHandler(c *gin.Context) {
	hr.transact(c, newBuyProcessor)
}

type buyProcessor struct {
	tx        *gorm.DB
	context   *gin.Context
	uInvestId int
}

func newBuyProcessor(c *gin.Context, tx *gorm.DB, uInvestId int) requestProcessor {
	return &buyProcessor{tx: tx, context: c, uInvestId: uInvestId}
}

func (bp *buyProcessor) process() (any, error) {
	version, err := share.GetAPIVersionFromHeader(bp.context.Request)
	if err != nil {
		return nil, err
	}

	switch version {
	case "v1":
		return bp.V1()
	default:
		return nil, share.NewApiErr(http.StatusNotFound, kUnrecognisedApiError, "")
	}
}

func (upr *buyProcessor) V1() (any, error) {
	return nil, nil
}
