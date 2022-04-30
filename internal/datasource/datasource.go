package datasource

import (
	"context"
	"errors"
	"fmt"
	"log"
)

var (
	ErrFetchingDatasource = errors.New("failed fetching datasource by id")
	ErrNotImplemented     = errors.New("not implemented")
)

// a representation of the comment strcture for our service
type Datasource struct {
	ID   string
	Name string
	URL  string
}

// Store - this interface defines all of the methods that our service needs in order to operate
type Store interface {
	GetDatasource(context.Context, string) (Datasource, error)
	PostDatasource(context.Context, Datasource) (Datasource, error)
	DeleteDatasource(context.Context, string) error
	UpdateDatasource(context.Context, string, Datasource) (Datasource, error)
}

// is the struct on which all our logic will be built on to of
type Service struct {
	Store Store
}

// returns a pointer to a new service
func NewService(store Store) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetDatasource(ctx context.Context, id string) (Datasource, error) {
	log.Println("Retrieving Datasource")

	ds, err := s.Store.GetDatasource(ctx, id)
	if err != nil {
		// log detailed error
		fmt.Println(err)
		// return only a "meaningful" error to the user
		return Datasource{}, ErrFetchingDatasource
	}

	return ds, nil
}

func (s *Service) UpdateDatasource(ctx context.Context, ID string, updatedDs Datasource) (Datasource, error) {
	ds, err := s.Store.UpdateDatasource(ctx, ID, updatedDs)
	if err != nil {
		fmt.Println("Error updating datasource")
		return Datasource{}, err
	}
	return ds, nil
}

func (s *Service) DeleteDatasource(ctx context.Context, id string) error {
	return s.Store.DeleteDatasource(ctx, id)
}

func (s *Service) PostDatasource(ctx context.Context, ds Datasource) (Datasource, error) {
	insertedDs, err := s.Store.PostDatasource(ctx, ds)
	if err != nil {
		return Datasource{}, err
	}
	return insertedDs, nil
}
