package repository

import (
	"context"
	"log"
	"net/http"

	"github.com/toggle-feature/entity"
	customErr "github.com/toggle-feature/utility/error"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var DatabaseName = "toggle_features"

type ToggleFeatureRepository struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewToggleFeatureRepository(client *mongo.Client) *ToggleFeatureRepository {
	collection := client.Database(DatabaseName).Collection(DatabaseName)
	createUniqueIndex(collection)

	return &ToggleFeatureRepository{
		client: client,
		coll:   collection,
	}
}

// createUniqueIndex is method set unique index for name
func createUniqueIndex(coll *mongo.Collection) {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"name": 1},               // Field to index (1 for ascending order)
		Options: options.Index().SetUnique(true), // Set the index as unique
	}

	indexName, err := coll.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}
	log.Printf("Created index: %s", indexName)
}

func (r *ToggleFeatureRepository) SelectAll(names []string) ([]entity.ToggleFeature, error) {
	filter := bson.M{}
	if len(names) != 0 {
		filter = bson.M{"name": bson.M{"$in": names}}
	}

	cursor, err := r.coll.Find(context.TODO(), filter)
	if err != nil {
		return []entity.ToggleFeature{}, customErr.New(err.Error(), http.StatusUnprocessableEntity)
	}

	var results []entity.ToggleFeature
	if err = cursor.All(context.TODO(), &results); err != nil {
		return []entity.ToggleFeature{}, customErr.New(err.Error(), http.StatusInternalServerError)
	}

	return results, nil
}

func (r *ToggleFeatureRepository) Select(id primitive.ObjectID) (entity.ToggleFeature, error) {
	filter := bson.M{"_id": id}

	var toggleFeature entity.ToggleFeature
	err := r.coll.FindOne(context.TODO(), filter).Decode(&toggleFeature)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.ToggleFeature{}, customErr.New("data not found", http.StatusNotFound)
		}

		return entity.ToggleFeature{}, customErr.New(err.Error(), http.StatusInternalServerError)
	}

	return toggleFeature, nil
}

func (r *ToggleFeatureRepository) Insert(toggleFeatureParams entity.ToggleFeatureParams) (entity.ToggleFeature, error) {
	newToggleFeature := entity.ToggleFeature{Name: toggleFeatureParams.Name, Description: toggleFeatureParams.Description, Active: toggleFeatureParams.Active}

	result, err := r.coll.InsertOne(context.TODO(), newToggleFeature)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return entity.ToggleFeature{}, customErr.New("toggle name already exist", http.StatusUnprocessableEntity)
		} else {
			return entity.ToggleFeature{}, customErr.New(err.Error(), http.StatusInternalServerError)
		}
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return entity.ToggleFeature{}, customErr.New("InsertedID is not of type ObjectID", http.StatusInternalServerError)
	}

	return entity.ToggleFeature{
		ID:          objectID,
		Name:        toggleFeatureParams.Name,
		Description: toggleFeatureParams.Description,
		Active:      toggleFeatureParams.Active,
	}, nil
}

func (r *ToggleFeatureRepository) Update(id primitive.ObjectID, toggleFeature entity.ToggleFeature) (entity.ToggleFeature, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": toggleFeature}
	_, err := r.coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return entity.ToggleFeature{}, customErr.New("toggle name already exist", http.StatusInternalServerError)
		} else {
			return entity.ToggleFeature{}, customErr.New(err.Error(), http.StatusInternalServerError)
		}
	}

	return toggleFeature, nil
}

func (r *ToggleFeatureRepository) Delete(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return customErr.New(err.Error(), http.StatusInternalServerError)
	}

	return nil
}
