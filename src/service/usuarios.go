package service

import (
	"../models/usuarios"
	"github.com/go-playground/validator/v10"
	"practice-api-gin-one/data/request"
	"practice-api-gin-one/data/response"
	"practice-api-gin-one/helper"
	"practice-api-gin-one/repository"
)

type TagsService interface {
	Create(tags request.CreateTagsRequest)
	Update(tags request.UpdateTagsRequest)
	Delete(tagsId int)
	FindById(tagsId int) response.TagsResponse
	FindAll() []response.TagsResponse
}
type TagsServiceImpl struct {
	TagRepository repository.TagsRepository
	Validate      *validator.Validate
}

func NewTagServiceImpl(tagRepository repository.TagsRepository, validate *validator.Validate) TagsService {
	return &TagsServiceImpl{
		TagRepository: tagRepository,
		Validate:      validate,
	}
}

func (t TagsServiceImpl) Create(tag request.CreateTagsRequest) {
	err := t.Validate.Struct(tag)
	helper.ErrorPanic(err)
	tagModel := model.Tags{
		Name: tag.Name,
	}
	t.TagRepository.Save(tagModel)
}

func (t TagsServiceImpl) Update(tag request.UpdateTagsRequest) {
	tagData, err := t.TagRepository.FindById(tag.Id)
	helper.ErrorPanic(err)
	tagData.Name = tag.Name
	t.TagRepository.Update(tagData)
}

func (t TagsServiceImpl) Delete(tagId int) {
	t.TagRepository.Delete(tagId)
}

func (t TagsServiceImpl) FindById(tagId int) response.TagsResponse {
	tagData, err := t.TagRepository.FindById(tagId)
	helper.ErrorPanic(err)

	tagResponse := response.TagsResponse{
		Id:   tagData.Id,
		Name: tagData.Name,
	}
	return tagResponse
}

func (t TagsServiceImpl) FindAll() []response.TagsResponse {
	result := t.TagRepository.FindAll()

	var tags []response.TagsResponse
	for _, value := range result {
		tag := response.TagsResponse{
			Id:   value.Id,
			Name: value.Name,
		}
		tags = append(tags, tag)
	}
	return tags
}

