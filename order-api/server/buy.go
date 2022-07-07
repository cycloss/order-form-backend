package server

import (
	"net/http"

	"github.com/cycloss/aj-bell-test/share"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (hr *HandlerWrapper) BuyHandler(c *gin.Context) {
	hr.transact(c, newBuyProcessor)
}

type BuyProcessor struct {
	tx        *gorm.DB
	context   *gin.Context
	uInvestId int
}

func newBuyProcessor(c *gin.Context, tx *gorm.DB, uInvestId int) requestProcessor {
	return &BuyProcessor{tx: tx, context: c, uInvestId: uInvestId}
}

func (bp *BuyProcessor) process() (any, error) {
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

func (upr *BuyProcessor) V1() (any, error) {
	return "success", nil
}
