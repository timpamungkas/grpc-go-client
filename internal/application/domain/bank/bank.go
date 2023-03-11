package bank

const (
	TransactionStatusIn  string = "IN"
	TransactionStatusOut string = "OUT"
)

type Transaction struct {
	Amount          float64
	TransactionType string
	Notes           string
}

type TransferTransaction struct {
	FromAccountNumber string
	ToAccountNumber   string
	Currency          string
	Amount            float64
}
