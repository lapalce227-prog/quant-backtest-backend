package service

import (
	"encoding/json"

	"tradingsystem/internal/dto"
	"tradingsystem/internal/model"
	"tradingsystem/internal/repository"
)

type StrategyService struct {
	repo *repository.StrategyRepository
}

func NewStrategyService(repo *repository.StrategyRepository) *StrategyService {
	return &StrategyService{repo: repo}
}

func (s *StrategyService) Create(req dto.CreateStrategyRequest) (*dto.StrategyResponse, error) {
	parameters, err := json.Marshal(req.Parameters)
	if err != nil {
		return nil, err
	}

	strategy := &model.Strategy{
		Name:         req.Name,
		Description:  req.Description,
		StrategyType: req.StrategyType,
		Parameters:   string(parameters),
	}
	if err := s.repo.Create(strategy); err != nil {
		return nil, err
	}

	return toStrategyResponse(strategy)
}

func (s *StrategyService) List() ([]dto.StrategyResponse, error) {
	strategies, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	items := make([]dto.StrategyResponse, 0, len(strategies))
	for _, strategy := range strategies {
		item, err := toStrategyResponse(&strategy)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, nil
}

func (s *StrategyService) GetByID(id uint) (*dto.StrategyResponse, error) {
	strategy, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return toStrategyResponse(strategy)
}

func (s *StrategyService) Update(id uint, req dto.UpdateStrategyRequest) (*dto.StrategyResponse, error) {
	strategy, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	parameters, err := json.Marshal(req.Parameters)
	if err != nil {
		return nil, err
	}

	strategy.Name = req.Name
	strategy.Description = req.Description
	strategy.StrategyType = req.StrategyType
	strategy.Parameters = string(parameters)

	if err := s.repo.Save(strategy); err != nil {
		return nil, err
	}

	return toStrategyResponse(strategy)
}

func (s *StrategyService) Delete(id uint) error {
	return s.repo.DeleteByID(id)
}

func toStrategyResponse(strategy *model.Strategy) (*dto.StrategyResponse, error) {
	parameters := make(map[string]interface{})
	if err := json.Unmarshal([]byte(strategy.Parameters), &parameters); err != nil {
		return nil, err
	}

	return &dto.StrategyResponse{
		ID:           strategy.ID,
		Name:         strategy.Name,
		Description:  strategy.Description,
		StrategyType: strategy.StrategyType,
		Parameters:   parameters,
		CreatedAt:    strategy.CreatedAt,
		UpdatedAt:    strategy.UpdatedAt,
	}, nil
}
