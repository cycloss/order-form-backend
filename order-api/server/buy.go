package server

import (
	"fmt"
	"net/http"

	"github.com/Rhymond/go-money"
	"github.com/cycloss/aj-bell-test/share"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (hr *HandlerWrapper) BuyHandler(c *gin.Context) {
	hr.transact(c, newBuyProcessor)
}

type BuyProcessor struct {
	tx       *gorm.DB
	context  *gin.Context
	username string
}

func newBuyProcessor(c *gin.Context, tx *gorm.DB, username string) requestProcessor {
	return &BuyProcessor{tx: tx, context: c, username: username}
}

func (bp *BuyProcessor) process() (any, error) {
	version, err := share.GetAPIVersionFromHeader(bp.context.Request)
	if err != nil {
		return nil, err
	}

	switch version {
	case "v1":
		return bp.v1()
	default:
		return nil, share.NewApiErr(http.StatusNotFound, kUnrecognisedApiError, "")
	}
}

func (upr *BuyProcessor) v1() (any, error) {

	bo, err := upr.unmarshalBuyOrder()
	if err != nil {
		return nil, err
	}
	return bo.Process(upr.tx)
}

type BuyOrder struct {
	AssetIsin string `json:"asset-isin" binding:"required"`
	Amount    int64  `json:"amount" binding:"required"`
	Username  string
}

func (upp *BuyProcessor) unmarshalBuyOrder() (*BuyOrder, error) {
	var bo BuyOrder
	// do not need to use BindJson as we will handle the error ourselves (Bind immediately returns 400 if it fails)
	err := upp.context.ShouldBindJSON(&bo)
	if err != nil {
		// error will contain the reason for binding failure
		return nil, share.NewApiErr(http.StatusBadRequest, err.Error(), "")
	}
	bo.Username = upp.username
	return &bo, err
}

// handles a buy order and returns a message on success
func (bo *BuyOrder) Process(tx *gorm.DB) (any, error) {

	investor, err := investorForUsername(bo.Username, tx)
	if err != nil {
		return nil, err

	}

	asset, err := getAssetBundle(bo.AssetIsin, tx)
	if err != nil {
		return nil, err
	}

	currency, err := currencyForCode(asset.PriceCode, tx)
	if err != nil {
		return "", err
	}

	currencyHolding, err := getCurrencyHoldingBundle(investor.Id, currency.Id, tx)
	if err != nil {
		return nil, err
	}

	normalisedTotalCost := asset.Price * bo.Amount
	value := money.New(normalisedTotalCost, asset.PriceCode)
	available := money.New(currencyHolding.Amount, currencyHolding.Code)

	res, err := available.Subtract(value)
	if err != nil {
		return "", err
	}

	if res.Amount() < 0 {
		return "", share.NewApiErr(http.StatusUnprocessableEntity, "Currency holdings insufficient", "")
	}

	err = updateCreateCurrencyHolding(investor.Id, currency.Id, value.Negative(), tx)
	if err != nil {
		return nil, err
	}

	err = updateCreateAssetHolding(investor.Id, asset.Id, int(bo.Amount), tx)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("successfully bought %d units of %s for %s", bo.Amount, bo.AssetIsin, value.Display()), nil
}
