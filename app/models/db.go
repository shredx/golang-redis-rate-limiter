package models

import "github.com/go-redis/redis"

/*
 * This file contains the resources required for making database connection
 */

//DB is the db client for the redis persistent storage
var DB *redis.Client
