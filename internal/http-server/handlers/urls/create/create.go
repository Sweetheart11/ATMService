package create

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Sweetheart11/ATMService/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Request struct {
	Username string `json:"username"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Msg    string `json:"message,omitempty"`
}

type AccCreator interface {
	CreateAccount(username string) (int, error)
}

func New(log *slog.Logger, accCreator AccCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.acc.create.New"

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

		if req.Username == "" {
			log.Error("username is empty")

			responseError(w, r, "username is empty")

			return

		}

		id, err := accCreator.CreateAccount(req.Username)
		if errors.Is(err, storage.ErrUserExists) {
			log.Info("user already exists", slog.String("url", req.Username))

			responseError(w, r, "account with that username already exists")

			return
		}
		if err != nil {
			log.Error("failed to create new account", slog.StringValue(err.Error()))

			responseError(w, r, "failed to to create new account")

			return
		}

		log.Info("new account created", slog.Int("id", id))

		responseOK(w, r, fmt.Sprintf("account %s with id %d created", req.Username, id))
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
