package test

import (
	"github.com/Rhymond/go-money"
	"github.com/cycloss/aj-bell-test/share"
)

// must be in sub package to be imported into other programs
// package main can never be imported

const (
	brIsin     = "IE00B52L4369"
	tnIsin     = "GB00BQ1YHQ70"
	vgFtseIsin = "GB00B3X7QG63"
	lgjIsin    = "GB00BG0QP828"
	vgUsIsin   = "GB00BPN5P238"
	vgUkIsin   = "IE00B1S74Q32"
)

const (
	brPrice     = 150
	tnPrice     = 254
	vgFtsePrice = 065
	lgjPrice    = 129
	vgUsPrice   = 203
	vgUkPrice   = 163
)

// use of pointer allows struct to be updated by gorm and then passed back
func createDummyAssetMap(currencyId int) map[string]*share.Asset {
	return map[string]*share.Asset{
		brIsin:     share.NewAsset(brIsin, "BlackRock Institutional Cash Series Sterling Liquidity Agency Inc", brPrice, currencyId),
		tnIsin:     share.NewAsset(tnIsin, "Threadneedle UK Property Authorised Investment Net GBP 1 Acc", tnPrice, currencyId),
		vgFtseIsin: share.NewAsset(vgFtseIsin, "Vanguard FTSE U.K. All Share Index Unit Trust Accumulation", vgFtsePrice, currencyId),
		lgjIsin:    share.NewAsset(lgjIsin, "Legal & General Japan Index Trust C Class Accumulation", lgjPrice, currencyId),
		vgUsIsin:   share.NewAsset(vgUsIsin, "Vanguard US Equity Index Institutional Plus GBP Accumulation", vgUsPrice, currencyId),
		vgUkIsin:   share.NewAsset(vgUkIsin, "Vanguard U.K. Investment Grade Bond Index Fund GBP Accumulation", vgUkPrice, currencyId),
	}
}

func createDummyCurrency() *share.Currency {
	return share.NewCurrency(money.GBP)
}

const (
	testInvestorName = "AJ Bell"
	testInvestorPass = "pass"
)

func createDummyInvestor() *share.Investor {
	return share.NewInvestor(testInvestorName, testInvestorPass)
}
