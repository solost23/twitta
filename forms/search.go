package forms

import "Twitta/pkg/utils"

type SearchForm struct {
	utils.PageForm
	Keyword string `form:"keyword"`
}
