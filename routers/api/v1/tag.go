package v1

import (
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/rcrick/blog-api/models"
	"github.com/rcrick/blog-api/pkg/e"
	"github.com/rcrick/blog-api/pkg/logging"
	"github.com/rcrick/blog-api/pkg/setting"
	"github.com/rcrick/blog-api/pkg/util"
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
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}
	code := e.SUCCESS
	data["lists"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
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
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("name can not be empty")
	valid.MaxSize(name, 100, "name").Message("name's length can not over 100")
	valid.Required(createdBy, "created_by").Message("created_by can not be empty")
	valid.MaxSize(createdBy, 100, "created_by").Message("created_by's length can not over 100")
	valid.Range(state, 0, 1, "state").Message("state only can be 0 or 1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistTagByName(name) {
			if models.AddTag(name, state, createdBy) {
				code = e.SUCCESS
			} else {
				code = e.ERROR
			}
		} else {
			code = e.ERROR_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
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
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("state only can be 0 or 1")
	}

	valid.Required(id, "id").Message("id can not be empty")
	valid.Required(modifiedBy, "modified_by").Message("modified_by name can not be empty")
	valid.MaxSize(modifiedBy, 100, "modifiedBy").Message("modifiedBy's length can not over 100")
	valid.MaxSize(name, 100, "name").Message("name's length can not over 100")

	code := e.INVALID_PARAMS
	fmt.Println(valid.Errors)
	if !valid.HasErrors() {
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}
			if models.EditTag(id, data) {
				code = e.SUCCESS
			}
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
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
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Required(id, "id").Message("id can not be empty")
	valid.Min(id, 1, "id").Message("id must great than 0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistTagByID(id) {
			code = e.ERROR_NOT_EXIST_TAG
		} else {
			if models.DeleteTag(id) {
				code = e.SUCCESS
			}
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
