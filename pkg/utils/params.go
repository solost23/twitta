package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var Trans ut.Translator

func DefaultGetValidParams(c *gin.Context, params interface{}) error {
	if err := InitTrans("zh"); err != nil {
		return err
	}
	if err := c.ShouldBind(params); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}
		return errors.New(removeTopStruct(errs.Translate(Trans)))
	}
	return nil
}

func removeTopStruct(fields map[string]string) (result string) {
	for _, err := range fields {
		result += err + ","
	}
	return result
}

func GetValidUriParams(c *gin.Context, params interface{}) error {
	if err := InitTrans("zh"); err != nil {
		return err
	}
	if err := c.ShouldBindUri(params); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}
		return errors.New(removeTopStruct(errs.Translate(Trans)))
	}
	return nil
}
