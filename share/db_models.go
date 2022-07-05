package share

type Asset struct {
	Id   int `gorm:"primary_key, AUTO_INCREMENT"`
	Isin string
	Name string
}

type Price struct {
	Id         int `gorm:"primary_key, AUTO_INCREMENT"`
	AssetId    int
	Price      float64
	CurrencyId int
}

type Currency struct {
	Id   int `gorm:"primary_key, AUTO_INCREMENT"`
	Code string
}

type Investor struct {
	Id   int `gorm:"primary_key, AUTO_INCREMENT"`
	Name string
}

type AssetHolding struct {
	Id         int `gorm:"primary_key, AUTO_INCREMENT"`
	InvestorId int
	AssetId    int
	Units      int
}

type CurrencyHolding struct {
	Id         int `gorm:"primary_key, AUTO_INCREMENT"`
	InvestorId int
	CurrencyId int
	Amount     int
}
