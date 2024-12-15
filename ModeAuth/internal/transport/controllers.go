package transport

import (
	"ModeAuth/internal/service"
	errors2 "ModeAuth/internal/service/errors"
	"ModeAuth/internal/shared/dto"
	"ModeAuth/pkg/logging"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type Controller struct {
	aService service.IAuth
	sService service.IStates
}

func NewController(aService service.IAuth, sService service.IStates) *Controller {
	return &Controller{aService, sService}
}

func (c *Controller) CheckUser() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf(logging.INFO+"Received a new request from %s", req.RemoteAddr)
		var user dto.RequestData
		if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
			log.Printf(logging.ERROR+"[Transport] json.Decoder.Decode failed: %v", err)

			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		ctx := req.Context()

		userTransport, err := c.aService.CheckUser(ctx, user.ID, user.UserName)
		if err != nil {
			if errors.Is(err, errors2.UserNotExist) {
				http.Error(w, fmt.Sprintf("%s", err), http.StatusNotFound)
				return
			}
			if errors.Is(err, errors2.UserIsBot) {
				http.Error(w, fmt.Sprintf("%s", err), http.StatusForbidden)
				return
			}

			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(userTransport)
		if err != nil {
			log.Printf(logging.ERROR+"[Transport] json.Marshal failed: %v", err)

			http.Error(w, "ошибка при сериализация в JSON", http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(data); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}
}

func (c *Controller) CheckUserIsBlocked() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf(logging.INFO+"Received a new request from %s", req.RemoteAddr)
		var user dto.RequestData
		if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
			log.Printf(logging.ERROR+"[Transport] json.Decoder.Decode failed: %v", err)

			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		ctx := req.Context()

		durationTransport, err := c.aService.CheckUserIsBlocked(ctx, user.ID)
		if err != nil {
			if errors.Is(err, errors2.UserNotExist) {
				http.Error(w, fmt.Sprintf("%s", err), http.StatusNotFound)
				return
			}
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		data, err := json.Marshal(durationTransport)
		if err != nil {
			log.Printf(logging.ERROR+"[Transport] json.Marshal failed: %v", err)

			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(data); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}
}

func (c *Controller) GoWorkState() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf(logging.INFO+"Received a new request from %s", req.RemoteAddr)
		var user dto.RequestData
		if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
			log.Printf(logging.ERROR+"[Transport] json.Decoder.Decode failed: %v", err)

			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		ctx := req.Context()

		status, err := c.sService.WorkState(ctx, user.ID)
		if err != nil {
			if errors.Is(err, errors2.AlreadyStateExist) {
				http.Error(w, fmt.Sprintf("%s", err), http.StatusConflict)
				return
			}
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		if status {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (c *Controller) SendingState() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf(logging.INFO+"Received a new request from %s", req.RemoteAddr)
		var user dto.RequestData
		if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
			log.Printf(logging.ERROR+"[Transport] json.Decoder.Decode failed: %v", err)

			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		ctx := req.Context()

		status, err := c.sService.SendingState(ctx, user.ID)
		if err != nil {
			if errors.Is(err, errors2.AlreadyStateExist) {
				http.Error(w, fmt.Sprintf("%s", err), http.StatusConflict)
				return
			}
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		if status {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (c *Controller) CheckingState() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf(logging.INFO+"Received a new request from %s", req.RemoteAddr)
		var user dto.RequestData
		if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
			log.Printf(logging.ERROR+"[Transport] json.Decoder.Decode failed: %v", err)

			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		ctx := req.Context()

		status, err := c.sService.CheckingState(ctx, user.ID)
		if err != nil {
			if errors.Is(err, errors2.AlreadyStateExist) {
				http.Error(w, fmt.Sprintf("%s", err), http.StatusConflict)
				return
			}
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return
		}

		if status {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (c *Controller) GetReport() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf(logging.INFO+"Received a new request from %s", req.RemoteAddr)
		var user dto.RequestData
		if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
			log.Printf(logging.ERROR+"[Transport] json.Decoder.Decode failed: %v", err)

			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		ctx := req.Context()

		report, err := c.sService.TakeForCheckingState(ctx, user.ID)
		if err != nil {
			if errors.Is(err, errors2.UserNotExist) {
				http.Error(w, fmt.Sprintf("%s", err), http.StatusNotFound)
				return
			}
			if errors.Is(err, errors2.CheckingStateIsNotSet) {
				http.Error(w, fmt.Sprintf("%s", err), http.StatusForbidden)
				return
			}
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		data, err := json.Marshal(report)
		if err != nil {
			log.Printf(logging.ERROR+"[Transport] json.Marshal failed: %v", err)

			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}

		if _, err := w.Write(data); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}
}
