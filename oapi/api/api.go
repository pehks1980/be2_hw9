package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Api struct {
	//pool       *pgxpool.Pool
	repository *Repo
}

// Make sure we conform to ServerInterface
var _ ServerInterface = (*Api)(nil)

func NewApi(ctx context.Context, dbconn string) (*Api, error) {

	newRepository, err := NewRepository(ctx, dbconn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	err = newRepository.InitSchema(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create DDL in database: %w", err)
	}

	return &Api{
		//pool:       pool,
		repository: newRepository,
	}, nil
}

func (a *Api) PricelistEntityCreateView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newPriceListEntity Pricelist
	if err := json.NewDecoder(r.Body).Decode(&newPriceListEntity); err != nil {
		//sendPetstoreError(w, http.StatusBadRequest, "Invalid format for NewPet")
		return
	}
	id, err := a.repository.AddPriceListEntity(ctx, newPriceListEntity)
	if err != nil {
		log.Printf("add repo error: %v", err)
		return
	}
	*newPriceListEntity.Id = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPriceListEntity)
}

func (a *Api) PricelistDeleteView(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	err := a.repository.DelPriceList(ctx, id)
	if err != nil {
		log.Printf("delete repo error: %v", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *Api) PricelistListView(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	var result []Pricelist

	result, _ = a.repository.GetPriceList(ctx, id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (a *Api) PricelistUpdateView(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	var updPriceListEntity Pricelist
	if err := json.NewDecoder(r.Body).Decode(&updPriceListEntity); err != nil {
		//sendPetstoreError(w, http.StatusBadRequest, "Invalid format for NewPet")
		return
	}
	err := a.repository.UpdatePriceList(ctx, id, updPriceListEntity)
	if err != nil {
		log.Printf("add repo error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
