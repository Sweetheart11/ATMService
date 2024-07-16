package balance

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Msg    string `json:"message,omitempty"`
}

type BalanceChecker interface {
	GetAccountBalance(id int) (float64, error)
}

func New(log *slog.Logger, balanceChecker BalanceChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.acc.balance"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userID, err := strconv.Atoi(chi.URLParam(r, "id"))
		fmt.Println(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("failed to get user id from url", slog.StringValue(err.Error()))

			responseError(w, r, "failed to get user id from url")

			return
		}

		ch := make(chan float64)

		go func() {
			balance, err := balanceChecker.GetAccountBalance(userID)
			if err != nil {
				log.Error("failed to get balance", slog.StringValue(err.Error()))
				responseError(w, r, "failed to get balance")
				return
			}

			ch <- balance
		}()

		res := <-ch

		log.Info("current balance", slog.String("deposit", fmt.Sprintf("%f", res)))

		responseOK(w, r, fmt.Sprintf("current balance of an account %d: %f", userID, res))
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, msg string) {
	render.JSON(w, r, Response{
		Status: "OK",
		Msg:    msg,
	})
}

func responseError(w http.ResponseWriter, r *http.Request, err string) {
	render.JSON(w, r, Response{
		Status: "Error",
		Error:  err,
	})
}
