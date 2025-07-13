package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/toggle-feature/entity"
	"github.com/toggle-feature/interfaces"
	"github.com/toggle-feature/utility/response"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ToggleFeatureHandler struct {
	Service interfaces.ToggleFeatureServiceInterface
}

func NewToggleFeatureHandler(service interfaces.ToggleFeatureServiceInterface) *ToggleFeatureHandler {
	return &ToggleFeatureHandler{
		Service: service,
	}
}

func (h *ToggleFeatureHandler) Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	query := r.URL.Query()
	toggle_features, err := h.Service.SelectAll(query["names"])
	if err != nil {
		return response.WriteError(w, err)
	}

	return response.WriteSuccess(w, entity.ToggleFeatureParser(toggle_features), "")
}

func (h *ToggleFeatureHandler) Get(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	id, err := primitive.ObjectIDFromHex(params.ByName("id"))
	if err != nil {
		return response.WriteError(w, err)
	}

	toggleFeature, err := h.Service.Select(id)
	if err != nil {
		return response.WriteError(w, err)
	}

	return response.WriteSuccess(w, toggleFeature.Parser(), "")
}

func (h *ToggleFeatureHandler) Insert(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return response.WriteError(w, err)
	}

	var toggleFeatureParams entity.ToggleFeatureParams
	err = json.Unmarshal(b, &toggleFeatureParams)
	if err != nil {
		return response.WriteError(w, err)
	}

	objectID, err := h.Service.Insert(toggleFeatureParams)
	if err != nil {
		return response.WriteError(w, err)
	}

	return response.WriteSuccess(w, objectID, "")
}

func (h *ToggleFeatureHandler) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	id, err := primitive.ObjectIDFromHex(params.ByName("id"))
	if err != nil {
		return response.WriteError(w, err)
	}

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return response.WriteError(w, err)
	}

	var toggleFeatureParams entity.ToggleFeatureParams
	err = json.Unmarshal(b, &toggleFeatureParams)
	if err != nil {
		return response.WriteError(w, err)
	}

	toggleFeature, err := h.Service.Update(id, toggleFeatureParams)
	if err != nil {
		return response.WriteError(w, err)
	}

	return response.WriteSuccess(w, toggleFeature.Parser(), "")
}

func (h *ToggleFeatureHandler) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	id, err := primitive.ObjectIDFromHex(params.ByName("id"))
	if err != nil {
		return response.WriteError(w, err)
	}

	err = h.Service.Delete(id)
	if err != nil {
		return response.WriteError(w, err)
	}

	return response.WriteSuccess(w, nil, "success delete")
}
