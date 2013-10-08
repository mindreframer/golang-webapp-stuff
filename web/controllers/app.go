package controllers

func New() *AppController {
	c := new(AppController)
	c.Init()

	return c
}
