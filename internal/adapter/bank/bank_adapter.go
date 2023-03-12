package bank

import (
	"context"
	"io"
	"log"

	dbank "github.com/timpamungkas/grpc-go-client/internal/application/domain/bank"
	"github.com/timpamungkas/grpc-go-client/internal/port"
	"github.com/timpamungkas/grpc-proto/protogen/go/bank"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BankAdapter struct {
	bankClient port.BankClientPort
}

func NewBankAdapter(conn *grpc.ClientConn) (*BankAdapter, error) {
	client := bank.NewBankServiceClient(conn)

	return &BankAdapter{
		bankClient: client,
	}, nil
}

func (a *BankAdapter) GetCurrentBalance(ctx context.Context, acct string) (*bank.CurrentBalanceResponse, error) {
	bankRequest := &bank.CurrentBalanceRequest{
		AccountNumber: acct,
	}

	bal, err := a.bankClient.GetCurrentBalance(ctx, bankRequest)

	if err != nil {
		log.Fatalln("Error on GetCurrentBalance : ", err)
	}

	return bal, nil
}

func (a *BankAdapter) FetchExchangeRates(ctx context.Context, fromCur string, toCur string) {
	bankRequest := &bank.ExchangeRateRequest{
		FromCurrency: fromCur,
		ToCurrency:   toCur,
	}

	exchangeRateStream, err := a.bankClient.FetchExchangeRates(ctx, bankRequest)

	if err != nil {
		log.Fatalln("Error on FetchExchangeRates : ", err)
	}

	for {
		rate, err := exchangeRateStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalln("Error on FetchExchangeRates : ", err)
		}

		log.Printf("Rate at %v from %v to %v is %v\n",
			rate.Timestamp, rate.FromCurrency, rate.ToCurrency, rate.Rate,
		)
	}
}

func (a *BankAdapter) SummarizeTransactions(ctx context.Context, acct string, tx []dbank.Transaction) {
	txStream, err := a.bankClient.SummarizeTransactions(ctx)

	if err != nil {
		log.Fatalln("Error on SummarizeTransactions : ", err)
	}

	for _, t := range tx {
		ttype := bank.TransactionType_TRANSACTION_TYPE_UNSPECIFIED

		if t.TransactionType == dbank.TransactionTypeIn {
			ttype = bank.TransactionType_TRANSACTION_TYPE_IN
		} else if t.TransactionType == dbank.TransactionTypeOut {
			ttype = bank.TransactionType_TRANSACTION_TYPE_OUT
		}

		bankRequest := &bank.Transaction{
			AccountNumber: acct,
			Type:          ttype,
			Amount:        t.Amount,
			Notes:         t.Notes,
		}

		txStream.Send(bankRequest)
	}

	summary, err := txStream.CloseAndRecv()

	if err != nil {
		log.Fatalln("Error on SummarizeTransactions : ", err)
	}

	log.Println(summary)
}

func (a *BankAdapter) TransferMultiple(ctx context.Context, trf []dbank.TransferTransaction) {
	trfStream, err := a.bankClient.TransferMultiple(ctx)

	if err != nil {
		// log.Fatalln("Error on TransferMultiple : ", err)
		st, _ := status.FromError(err)
		log.Println(st.Message())
	}

	trfChan := make(chan struct{})

	go func() {
		for _, tt := range trf {
			req := &bank.TransferRequest{
				FromAccountNumber: tt.FromAccountNumber,
				ToAccountNumber:   tt.ToAccountNumber,
				Currency:          tt.Currency,
				Amount:            tt.Amount,
			}

			trfStream.Send(req)
		}

		trfStream.CloseSend()
	}()

	go func() {
		for {
			res, err := trfStream.Recv()

			if err == io.EOF {
				break
			}

			// break transaction in case error
			// if err != nil {
			// 	st, _ := status.FromError(err)
			// 	log.Fatalln("[FATAL] Caught error on client : ", st.Message())
			// } else {
			// 	log.Printf("Transfer status : %v on %v\n", res.Status, res.Timestamp)
			// }

			if err != nil {
				st := status.Convert(err)

				if st.Code() == codes.FailedPrecondition {
					log.Fatalln("[FATAL] Failed precondition : ", st.Message())
				}

				for _, detail := range st.Details() {
					switch t := detail.(type) {
					case *errdetails.PreconditionFailure:
						for _, violation := range t.GetViolations() {
							log.Println("[VIOLATION]", violation)
						}
					case *errdetails.ErrorInfo:
						log.Printf("Error on : %v, with reason :%v\n", t.Domain, t.Reason)
						for k, v := range t.GetMetadata() {
							log.Printf("%v : %v\n", k, v)
						}
					}
				}

				break
			} else {
				log.Printf("Transfer status : %v on %v\n", res.Status, res.Timestamp)
			}

		}

		close(trfChan)
	}()

	<-trfChan
}
