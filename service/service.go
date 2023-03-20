package service

import (
	"fmt"
	"main/models"
	"main/transformer"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	db    DB
	cache Cache
}

type DB interface {
	CreateRecord(r *models.Record) error
	GetRecord(id string) (models.Record, error)
	GetAllRecords() ([]models.Record, error)
	UpdateRecord(r *models.Record) error
	DeleteRecord(id string) error
}

type Cache interface {
	Set(value *models.Record)
	Get(key string) (*models.Record, bool)
	Delete(key string)
}

func NewService(db DB, cache Cache) Service {
	return Service{
		db:    db,
		cache: cache,
	}
}

func SwitchAndTransform(request models.TransformRequest) (string, error) {
	var tr transformer.Transformer

	switch request.Type {
	case "reverse":
		tr = transformer.NewReverseTransformer()
	case "caesar":
		tr = transformer.NewCaesarTransformer(request.CaesarShift)
	case "base64":
		tr = transformer.NewBase64Transformer()
	}

	TransformResult, err := tr.Transform(strings.NewReader(request.Input), false)
	if err != nil {
		return "", fmt.Errorf("service error: %w", err)
	}

	return TransformResult, nil
}

func (s Service) CreateRecord(request models.TransformRequest) (*models.Record, error) {
	TransformResult, err := SwitchAndTransform(request)
	if err != nil {
		return nil, err
	}

	result := &models.Record{
		ID:          uuid.NewString(),
		Type:        request.Type,
		CaesarShift: request.CaesarShift,
		Result:      TransformResult,
	}

	err = s.db.CreateRecord(result)
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

func (s Service) GetAllRecords() ([]models.Record, error) {
	values, err := s.db.GetAllRecords()
	if err != nil {
		return nil, fmt.Errorf("service error: %w", err)
	}

	return values, nil
}

func (s Service) GetRecord(id string) (*models.Record, error) {
	res, ok := s.cache.Get(id)
	if ok {
		return res, nil
	}

	result, err := s.db.GetRecord(id)
	if err != nil {
		return nil, fmt.Errorf("service error: %w", err)
	}

	//nolint:nilnil
	if result.ID == "" {
		return nil, nil
	}

	s.cache.Set(&result)

	return &result, nil
}

func (s Service) UpdateRecord(id string, request models.TransformRequest) *models.Record {
	TransformResult, err := SwitchAndTransform(request)
	if err != nil {
		return nil
	}

	result, _ := s.db.GetRecord(id)
	result.Type = request.Type
	result.CaesarShift = request.CaesarShift
	result.Result = TransformResult

	if result.ID == "" {
		result.ID = id

		err = s.db.CreateRecord(&result)
		if err != nil {
			return nil
		}

		return &result
	}

	err = s.db.UpdateRecord(&result)
	if err != nil {
		return nil
	}

	res, err := s.db.GetRecord(result.ID)
	if err != nil {
		return nil
	}

	s.cache.Set(&res)

	return &res
}
