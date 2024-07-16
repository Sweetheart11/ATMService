package deposit

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Request struct {
	Amount string `json:"amount"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Msg    string `json:"message,omitempty"`
}

type DepositMaker interface {
	DepositToAccount(id int, amount float64) (float64, error)
}

func New(log *slog.Logger, depositMaker DepositMaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.acc.deposit"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			responseError(w, r, "empty request")

			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.StringValue(err.Error()))

			responseError(w, r, "failed to decode request")

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		amount, err := strconv.ParseFloat(req.Amount, 64)
		if err != nil {
			log.Error("failed to parse amount", slog.StringValue(err.Error()))

			responseError(w, r, "failed to parse amount")

			return
		}

		userID, err := strconv.Atoi(chi.URLParam(r, "id"))
		fmt.Println(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("failed to get user id from url", slog.StringValue(err.Error()))

			responseError(w, r, "failed to get user id from url")

			return
		}

		ch := make(chan float64)

		go func() {
			balance, err := depositMaker.DepositToAccount(userID, amount)
			if err != nil {
				log.Error("failed to deposit to account", slog.StringValue(err.Error()))
				responseError(w, r, "failed to deposit to account")
				return
			}

			ch <- balance
		}()

		res := <-ch

		log.Info("new deposit", slog.String("deposit", fmt.Sprintf("%f", amount)), slog.String("new balance", fmt.Sprintf("%f", res)))

		responseOK(w, r, fmt.Sprintf("made deposit %f to account %d. new balance: %f", amount, userID, res))
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
