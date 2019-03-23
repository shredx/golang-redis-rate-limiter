package models

import "github.com/go-redis/redis"

/*
 * This file contains the models defintions for api usage
 */

//UsageChannel recieves the redis api usage channel
var UsageChannel <-chan *redis.Message

//ResetChannel recieves the redis api usage reset channel
var ResetChannel <-chan *redis.Message

//RateChannel recieves the rate updation requests
var RateChannel chan RateRequest = make(chan RateRequest)

//BlockChannel is the name of the channel that should be used for blocking the api subscription key
var BlockChannel string

//Apiusage represents the api usage of the key
type ApiUsage struct {
	Key          string `json:"key"`     //Key is the key associated with the api usage
	MaxUsage     int    `json:"usage"`   //MaxUsage is the maximum allowed usage
	CurrentUsage int    `json:"current"` //CurrentUsage is the current usage
	Email        string `json:"email"`   //Email is the email associated with the key
}

//Type is the type of the request for rate limitng go routine
type Type uint

const (
	//ADD will add the current usage of an api with 1
	ADD Type = iota
	//RESET will reset the the current api usage to 0
	RESET
)

//RateRequest is the request to update the rate
type RateRequest struct {
	KeyHash string //KeyHash is used to store the rate info in redis cache
	Type    Type   //Type of the rate update request
}
