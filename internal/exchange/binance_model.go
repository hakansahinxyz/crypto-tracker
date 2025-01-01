package exchange

type AccountInfo struct {
	Balances []struct {
		Asset  string `json:"asset"`
		Free   string `json:"free"`
		Locked string `json:"locked"`
	} `json:"balances"`
}

type MarginBalance struct {
	Asset    string `json:"asset"`
	Borrowed string `json:"borrowed"`
	Free     string `json:"free"`
	Interest string `json:"interest"`
	Locked   string `json:"locked"`
	NetAsset string `json:"netAsset"`
}

type MarginAccountResponse struct {
	MarginLevel                string          `json:"marginLevel"`
	CollateralMarginLevel      string          `json:"collateralMarginLevel"`
	TotalCollateralValueInUSDT string          `json:"totalCollateralValueInUSDT"`
	Assets                     []MarginBalance `json:"userAssets"`
}

type FutureBalance struct {
	Asset         string `json:"asset"`
	Balance       string `json:"balance"`
	UnrealizedPNL string `json:"crossUnPnl"`
}

type FutureBalanceResponse []FutureBalance

type Price struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}
