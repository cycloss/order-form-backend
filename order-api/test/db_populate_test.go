package test

import (
	"fmt"
	"testing"

	"github.com/Rhymond/go-money"
	"github.com/cycloss/aj-bell-test/share"
	"gorm.io/gorm"
)

func TestDbPopulate(t *testing.T) {
	db := share.MustConnectDb()
	t.Log("populating db...")
	err := db.Transaction(func(tx *gorm.DB) error {
		err := dbClear(tx)
		if err != nil {
			return err
		}
		return dbPopulate(tx)
	})
	if err != nil {
		t.Error(err)
	}
	t.Log("db populated")
}

func dbPopulate(tx *gorm.DB) error {
	assetMap, err := insertDummyAssets(tx)
	if err != nil {
		return err
	}
	currency, err := insertDummyCurrency(tx)
	if err != nil {
		return err
	}

	err = insertDummyPrices(tx, assetMap, currency.Id)
	if err != nil {
		return err
	}

	investor, err := insertDummyInvestor(tx)
	if err != nil {
		return err
	}

	err = insertDummyAssetHoldings(tx, assetMap, investor.Id)
	if err != nil {
		return err
	}

	err = insertDummyCurrencyHoldings(tx, investor.Id, currency.Id)
	if err != nil {
		return err
	}

	return nil
}

func insertDummyAssets(tx *gorm.DB) (map[string]*share.Asset, error) {
	assetMap := createDummyAssetMap()
	for _, v := range assetMap {
		err := tx.Omit("id").Create(v).Error
		if err != nil {
			return nil, err
		}
	}
	return assetMap, nil
}

func insertDummyCurrency(tx *gorm.DB) (*share.Currency, error) {
	currency := createDummyCurrency()

	err := tx.Omit("id").Create(currency).Error
	if err != nil {
		return nil, err
	}

	return currency, nil
}

func insertDummyInvestor(tx *gorm.DB) (*share.Investor, error) {
	investor := createDummyInvestor()

	err := tx.Omit("id").Create(investor).Error
	if err != nil {
		return nil, err
	}
	return investor, nil
}

func insertDummyPrices(tx *gorm.DB, assetMap map[string]*share.Asset, currencyId int) error {

	err := insertPrice(tx, currencyId, brIsin, brPrice, assetMap)
	if err != nil {
		return err
	}
	err = insertPrice(tx, currencyId, tnIsin, tnPrice, assetMap)
	if err != nil {
		return err
	}
	err = insertPrice(tx, currencyId, vgFtseIsin, vgFtsePrice, assetMap)
	if err != nil {
		return err
	}
	err = insertPrice(tx, currencyId, lgjIsin, lgjPrice, assetMap)
	if err != nil {
		return err
	}
	err = insertPrice(tx, currencyId, vgUsIsin, vgUsPrice, assetMap)
	if err != nil {
		return err
	}
	err = insertPrice(tx, currencyId, vgUkIsin, vgUkPrice, assetMap)
	if err != nil {
		return err
	}
	return nil
}

func insertPrice(tx *gorm.DB, currencyId int, isin string, assetPrice float64, assetMap map[string]*share.Asset) error {
	asset, ok := assetMap[isin]
	if !ok {
		return fmt.Errorf("asset with ISIN %s not found in map", isin)
	}
	price := share.NewPrice(asset.Id, currencyId, assetPrice)
	return tx.Omit("id").Create(price).Error
}

// 220,000 in pennies
const portfolioAmount = 22000000

func insertDummyAssetHoldings(tx *gorm.DB, assetMap map[string]*share.Asset, investorId int) error {

	portfolio := money.New(portfolioAmount, money.GBP)
	parties, err := portfolio.Allocate(20, 20, 10, 5, 15, 30)
	if err != nil {
		return err
	}

	err = insertAssetHolding(tx, investorId, brIsin, brPrice, parties[0].Amount(), assetMap)
	if err != nil {
		return err
	}
	err = insertAssetHolding(tx, investorId, tnIsin, tnPrice, parties[1].Amount(), assetMap)
	if err != nil {
		return err
	}
	err = insertAssetHolding(tx, investorId, vgFtseIsin, vgFtsePrice, parties[2].Amount(), assetMap)
	if err != nil {
		return err
	}
	err = insertAssetHolding(tx, investorId, lgjIsin, lgjPrice, parties[3].Amount(), assetMap)
	if err != nil {
		return err
	}
	err = insertAssetHolding(tx, investorId, vgUsIsin, vgUsPrice, parties[4].Amount(), assetMap)
	if err != nil {
		return err
	}
	err = insertAssetHolding(tx, investorId, vgUkIsin, vgUkPrice, parties[5].Amount(), assetMap)
	if err != nil {
		return err
	}
	return nil
}

func insertAssetHolding(tx *gorm.DB, investorId int, isin string, assetPrice float64, amount int64, assetMap map[string]*share.Asset) error {
	asset, ok := assetMap[isin]
	if !ok {
		return fmt.Errorf("asset with ISIN %s not found in map", isin)
	}
	// multiply price by 100 as amount is also 100 more than it is nominally
	units := float64(amount) / (assetPrice * 100)
	// there will be some currency left over when converted to an int
	assetHolding := share.NewAssetHolding(investorId, asset.Id, int(units))
	return tx.Omit("id").Create(assetHolding).Error
}

const oneHundredThousand = 10000000

func insertDummyCurrencyHoldings(tx *gorm.DB, investorId, currencyId int) error {
	currencyHolding := share.NewCurrencyHolding(investorId, currencyId, oneHundredThousand)

	err := tx.Omit("id").Create(currencyHolding).Error
	if err != nil {
		return err
	}
	return nil
}
