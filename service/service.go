package service

import (
	"fmt"
	"main/cache"
	"main/repo"
	"main/transformer"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ServiceInterface interface {
	NewRecord(request TransformRequest) *repo.Record
	GetRecord(id string) (*repo.Record, error)
	GetRecords() ([]repo.Record, error)
	UpdateRecord(id string, request TransformRequest) *repo.Record
	DeleteRecord(id string) error
}

type service struct {
	db    DBLayer
	cache cache.CacheInterface
}

type DBLayer interface {
	NewRecord(r *repo.Record) error
	GetRecord(id string) (repo.Record, error)
	GetRecords() ([]repo.Record, error)
	UpdateRecord(r *repo.Record) error
	DeleteRecord(id string) error
}

type TransformRequest struct {
	Type        string `json:"type"`
	CaesarShift int    `json:"shift,omitempty"`
	Input       string `json:"input,omitempty"`
}

func NewService(db DBLayer, cache cache.CacheInterface) service {
	return service{
		db:    db,
		cache: cache,
	}
}

func (s service) NewRecord(request TransformRequest) *repo.Record {

	result := new(repo.Record)
	result.ID = uuid.NewString()
	result.Type = request.Type
	result.CaesarShift = request.CaesarShift
	var tr transformer.Transformer
	switch {
	case request.Type == "reverse":
		tr = transformer.NewReverseTransformer()
	case request.Type == "caesar":
		tr = transformer.NewCaesarTransformer(request.CaesarShift)
	case request.Type == "base64":
		tr = transformer.NewBase64Transformer()
	}
	var err error
	result.Result, err = tr.Transform(strings.NewReader(request.Input), false)
	if err != nil {
		fmt.Errorf("service error transforming")
		return nil
	}
	result.CreatedAt = time.Now().Unix()

	err = s.db.NewRecord(result)
	if err != nil {
		fmt.Errorf("service error creating new record")
		return nil
	}
	return result
}

func (s service) DeleteRecord(id string) error {
	err := s.db.DeleteRecord(id)
	if err != nil {
		fmt.Errorf("service error deleting record")
		return err
	}
	return nil
}

func (s service) GetRecords() ([]repo.Record, error) {
	values, err := s.db.GetRecords()
	if err != nil {
		fmt.Errorf("service error getting all records")
		return nil, err
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i].CreatedAt > values[j].CreatedAt
	})

	return values, nil
}

func (s service) GetRecord(id string) (*repo.Record, error) {
	res := s.cache.Get(id)
	if res != nil {
		return res, nil
	}
	result, err := s.db.GetRecord(id)
	if err != nil {
		fmt.Errorf("service error getting one record")
		return nil, err
	}
	s.cache.Set(&result)
	return &result, nil
}

func (s service) UpdateRecord(id string, request TransformRequest) *repo.Record {

	var tr transformer.Transformer
	switch {
	case request.Type == "reverse":
		tr = transformer.NewReverseTransformer()
	case request.Type == "caesar":
		tr = transformer.NewCaesarTransformer(request.CaesarShift)
	case request.Type == "base64":
		tr = transformer.NewBase64Transformer()
	}
	transform_result, err := tr.Transform(strings.NewReader(request.Input), false)
	if err != nil {
		fmt.Errorf("service error (updating) transforming")
		return nil
	}

	result, err := s.db.GetRecord(id)
	result.Type = request.Type
	result.CaesarShift = request.CaesarShift
	result.Result = transform_result
	result.UpdatedAt = time.Now().Unix()
	if err != nil {
		result.ID = id
		result.UpdatedAt = 0
		result.CreatedAt = time.Now().Unix()
		err = s.db.NewRecord(&result)
		if err != nil {
			fmt.Errorf("service error creating(update) record")
			return nil
		}
	} else {
		err = s.db.UpdateRecord(&result)
		if err != nil {
			fmt.Errorf("service error updating record")
			return nil
		}
	}
	s.cache.Set(&result)

	return &result
}
