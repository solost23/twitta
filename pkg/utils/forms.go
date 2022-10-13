package utils

type UIdForm struct {
	Id string `uri:"id" comment:"id" binding:"min=1"`
}

type PageForm struct {
	Page int `form:"page" comment:"当前页码"`
	Size int `form:"size" comment:"每页显示记录数"`
}
