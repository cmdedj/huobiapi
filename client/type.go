package client

type ParamData = map[string]string

type Balance struct {
	Currency string `json:"currency"`
	Type     string `json:"type"`
	Balance  string `json:"balance"`
}

type DepositAndWithdraw struct {
	Id         int64   `json:"id"`
	Type       string  `json:"type"`
	Currency   string  `json:"currency"`
	Chain      string  `json:"chain"`
	TxHash     string  `json:"tx-hash"`
	Amount     float64 `json:"amount"`
	Address    string  `json:"address"`
	AddressTag string  `json:"address-tag"`
	Fee        int64   `json:"fee"`
	State      string  `json:"state"`
	CreatedAt  int64   `json:"created-at"`
	UpdatedAt  int64   `json:"updated-at"`
}

type Order struct {
	Id              int64  `json:"id"`
	Symbol          string `json:"symbol"`
	AccountId       int64  `json:"account-id"`
	Amount          string `json:"amount"`
	Price           string `json:"price"`
	CreatedAt       int64  `json:"created-at"`
	Type            string `json:"type"`
	FieldAmount     string `json:"field-amount"`
	FieldCashAmount string `json:"field-cash-amount"`
	FieldFees       string `json:"field-fees"`
	FinishedAt      int64  `json:"finished-at"`
	UserId          int64  `json:"user-id"`
	Source          string `json:"source"`
	State           string `json:"state"`
	CanceledAt      int64  `json:"canceled-at"`
}
