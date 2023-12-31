package main_gateways_mongodb_repositories

import (
	main_configs_mongo "baseapplicationgo/main/configs/mongodb"
	main_configs_mongo_collections "baseapplicationgo/main/configs/mongodb/collections"
	main_domains "baseapplicationgo/main/domains"
	main_gateways_mongodb_documents "baseapplicationgo/main/gateways/mongodb/documents"
	main_gateways_mongodb_utils "baseapplicationgo/main/gateways/mongodb/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

const _USERS_COLLECTION_NAME = main_configs_mongo_collections.USERS_COLLECTION_NAME
const _USERS_IDX_INDICATOR_MONGO_ID = main_configs_mongo.DEFAULT_BSON_ID_NAME

const _USER_REPO_ID = "_id"
const _USER_REPO_NAME = "name"
const _USER_REPO_DOCUMENT_NUMBER = "documentNumber"
const _USER_REPO_BIRTHDAY = "birthday"
const _USER_REPO_CREATED_DATE = "createdDate"
const _USER_REPO_LAST_MODIFIED_DATE = "lastModifiedDate"

type UserRepository struct {
	database *mongo.Database
}

func NewUserRepository() *UserRepository {
	return &UserRepository{database: main_configs_mongo.GetMongoDBClient()}
}

func (this *UserRepository) Save(user main_gateways_mongodb_documents.UserDocument) (main_gateways_mongodb_documents.UserDocument, error) {
	collection := this.database.Collection(_USERS_COLLECTION_NAME)
	indexModel := mongo.IndexModel{
		Keys: bson.D{{_USER_REPO_DOCUMENT_NUMBER, 1}},
	}
	name, err := this.database.Collection(_USERS_COLLECTION_NAME).Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		panic(err)
	}
	fmt.Println("Name of Index Created: " + name)

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return main_gateways_mongodb_documents.UserDocument{}, err
	}

	oid, _ := result.InsertedID.(primitive.ObjectID)
	user.Id = oid
	return user, nil
}

func (this *UserRepository) FindById(id string) (*main_gateways_mongodb_documents.UserDocument, error) {
	collection := this.database.Collection(_USERS_COLLECTION_NAME)
	var result main_gateways_mongodb_documents.UserDocument
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &result, nil
	}
	filter := bson.D{{_USERS_IDX_INDICATOR_MONGO_ID, objectId}}
	err2 := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err2 != nil {
		if errors.Is(err2, mongo.ErrNoDocuments) {
			return &result, nil
		}
		return &result, err2
	}
	return &result, nil
}

func (this *UserRepository) FindByDocumentNumber(documentNumber string) (*main_gateways_mongodb_documents.UserDocument, error) {
	collection := this.database.Collection(_USERS_COLLECTION_NAME)
	filter := bson.D{{_USER_REPO_DOCUMENT_NUMBER, documentNumber}}
	var result main_gateways_mongodb_documents.UserDocument
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &result, nil
		}
		return &result, err
	}
	return &result, nil
}

