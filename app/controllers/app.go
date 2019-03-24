package controllers

import (
	"github.com/revel/revel"
)

//App is the main controller of the application
type App struct {
	*revel.Controller
}

//Index serves the index page of the application
func (c App) Index() revel.Result {
	return c.Render()
}
