package service

import (
	"fmt"
	"main/repo"
	"main/transformer"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	db    DB
	cache Cache
}

type DB interface {
	NewRecord(r *repo.Record) error
	GetRecord(id string) (repo.Record, error)
	GetAllRecords() ([]repo.Record, error)
	UpdateRecord(r *repo.Record) error
	DeleteRecord(id string) error
}

type Cache interface {
	Set(value *repo.Record)
	Get(key string) (*repo.Record, bool)
	Delete(key string)
}

func NewService(db DB, cache Cache) Service {
	return Service{
		db:    db,
		cache: cache,
	}
}

func (s Service) CreateRecord(request repo.TransformRequest) (*repo.Record, error) {
	var tr transformer.Transformer

	switch {
	case request.Type == "reverse":
		tr = transformer.NewReverseTransformer()
	case request.Type == "caesar":
		tr = transformer.NewCaesarTransformer(request.CaesarShift)
	case request.Type == "base64":
		tr = transformer.NewBase64Transformer()
	}

	res, err := tr.Transform(strings.NewReader(request.Input), false)
	if err != nil {
		return nil, fmt.Errorf("service error: %w", err)
	}

	result := &repo.Record{
		ID:          uuid.NewString(),
		Type:        request.Type,
		CaesarShift: request.CaesarShift,
		Result:      res,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   0,
	}

	err = s.db.NewRecord(result)
	if err != nil {
		return nil, fmt.Errorf("service error: %w", err)
	}

	return result, nil
}

func (s Service) DeleteRecord(id string) error {
	err := s.db.DeleteRecord(id)
	if err != nil {
		return fmt.Errorf("service error: %w", err)
	}

	s.cache.Delete(id)

	return nil
}

func (s Service) GetAllRecords() ([]repo.Record, error) {
	values, err := s.db.GetAllRecords()
	if err != nil {
		return nil, fmt.Errorf("service error: %w", err)
	}

	return values, nil
}

func (s Service) GetRecord(id string) (*repo.Record, error) {
	res, ok := s.cache.Get(id)
	if ok {
		return res, nil
	}

	result, err := s.db.GetRecord(id)
	if err != nil {
		return nil, fmt.Errorf("service error: %w", err)
	}

	s.cache.Set(&result)

	return &result, nil
}

func (s Service) UpdateRecord(id string, request repo.TransformRequest) *repo.Record {
	var tr transformer.Transformer

	switch {
	case request.Type == "reverse":
		tr = transformer.NewReverseTransformer()
	case request.Type == "caesar":
		tr = transformer.NewCaesarTransformer(request.CaesarShift)
	case request.Type == "base64":
		tr = transformer.NewBase64Transformer()
	}

	TransformResult, err := tr.Transform(strings.NewReader(request.Input), false)
	if err != nil {
		return nil
	}

	result, err := s.db.GetRecord(id)
	result.Type = request.Type
	result.CaesarShift = request.CaesarShift
	result.Result = TransformResult

	if err != nil {
		result.ID = id
		result.CreatedAt = time.Now().Unix()

		err = s.db.NewRecord(&result)
		if err != nil {
			return nil
		}

		return &result
	}

	result.UpdatedAt = time.Now().Unix()

	err = s.db.UpdateRecord(&result)
	if err != nil {
		return nil
	}

	s.cache.Set(&result)

	return &result
}
