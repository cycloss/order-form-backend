package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Rhymond/go-money"
	"github.com/cycloss/aj-bell-test/share"
	"gorm.io/gorm"
)

func investorForUsername(username string, tx *gorm.DB) (*share.Investor, error) {

	var investor share.Investor
	err := tx.Where("name = ?", username).First(&investor).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// record was not found
		return nil, share.NewApiErr(http.StatusUnprocessableEntity, fmt.Sprintf("Username: '%s' was not found in the database", username), "")
	} else if err != nil {
		return nil, err
	}
	return &investor, nil
}

func currencyForCode(code string, tx *gorm.DB) (*share.Currency, error) {

	var currency share.Currency
	err := tx.Where("code = ?", code).First(&currency).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// record was not found
		return nil, share.NewApiErr(http.StatusUnprocessableEntity, fmt.Sprintf("Currency: '%s' was not found in the database", code), "")
	} else if err != nil {
		return nil, err
	}
	return &currency, nil
}

func updateCreateAssetHolding(investorId, assetId, unitsToAdd int, tx *gorm.DB) error {
	queryRes := tx.Model(&share.AssetHolding{}).Where("investor_id = ? AND asset_id = ?", investorId, assetId).Update("units", gorm.Expr("units + ?", unitsToAdd))
	if err := queryRes.Error; err != nil {
		return err
	}

	if queryRes.RowsAffected < 1 {
		newHolding := share.NewAssetHolding(investorId, assetId, unitsToAdd)
		err := tx.Omit("id").Create(&newHolding).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func updateCreateCurrencyHolding(investorId, currencyId int, currencyToAdd *money.Money, tx *gorm.DB) error {
	queryRes := tx.Model(&share.CurrencyHolding{}).Where("investor_id = ? AND currency_id = ?", investorId, currencyId).Update("amount", gorm.Expr("amount + ?", currencyToAdd.Amount()))
	if err := queryRes.Error; err != nil {
		return err
	}
	if queryRes.RowsAffected < 1 {
		newHolding := share.NewCurrencyHolding(investorId, currencyId, int(currencyToAdd.Amount()))
		err := tx.Omit("id").Create(&newHolding).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func getAssetBundle(assetIsin string, tx *gorm.DB) (*assetBundle, error) {
	var assetBundle assetBundle
	err := tx.Model(&share.Asset{}).Select("assets.id id, price, code price_code").Joins("INNER JOIN currencies ON assets.currency_id = currencies.id").Where("isin = ?", assetIsin).First(&assetBundle).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// record was not found
		return nil, share.NewApiErr(http.StatusUnprocessableEntity, fmt.Sprintf("Asset with ISIN: '%s' was not found in the database", assetIsin), "")
	} else if err != nil {
		return nil, err
	}
	return &assetBundle, nil
}

type assetBundle struct {
	Id        int
	Price     int64
	PriceCode string
}

func getCurrencyHoldingBundle(investorId, currencyId int, tx *gorm.DB) (*currencyHoldingBundle, error) {
	var currencyHoldingBundle currencyHoldingBundle
	err := tx.Model(&share.CurrencyHolding{}).Select("currency_holdings.id id, amount, code").Joins("INNER JOIN currencies ON currency_holdings.currency_id = currencies.id").Where("investor_id = ? AND currency_id = ?", investorId, currencyId).First(&currencyHoldingBundle).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// record was not found
		return nil, share.NewApiErr(http.StatusUnprocessableEntity, "Currency holding insufficient", "")
	} else if err != nil {
		return nil, err
	}
	return &currencyHoldingBundle, nil
}

type currencyHoldingBundle struct {
	Id     int
	Amount int64
	Code   string
}

func getAssetHolding(investorId, assetId int, tx *gorm.DB) (*share.AssetHolding, error) {
	var assetHolding share.AssetHolding
	err := tx.Where("investor_id = ? AND asset_id = ?", investorId, assetId).First(&assetHolding).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// record was not found
		return nil, share.NewApiErr(http.StatusUnprocessableEntity, "Asset holding insufficient", "")
	} else if err != nil {
		return nil, err
	}
	return &assetHolding, err
}
