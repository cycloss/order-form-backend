package test

import (
	"testing"

	"github.com/cycloss/aj-bell-test/share"
	"gorm.io/gorm"
)

func TestDbClear(t *testing.T) {
	db := share.MustConnectDb()
	t.Log("clearing db...")
	err := db.Transaction(func(tx *gorm.DB) error {
		return dbClear(tx)
	})
	if err != nil {
		t.Error(err)
	}
	t.Log("db cleared")
}

func dbClear(tx *gorm.DB) error {

	err := tx.Delete(&share.CurrencyHolding{}, "1 = 1").Error
	if err != nil {
		return err
	}
	err = tx.Delete(&share.AssetHolding{}, "1 = 1").Error
	if err != nil {
		return err
	}
	err = tx.Delete(&share.Investor{}, "1 = 1").Error
	if err != nil {
		return err
	}
	err = tx.Delete(&share.Price{}, "1 = 1").Error
	if err != nil {
		return err
	}
	err = tx.Delete(&share.Currency{}, "1 = 1").Error
	if err != nil {
		return err
	}
	err = tx.Delete(&share.Asset{}, "1 = 1").Error
	if err != nil {
		return err
	}
	return nil
}