func (this *UserRepository) FindByFilter(filter main_domains.FindUserFilter,
	pageable main_domains.Pageable) (*main_domains.Page, error) {
	collection := this.database.Collection(_USERS_COLLECTION_NAME)

	log.Println(len(pageable.GetSort()))
	if pageable.HasEmptySort() {
		defaultSort := make(map[string]int)
		defaultSort[_USER_REPO_ID] = 1
		pageable.SetSort(defaultSort)
	}

	opt := main_gateways_mongodb_utils.NewPageableUtils().GetOptsFromPageable(pageable)

	//var filterCriterias bson.D
	//if len(filter.GetName()) > 0 {
	//	filterCriterias = append(filterCriterias,
	//		bson.E{Key: _USER_REPO_NAME, Value: bson.M{"$in": filter.GetName()}})
	//}
	//if len(filter.GetDocumentNumber()) > 0 {
	//	filterCriterias = append(filterCriterias,
	//		bson.E{Key: _USER_REPO_DOCUMENT_NUMBER, Value: bson.M{"$in": filter.GetDocumentNumber()}})
	//}
	//if len(filter.GetBirthday()) > 0 {
	//	filterCriterias = append(filterCriterias,
	//		bson.E{Key: _USER_REPO_BIRTHDAY, Value: bson.M{"$in": filter.GetBirthday()}})
	//}
	//if filter.GetStartCreatedDate().IsZero() {
	//	filterCriterias = append(filterCriterias,
	//		bson.E{Key: _USER_REPO_CREATED_DATE, Value: bson.M{"$gte": filter.GetStartCreatedDate()}})
	//}
	//if filter.GetEndCreatedDate().IsZero() {
	//	filterCriterias = append(filterCriterias,
	//		bson.E{Key: _USER_REPO_CREATED_DATE, Value: bson.M{"$lt": filter.GetEndCreatedDate()}})
	//}
	//if filter.GetStartLastModifiedDate().IsZero() {
	//	filterCriterias = append(filterCriterias,
	//		bson.E{Key: _USER_REPO_LAST_MODIFIED_DATE, Value: bson.M{"$gte": filter.GetStartLastModifiedDate()}})
	//}
	//if filter.GetEndLastModifiedDate().IsZero() {
	//	filterCriterias = append(filterCriterias,
	//		bson.E{Key: _USER_REPO_LAST_MODIFIED_DATE, Value: bson.M{"$lt": filter.GetEndLastModifiedDate()}})
	//}

	filterCriteria := bson.M{}
	if len(filter.GetName()) > 0 {
		filterCriteria[_USER_REPO_NAME] = bson.M{"$in": filter.GetName()}
	}
	if len(filter.GetDocumentNumber()) > 0 {
		filterCriteria[_USER_REPO_DOCUMENT_NUMBER] = bson.M{"$in": filter.GetDocumentNumber()}
	}

	//
	//filterCriteria := bson.M{
	//	_USER_REPO_NAME:               bson.M{"$in": filter.GetName()},
	//	_USER_REPO_DOCUMENT_NUMBER:    bson.M{"$in": filter.GetDocumentNumber()},
	//	_USER_REPO_BIRTHDAY:           bson.M{"$in": filter.GetBirthday()},
	//	_USER_REPO_CREATED_DATE:       bson.M{"eventDateTime": bson.M{"$gte": filter.GetStartCreatedDate(), "$lt": filter.GetEndCreatedDate()}},
	//	_USER_REPO_LAST_MODIFIED_DATE: bson.M{"eventDateTime": bson.M{"$gte": filter.GetStartLastModifiedDate(), "$lt": filter.GetEndLastModifiedDate()}},
	//}

	//if len(filter.GetDocumentNumber()) > 0 {
	//	filterCriterias = append(filterCriterias,
	//		bson.E{Key: _USER_REPO_DOCUMENT_NUMBER, Value: bson.M{"$in": filter.GetDocumentNumber()}})
	//}
	//if len(filter.GetBirthday()) > 0 {
	//	filterCriterias = append(filterCriterias,
	//		bson.E{Key: _USER_REPO_BIRTHDAY, Value: bson.M{"$in": filter.GetBirthday()}})
	//}
	//if filter.GetStartCreatedDate().IsZero() {
	//	filterCriterias = append(filterCriterias,
	//		bson.E{Key: _USER_REPO_CREATED_DATE, Value: bson.M{"$gte": filter.GetStartCreatedDate()}})
	//}
	//if filter.GetEndCreatedDate().IsZero() {
	//	filterCriterias = append(filterCriterias,
	//		bson.E{Key: _USER_REPO_CREATED_DATE, Value: bson.M{"$lt": filter.GetEndCreatedDate()}})
	//}
	//if filter.GetStartLastModifiedDate().IsZero() {
	//	filterCriterias = append(filterCriterias,
	//		bson.E{Key: _USER_REPO_LAST_MODIFIED_DATE, Value: bson.M{"$gte": filter.GetStartLastModifiedDate()}})
	//}
	//if filter.GetEndLastModifiedDate().IsZero() {
	//	filterCriterias = append(filterCriterias,
	//		bson.E{Key: _USER_REPO_LAST_MODIFIED_DATE, Value: bson.M{"$lt": filter.GetEndLastModifiedDate()}})
	//}

	var results []main_gateways_mongodb_documents.UserDocument
	cursor, err := collection.Find(context.TODO(), filterCriteria, opt)
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		res, _ := json.Marshal(result)
		fmt.Println(string(res))
	}

	numberDocs, err := collection.CountDocuments(context.TODO(), filterCriteria)
	if err != nil {
		return nil, err
	}

	var contents []main_domains.Content
	for _, value := range results {
		contents = append(contents, *main_domains.NewContent(value))
	}

	return main_domains.NewPage(contents, pageable.GetPage(), pageable.GetSize(), numberDocs), nil
}
