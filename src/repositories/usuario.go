package repository

import (
	"practice-api-gin-one/model"
)

type TagsRepository interface {
	Save(tags model.Tags)
	Update(tags model.Tags)
	Delete(tagsId int)
	FindById(tagsId int) (tags model.Tags, err error)
	FindAll() []model.Tags
}

type usuario interface {
	Save(tags model.Usuarios)
	Update(tags model.Usuarios)
	Delete(tagsId int)
	FindById(tagsId int) (tags model.Usuarios, err error)
	FindAll() []model.Usuarios
}