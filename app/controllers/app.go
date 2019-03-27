package controllers

import (
	"encoding/json"
	"strings"

	"github.com/revel/revel"
	"github.com/shredx/golang-redis-rate-limiter/app/models"
)

//App is the main controller of the application
type App struct {
	*revel.Controller
}

//Index serves the index page of the application
func (c App) Index() revel.Result {
	return c.Render()
}

//Usage returns the usage of an api
func (c App) Usage(token string) revel.Result {
	val, err := models.DB.Get(token).Result()
	if err != nil {
		//Error while getting the key hash of the rate request
		revel.AppLog.Error("Couldn't find the key in cache for", token)
		revel.AppLog.Error(err.Error())
		return c.RenderJSON(map[string]string{"Error": "Couldn't find the key in cache for " + token})
	}

	//parising the current api usage
	var curr models.APIUsage
	dec := json.NewDecoder(strings.NewReader(val))
	err = dec.Decode(&curr)
	if err != nil {
		//error while decoing the current usage
		revel.AppLog.Error("Error while decoing the current usage from redis for", token)
		revel.AppLog.Error(err.Error())
		return c.RenderJSON(map[string]string{"Error": "We are facing a technical difficulty while parsing your token usage of " + token})
	}

	return c.RenderJSON(curr)
}

//Reset resets the usage of an api token
func (c App) Reset(token string) revel.Result {
	models.DB.Publish(models.ResetChannelName, token)

	return c.RenderJSON(map[string]string{"Result": "Successfully reset the usage"})
}
