package server

import (
	"fmt"
	"net/http"

	"github.com/Rhymond/go-money"
	"github.com/cycloss/aj-bell-test/share"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (hr *HandlerWrapper) InvestHandler(c *gin.Context) {
	hr.transact(c, newInvestProcessor)
}

type InvestProcessor struct {
	tx       *gorm.DB
	context  *gin.Context
	username string
}

func newInvestProcessor(c *gin.Context, tx *gorm.DB, username string) requestProcessor {
	return &InvestProcessor{tx: tx, context: c, username: username}
}

func (bp *InvestProcessor) process() (any, error) {
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

func (upr *InvestProcessor) v1() (any, error) {

	bo, err := upr.unmarshalInvestOrder()
	if err != nil {
		return nil, err
	}
	return bo.Process(upr.tx)
}

type InvestOrder struct {
	AssetIsin      string `json:"asset-isin" binding:"required"`
	CurrencyAmount int64  `json:"currency-amount" binding:"required"`
	Username       string
}

func (upp *InvestProcessor) unmarshalInvestOrder() (*InvestOrder, error) {
	var bo InvestOrder
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
func (io *InvestOrder) Process(tx *gorm.DB) (any, error) {

	investor, err := investorForUsername(io.Username, tx)
	if err != nil {
		return nil, err

	}

	asset, err := getAssetBundle(io.AssetIsin, tx)
	if err != nil {
		return nil, err
	}

	currency, err := currencyForCode(asset.PriceCode, tx)
	if err != nil {
		return nil, err
	}

	currencyHolding, err := getCurrencyHoldingBundle(investor.Id, currency.Id, tx)
	if err != nil {
		return nil, err
	}

	// how many assets can be bought for given currency amount
	purchasableUnits := io.CurrencyAmount / asset.Price
	requestedPurchaseValue := money.New(purchasableUnits*asset.Price, asset.PriceCode)

	availableCurrency := money.New(currencyHolding.Amount, currencyHolding.Code)

	net, err := availableCurrency.Subtract(requestedPurchaseValue)
	if err != nil {
		return nil, err
	}
	if net.Amount() < 0 {
		return nil, share.NewApiErr(http.StatusUnprocessableEntity, "Currency holding insufficient", "")
	}

	err = updateCreateCurrencyHolding(investor.Id, currency.Id, requestedPurchaseValue.Negative(), tx)
	if err != nil {
		return nil, err
	}

	err = updateCreateAssetHolding(investor.Id, asset.Id, int(purchasableUnits), tx)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("successfully invested in %d units of %s for %s", purchasableUnits, io.AssetIsin, requestedPurchaseValue.Display()), nil
}
