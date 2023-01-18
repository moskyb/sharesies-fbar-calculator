package portfolio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type PortfolioSliceValue struct {
	Percent int     `json:"percent"`
	Value   float64 `json:"value"`
}

type PortfolioDate time.Time

func (d *PortfolioDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	*d = PortfolioDate(t)
	return nil
}

func (d PortfolioDate) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	s := t.Format("2006-01-02")
	return []byte(s), nil
}

type PortfolioHistoryItem struct {
	Date                  PortfolioDate `json:"date"`
	PortfolioValue        float64       `json:"portfolio_value"`
	TotalReturn           float64       `json:"total_return"`
	UnrealisedTotalReturn float64       `json:"unrealised_total_return"`
	CostBasis             float64       `json:"cost_basis"`
}

type SharesiesPortfolio struct {
	Date *PortfolioDate `json:"date"`

	UUID                       string  `json:"uuid"`
	Currency                   string  `json:"currency"`
	PortfolioValue             float64 `json:"portfolio_value"`
	TotalReturn                float64 `json:"total_return"`
	SimpleReturn               float64 `json:"simple_return"`
	CostBasis                  float64 `json:"cost_basis"`
	CostBasisMax               float64 `json:"cost_basis_max"`
	PortfolioRiskType          string  `json:"portfolio_risk_type"`
	Dividends                  float64 `json:"dividends"`
	ManagedFundTransactionFees float64 `json:"managed_fund_transaction_fees"`
	TaxPaid                    float64 `json:"tax_paid"`
	TransactionFees            float64 `json:"transaction_fees"`
	ADRFees                    float64 `json:"adr_fees"`
	RealisedCapitalGain        float64 `json:"realised_capital_gain"`
	RealisedCurrencyGain       float64 `json:"realised_currency_gain"`
	RealisedCostBasis          float64 `json:"realised_cost_basis"`

	PortfolioRisk struct {
		Higher PortfolioSliceValue `json:"higher"`
		Medium PortfolioSliceValue `json:"medium"`
		Lower  PortfolioSliceValue `json:"lower"`
	} `json:"portfolio_risk"`

	PortfolioCountries struct {
		NZ PortfolioSliceValue `json:"nzl"`
		US PortfolioSliceValue `json:"usa"`
		AU PortfolioSliceValue `json:"aus"`
	} `json:"portfolio_countries"`

	PortfolioInstrumentTypes struct {
		ETF         PortfolioSliceValue `json:"etf"`
		Company     PortfolioSliceValue `json:"company"`
		ManagedFund PortfolioSliceValue `json:"mutual"`
	} `json:"portfolio_instrument_types"`

	PortfolioHistory []PortfolioHistoryItem `json:"portfolio_history"`

	UnrealisedDividends                  float64 `json:"unrealised_dividends"`
	UnrealisedManagedFundTransactionFees float64 `json:"unrealised_managed_fund_transaction_fees"`
	UnrealisedTotalReturn                float64 `json:"unrealised_total_return"`
	UnrealisedSimpleReturn               float64 `json:"unrealised_simple_return"`
	UnrealisedTaxPaid                    float64 `json:"unrealised_tax_paid"`
	UnrealisedTotalTransactionFees       float64 `json:"unrealised_total_transaction_fees"`
	UnrealisedADRFees                    float64 `json:"unrealised_adr_fees"`
	UnrealisedCapitalGain                float64 `json:"unrealised_capital_gain"`
	UnrealisedCurrencyGain               float64 `json:"unrealised_currency_gain"`
}

func Fetch(id, token string) (*SharesiesPortfolio, error) {
	url := fmt.Sprintf("https://portfolio.sharesies.nz/api/v1/portfolios/%s", id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "https://github.com/moskyb/sharesies-fbar-calculator (please send me an email at ben@mosk.nz if i'm causing trouble!)")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var portfolio SharesiesPortfolio
	err = json.NewDecoder(resp.Body).Decode(&portfolio)
	if err != nil {
		return nil, err
	}

	return &portfolio, nil
}
