package forms

import "twitta/pkg/utils"

type SearchForm struct {
	utils.PageForm
	Keyword string `form:"keyword"`
}
