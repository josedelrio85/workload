package workload

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Handler is an struct to handle models, database and http handler
type Handler struct {
	DbHandler     Storer
	ClientHandler http.Handler
	ActiveModels  []Modelable
}

// TODO | receive single entities instead arrays, maybe add
// TODO | add an array of uint to handle get multiple methods

// Input is a struct to handle the input data
type Input struct {
	InputType int64     `json:"type"`
	Project   []Project `json:"project"`
	Work      []Work    `json:"work"`
	Status    []Status  `json:"status"`
	Person    []Person  `json:"person"`
}

// Get searchs for active input and call its get method
func (h *Handler) Get() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := Input{}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil && err != io.EOF {
			resp := Response{
				Writer:  w,
				Message: "Error decoding project input",
				Code:    http.StatusInternalServerError,
				Error:   err,
			}
			resp.Ko()
			return
		}

		for _, am := range h.ActiveModels {
			if am.GetType() != input.InputType {
				continue
			}

			err := am.GetElements(input, h.DbHandler)
			if err != nil {
				msg := fmt.Sprintf("Error retrieving elements type %d", input.InputType)
				resp := Response{
					Writer:  w,
					Message: msg,
					Code:    http.StatusInternalServerError,
					Error:   err,
				}
				resp.Ko()
				return
			}
			am.Response(w)
			return
		}
		msg := fmt.Sprintf("Error retrieving type %d", input.InputType)
		resp := Response{
			Writer:  w,
			Message: msg,
			Code:    http.StatusUnprocessableEntity,
			Error:   nil,
		}
		resp.Ko()
		return
	})
}

// Put searchs for active input and call its create method
func (h *Handler) Put() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := Input{}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil && err != io.EOF {
			msg := "Error decoding input"
			resp := Response{
				Writer:  w,
				Message: msg,
				Code:    http.StatusInternalServerError,
				Error:   err,
			}
			resp.Ko()
			return
		}

		for _, am := range h.ActiveModels {
			if am.GetType() != input.InputType {
				continue
			}
			err := am.CreateElement(input, h.DbHandler)
			if err != nil {
				msg := fmt.Sprintf("Error retrieving elements type %d", input.InputType)
				resp := Response{
					Writer:  w,
					Message: msg,
					Code:    http.StatusInternalServerError,
					Error:   err,
				}
				resp.Ko()
				return
			}
			am.Response(w)
			return
		}
		msg := fmt.Sprintf("Error retrieving type %d", input.InputType)
		resp := Response{
			Writer:  w,
			Message: msg,
			Code:    http.StatusUnprocessableEntity,
			Error:   nil,
		}
		resp.Ko()
		return
	})
}

// Update searchs for active input and call its update method
func (h *Handler) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := Input{}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil && err != io.EOF {
			msg := "Error decoding project input"
			resp := Response{
				Writer:  w,
				Message: msg,
				Code:    http.StatusInternalServerError,
				Error:   err,
			}
			resp.Ko()
			return
		}

		for _, am := range h.ActiveModels {
			if am.GetType() != input.InputType {
				continue
			}

			err := am.UpdateElement(input, h.DbHandler)
			if err != nil {
				code := http.StatusInternalServerError
				if strings.Contains(err.Error(), "0 rows") {
					code = http.StatusUnprocessableEntity
				}
				resp := Response{
					Writer:  w,
					Message: err.Error(),
					Code:    code,
					Error:   err,
				}
				resp.Ko()
				return
			}
			am.Response(w)
			return
		}
		msg := fmt.Sprintf("Error retrieving type %d", input.InputType)
		resp := Response{
			Writer:  w,
			Message: msg,
			Code:    http.StatusUnprocessableEntity,
			Error:   nil,
		}
		resp.Ko()
		return
	})
}

// Delete searchs for active input and call its delete method
func (h *Handler) Delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := Input{}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil && err != io.EOF {
			msg := "Error decoding project input"
			resp := Response{
				Writer:  w,
				Message: msg,
				Code:    http.StatusInternalServerError,
				Error:   err,
			}
			resp.Ko()
			return
		}

		for _, am := range h.ActiveModels {
			if am.GetType() != input.InputType {
				continue
			}

			err := am.DeleteElement(input, h.DbHandler)
			if err != nil {
				code := http.StatusInternalServerError
				if strings.Contains(err.Error(), "0 rows") {
					code = http.StatusUnprocessableEntity
				}
				resp := Response{
					Writer:  w,
					Message: err.Error(),
					Code:    code,
					Error:   err,
				}
				resp.Ko()
				return
			}
			am.Response(w)
			return
		}
		msg := fmt.Sprintf("Error retrieving type %d", input.InputType)
		resp := Response{
			Writer:  w,
			Message: msg,
			Code:    http.StatusUnprocessableEntity,
			Error:   nil,
		}
		resp.Ko()
		return
	})
}
