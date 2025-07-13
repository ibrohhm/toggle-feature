package service

import (
	"github.com/toggle-feature/entity"
	"github.com/toggle-feature/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ToggleFeatureService struct {
	Repo interfaces.ToggleFeatureRepositoryInterface
}

func NewToggleFeatureService(repo interfaces.ToggleFeatureRepositoryInterface) *ToggleFeatureService {
	return &ToggleFeatureService{
		Repo: repo,
	}
}

func (s *ToggleFeatureService) SelectAll(names []string) ([]entity.ToggleFeature, error) {
	return s.Repo.SelectAll(names)
}

func (s *ToggleFeatureService) Select(id primitive.ObjectID) (entity.ToggleFeature, error) {
	return s.Repo.Select(id)
}

func (s *ToggleFeatureService) Insert(toggleFeatureParams entity.ToggleFeatureParams) (entity.ToggleFeature, error) {
	return s.Repo.Insert(toggleFeatureParams)
}

func (s *ToggleFeatureService) Update(id primitive.ObjectID, toggleFeatureParams entity.ToggleFeatureParams) (entity.ToggleFeature, error) {
	oldToggleFeature, err := s.Repo.Select(id)
	if err != nil {
		return entity.ToggleFeature{}, err
	}

	oldToggleFeature.Name = toggleFeatureParams.Name
	oldToggleFeature.Description = toggleFeatureParams.Description
	oldToggleFeature.Active = toggleFeatureParams.Active
	return s.Repo.Update(oldToggleFeature.ID, oldToggleFeature)
}

func (s *ToggleFeatureService) Delete(id primitive.ObjectID) error {
	return s.Repo.Delete(id)
}
