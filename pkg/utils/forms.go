package utils

type UIdForm struct {
	Id string `uri:"id" comment:"id" binding:"min=1"`
}
