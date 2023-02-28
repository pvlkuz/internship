package service

import (
	"main/repo"
	"main/transformer"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Service struct {
	db    DBLayer
	cache CacheInterface
}

type DBLayer interface {
	NewRecord(r *repo.Record) error
	GetRecord(id string) (repo.Record, error)
	GetAllRecords() ([]repo.Record, error)
	UpdateRecord(r *repo.Record) error
	DeleteRecord(id string) error
}

type CacheInterface interface {
	Set(value *repo.Record)
	Get(key string) (*repo.Record, bool)
	Delete(key string)
}

func NewService(db DBLayer, cache CacheInterface) Service {
	return Service{
		db:    db,
		cache: cache,
	}
}

func (s Service) CreateRecord(request repo.TransformRequest) *repo.Record {
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
		return nil
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
		return nil
	}

	return result
}

func (s Service) DeleteRecord(id string) error {
	err := s.db.DeleteRecord(id)
	if err != nil {
		return errors.Wrap(err, "service error while deleting")
	}

	return nil
}

func (s Service) GetRecords() ([]repo.Record, error) {
	values, err := s.db.GetAllRecords()
	if err != nil {
		return nil, errors.Wrap(err, "service error while reading all records")
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
		return nil, errors.Wrap(err, "service error while reading one record")
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
