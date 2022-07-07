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

// use of pointer allows struct to be updated by gorm and then passed back
func createDummyAssetMap() map[string]*share.Asset {
	return map[string]*share.Asset{
		brIsin:     share.NewAsset(brIsin, "BlackRock Institutional Cash Series Sterling Liquidity Agency Inc"),
		tnIsin:     share.NewAsset(tnIsin, "Threadneedle UK Property Authorised Investment Net GBP 1 Acc"),
		vgFtseIsin: share.NewAsset(vgFtseIsin, "Vanguard FTSE U.K. All Share Index Unit Trust Accumulation"),
		lgjIsin:    share.NewAsset(lgjIsin, "Legal & General Japan Index Trust C Class Accumulation"),
		vgUsIsin:   share.NewAsset(vgUsIsin, "Vanguard US Equity Index Institutional Plus GBP Accumulation"),
		vgUkIsin:   share.NewAsset(vgUkIsin, "Vanguard U.K. Investment Grade Bond Index Fund GBP Accumulation"),
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

const (
	brPrice     float64 = 1.82
	tnPrice     float64 = 2.54
	vgFtsePrice float64 = 0.65
	lgjPrice    float64 = 1.29
	vgUsPrice   float64 = 2.03
	vgUkPrice   float64 = 1.63
)
