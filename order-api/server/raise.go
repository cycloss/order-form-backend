package server

import (
	"fmt"
	"net/http"

	"github.com/Rhymond/go-money"
	"github.com/cycloss/aj-bell-test/share"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (hr *HandlerWrapper) RaiseHandler(c *gin.Context) {
	hr.transact(c, newRaiseProcessor)
}

type RaiseProcessor struct {
	tx       *gorm.DB
	context  *gin.Context
	username string
}

func newRaiseProcessor(c *gin.Context, tx *gorm.DB, username string) requestProcessor {
	return &RaiseProcessor{tx: tx, context: c, username: username}
}

func (sp *RaiseProcessor) process() (any, error) {
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

func (sp *RaiseProcessor) v1() (any, error) {

	bo, err := sp.unmarshalRaiseOrder()
	if err != nil {
		return nil, err
	}
	return bo.Process(sp.tx)
}

type RaiseOrder struct {
	AssetIsin      string `json:"asset-isin" binding:"required"`
	CurrencyAmount int64  `json:"currency-amount" binding:"required"`
	Username       string
}

func (sp *RaiseProcessor) unmarshalRaiseOrder() (*RaiseOrder, error) {
	var so RaiseOrder
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
func (ro *RaiseOrder) Process(tx *gorm.DB) (any, error) {

	investor, err := investorForUsername(ro.Username, tx)
	if err != nil {
		return nil, err

	}

	asset, err := getAssetBundle(ro.AssetIsin, tx)
	if err != nil {
		return nil, err
	}

	currency, err := currencyForCode(asset.PriceCode, tx)
	if err != nil {
		return nil, err
	}

	// how many assets can be sold for given currency amount
	unitsToSell := ro.CurrencyAmount / asset.Price
	requestedSaleValue := money.New(unitsToSell*asset.Price, asset.PriceCode)

	assetHolding, err := getAssetHolding(investor.Id, asset.Id, tx)
	if err != nil {
		return nil, err
	}

	if assetHolding.Units-int(unitsToSell) < 0 {
		return nil, share.NewApiErr(http.StatusUnprocessableEntity, "Asset holding insufficient", "")
	}

	err = updateCreateAssetHolding(investor.Id, asset.Id, -int(unitsToSell), tx)
	if err != nil {
		return nil, err
	}

	err = updateCreateCurrencyHolding(investor.Id, currency.Id, requestedSaleValue, tx)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("successfully raised %s by selling %d units of %s", requestedSaleValue.Display(), unitsToSell, ro.AssetIsin), nil
}
