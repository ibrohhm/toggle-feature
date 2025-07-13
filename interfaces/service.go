package interfaces

import (
	"github.com/toggle-feature/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ToggleFeatureServiceInterface interface {
	SelectAll(names []string) ([]entity.ToggleFeature, error)
	Select(id primitive.ObjectID) (entity.ToggleFeature, error)
	Insert(toggleFeatureParams entity.ToggleFeatureParams) (entity.ToggleFeature, error)
	Update(id primitive.ObjectID, toggleFeatureParams entity.ToggleFeatureParams) (entity.ToggleFeature, error)
	Delete(id primitive.ObjectID) error
}
