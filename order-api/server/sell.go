package server

import (
	"fmt"
	"net/http"

	"github.com/Rhymond/go-money"
	"github.com/cycloss/aj-bell-test/share"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (hr *HandlerWrapper) SellHandler(c *gin.Context) {
	hr.transact(c, newSellProcessor)
}

type SellProcessor struct {
	tx       *gorm.DB
	context  *gin.Context
	username string
}

func newSellProcessor(c *gin.Context, tx *gorm.DB, username string) requestProcessor {
	return &SellProcessor{tx: tx, context: c, username: username}
}

func (sp *SellProcessor) process() (any, error) {
	version, err := share.GetAPIVersionFromHeader(sp.context.Request)
	if err != nil {
		return nil, err
	}

	switch version {
	case "v1":
		return sp.v1()
	default:
		return nil, share.NewApiErr(http.StatusNotFound, kUnrecognisedApiError, "")
	}
}

func (sp *SellProcessor) v1() (any, error) {

	bo, err := sp.unmarshalSellOrder()
	if err != nil {
		return nil, err
	}
	return bo.Process(sp.tx)
}

type SellOrder struct {
	AssetIsin string `json:"asset-isin" binding:"required"`
	Amount    int64  `json:"amount" binding:"required"`
	Username  string
}

func (sp *SellProcessor) unmarshalSellOrder() (*SellOrder, error) {
	var so SellOrder
	// do not need to use BindJson as we will handle the error ourselves (Bind immediately returns 400 if it fails)
	err := sp.context.ShouldBindJSON(&so)
	if err != nil {
		// error will contain the reason for binding failure
		return nil, share.NewApiErr(http.StatusBadRequest, err.Error(), "")
	}
	so.Username = sp.username
	return &so, err
}

// handles a sell order and returns a message on success
func (so *SellOrder) Process(tx *gorm.DB) (any, error) {

	investor, err := investorForUsername(so.Username, tx)
	if err != nil {
		return nil, err

	}

	asset, err := getAssetBundle(so.AssetIsin, tx)
	if err != nil {
		return nil, err
	}

	assetHolding, err := getAssetHolding(investor.Id, asset.Id, tx)
	if err != nil {
		return nil, err
	}

	if assetHolding.Units-int(so.Amount) < 0 {
		return "", share.NewApiErr(http.StatusUnprocessableEntity, "Asset holding insufficient", "")
	}

	err = updateCreateAssetHolding(investor.Id, asset.Id, -int(so.Amount), tx)
	if err != nil {
		return nil, err
	}

	currency, err := currencyForCode(asset.PriceCode, tx)
	if err != nil {
		return "", err
	}

	normalisedTotalvalue := asset.Price * so.Amount
	value := money.New(normalisedTotalvalue, asset.PriceCode)

	err = updateCreateCurrencyHolding(investor.Id, currency.Id, value, tx)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("successfully sold %d units of %s for %s", so.Amount, so.AssetIsin, value.Display()), nil
}
