package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bueti/yeoman/internal/datasource"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type DatasourceService interface {
	PostDatasource(context.Context, datasource.Datasource) (datasource.Datasource, error)
	GetDatasource(ctx context.Context, ID string) (datasource.Datasource, error)
	UpdateDatasource(ctx context.Context, ID string, newCmt datasource.Datasource) (datasource.Datasource, error)
	DeleteDatasource(ctx context.Context, ID string) error
}

type Response struct {
	Message string
}

type PostDatasourceRequest struct {
	Name string `json:"name" validate:"required"`
	URL  string `json:"url" validate:"required"`
}

func convertPostDatasourceRequestToDatasource(c PostDatasourceRequest) datasource.Datasource {
	return datasource.Datasource{
		Name: c.Name,
		URL:  c.URL,
	}
}

func (h *Handler) PostDatasource(w http.ResponseWriter, r *http.Request) {
	var cmt PostDatasourceRequest
	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		return
	}

	validator := validator.New()
	err := validator.Struct(cmt)
	if err != nil {
		http.Error(w, "not a valid datasource", http.StatusBadRequest)
		return
	}

	convertedDatasource := convertPostDatasourceRequestToDatasource(cmt)
	postedCmt, err := h.Service.PostDatasource(r.Context(), convertedDatasource)
	if err != nil {
		log.Print(err)
		return
	}

	if err = json.NewEncoder(w).Encode(postedCmt); err != nil {
		panic(err)
	}
}

func (h *Handler) GetDatasource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmt, err := h.Service.GetDatasource(r.Context(), id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cmt); err != nil {
		panic(err)
	}
}

func (h *Handler) UpdateDatasource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var cmt datasource.Datasource
	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		return
	}

	cmt, err := h.Service.UpdateDatasource(r.Context(), id, cmt)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(cmt); err != nil {
		panic(err)
	}

}

func (h *Handler) DeleteDatasource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteDatasource(r.Context(), id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(Response{Message: "Successfully deleted."}); err != nil {
		panic(err)
	}
}
