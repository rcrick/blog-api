package v1

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/rcrick/blog-api/pkg/app"
	"github.com/rcrick/blog-api/pkg/e"
	"github.com/rcrick/blog-api/pkg/setting"
	"github.com/rcrick/blog-api/pkg/util"
	"github.com/rcrick/blog-api/service/tag_service"
	"github.com/unknwon/com"
	"net/http"
)

// @Summary Get multiple article tags
// @Produce  json
// @Param token header string true "Token"
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageSize: setting.PageSize,
		PageNum:  util.GetPage(c),
	}
	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}

type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Add Tag
// @Produce  json
// @Param token header string true "Token"
// @Param name query string true "Name"
// @Param state query int true "State"
// @Param created_by query string true "createdBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":400,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":10001,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddTagForm
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}
	exist, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exist {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST_TAG, nil)
		return
	}
	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" valid:"Required;Min(1)"`
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Edit Tag
// @Produce  json
// @Param token header string true "Token"
// @Param id path int true "Id"
// @Param name query string false "Name"
// @Param state query int false "State"
// @Param modified_by query string true "modifiedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":400,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":10001,"data":{},"msg":"ok"}"
// @Router /api/v1/tag/{id} [put]
func EditTag(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form = EditTagForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		State:      form.State,
		ModifiedBy: form.ModifiedBy,
	}

	exist, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	exist, err = tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exist {
		appG.Response(http.StatusBadRequest, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary Delete Tag
// @Produce  json
// @Param token header string true "Token"
// @Param id path int true "Id"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":400,"data":{},"msg":"ok"}"
// @Failure 200 {string} json "{"code":10002,"data":{},"msg":"ok"}"
// @Router /api/v1/tag/{id} [delete]
func DeleteTag(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Required(id, "id").Message("id can not be empty")
	valid.Min(id, 1, "id").Message("id must great than 0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{ID: id}
	exist, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	err = tagService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
