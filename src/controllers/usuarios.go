package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"PI-TATOVERING/src/models"
)

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	var albums = []album{
		{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
    c.IndentedJSON(http.StatusOK, albums)
}
// type TagController struct {
// 	tagService service.TagsService
// }

// func NewTagController(service service.TagsService) *TagController {
// 	return &TagController{tagService: service}
// }

// func (controller *TagController) Create(ctx *gin.Context) {
// 	createTagRequest := request.CreateTagsRequest{}
// 	err := ctx.ShouldBindJSON(&createTagRequest)
// 	helper.ErrorPanic(err)

// 	controller.tagService.Create(createTagRequest)

// 	webResponse := response.Response{
// 		Code:   200,
// 		Status: "Ok",
// 		Data:   nil,
// 	}

// 	ctx.JSON(http.StatusOK, webResponse)
// }

// func (controller *TagController) Update(ctx *gin.Context) {
// 	updateTagRequest := request.UpdateTagsRequest{}
// 	err := ctx.ShouldBindJSON(&updateTagRequest)
// 	helper.ErrorPanic(err)

// 	tagId := ctx.Param("tagId")
// 	id, err := strconv.Atoi(tagId)
// 	helper.ErrorPanic(err)

// 	updateTagRequest.Id = id

// 	controller.tagService.Update(updateTagRequest)

// 	webResponse := response.Response{
// 		Code:   200,
// 		Status: "Ok",
// 		Data:   nil,
// 	}

// 	ctx.JSON(http.StatusOK, webResponse)
// }

// func (controller *TagController) Delete(ctx *gin.Context) {
// 	tagId := ctx.Param("tagId")
// 	id, err := strconv.Atoi(tagId)
// 	helper.ErrorPanic(err)
// 	controller.tagService.Delete(id)

// 	webResponse := response.Response{
// 		Code:   200,
// 		Status: "Ok",
// 		Data:   nil,
// 	}

// 	ctx.JSON(http.StatusOK, webResponse)

// }

// func (controller *TagController) FindById(ctx *gin.Context) {
// 	tagId := ctx.Param("tagId")
// 	id, err := strconv.Atoi(tagId)
// 	helper.ErrorPanic(err)

// 	tagResponse := controller.tagService.FindById(id)

// 	webResponse := response.Response{
// 		Code:   200,
// 		Status: "Ok",
// 		Data:   tagResponse,
// 	}
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)
// }

// func (controller *TagController) FindAll(ctx *gin.Context) {
// 	tagResponse := controller.tagService.FindAll()

// 	webResponse := response.Response{
// 		Code:   200,
// 		Status: "Ok",
// 		Data:   tagResponse,
// 	}
// 	ctx.Header("Content-Type", "application/json")
// 	ctx.JSON(http.StatusOK, webResponse)

// }