package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type ToggleFeature struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description,omitempty"`
	Active      bool               `bson:"active"`
}

type ToggleFeatureParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type ToggleFeatureResponse struct {
	ID          primitive.ObjectID `json:"_id,omitempty"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Active      bool               `json:"active"`
}

func (toggle_feature *ToggleFeature) Parser() ToggleFeatureResponse {
	return ToggleFeatureResponse{
		ID:          toggle_feature.ID,
		Name:        toggle_feature.Name,
		Description: toggle_feature.Description,
		Active:      toggle_feature.Active,
	}
}

func ToggleFeatureParser(toggle_features []ToggleFeature) []ToggleFeatureResponse {
	toggle_feature_responses := []ToggleFeatureResponse{}
	for _, toggle_feature := range toggle_features {
		toggle_feature_responses = append(toggle_feature_responses, toggle_feature.Parser())
	}

	return toggle_feature_responses
}
