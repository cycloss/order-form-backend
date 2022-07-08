package share

type Asset struct {
	Id         int `gorm:"primary_key, AUTO_INCREMENT"`
	Isin       string
	Name       string
	Price      int
	CurrencyId int
}

func NewAsset(isin, name string, price int, currencyId int) *Asset {
	return &Asset{Isin: isin, Name: name, Price: price, CurrencyId: currencyId}
}

type Currency struct {
	Id   int `gorm:"primary_key, AUTO_INCREMENT"`
	Code string
}

func NewCurrency(code string) *Currency {
	return &Currency{Code: code}
}

type Investor struct {
	Id   int `gorm:"primary_key, AUTO_INCREMENT"`
	Name string
	Pass string
}

func NewInvestor(name, pass string) *Investor {
	return &Investor{Name: name, Pass: pass}
}

type AssetHolding struct {
	Id         int `gorm:"primary_key, AUTO_INCREMENT"`
	InvestorId int
	AssetId    int
	Units      int
}

func NewAssetHolding(investorId, assetId, units int) *AssetHolding {
	return &AssetHolding{InvestorId: investorId, AssetId: assetId, Units: units}
}

type CurrencyHolding struct {
	Id         int `gorm:"primary_key, AUTO_INCREMENT"`
	InvestorId int
	CurrencyId int
	Amount     int
}

func NewCurrencyHolding(investorId, currencyId, amount int) *CurrencyHolding {
	return &CurrencyHolding{InvestorId: investorId, CurrencyId: currencyId, Amount: amount}
}
