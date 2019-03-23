package controllers

import (
	"encoding/json"
	"strings"

	"github.com/go-redis/redis"
	"github.com/revel/revel"
	"github.com/shredx/golang-redis-rate-limiter/app/models"
)

/*
 * This file contains the logic required for the usage pubsub to work
 */

//InitUsageSubscriber inits the api usage subscriber function
func InitUsageSubscriber() {
	/*
	 * We will init the subcriber and store the usage channel
	 * Will init the subscriber and store the reset channel
	 * Will store the clock channel name
	 * Will init the usage go routine
	 * Will init the reste go routine
	 * Will init the rate process go routine
	 */
	//init he subscriber for usage channel
	sub := models.DB.Subscribe(revel.Config.StringDefault("redis.usagechannel", "apiusage"))
	sub.Receive()
	models.UsageChannel = sub.Channel()

	//initing the reste channel and storing it
	subR := models.DB.Subscribe(revel.Config.StringDefault("redis.resetchannel", "resetusage"))
	subR.Receive()
	models.ResetChannel = subR.Channel()

	//initing the clock channel name
	models.BlockChannel = revel.Config.StringDefault("redis.blockchannel", "block")

	//initing the api usage update go routine
	go Usage(models.UsageChannel, "API Usage updation request", models.ADD)

	//initing the api usage reste go routine
	go Usage(models.ResetChannel, "API Usage reset request", models.RESET)

	//rate process go routine to process the rate usage update and reset
	go RateProcess(models.RateChannel)
}

//Usage go routine will handle the api usage subscription messages
//desc is the usage description of the go routine
//t is the type of requests that needed to be catered
func Usage(ch <-chan *redis.Message, desc string, t models.Type) {
	/*
	 * We will go into an infinite for loop
	 * Will wait for usage request
	 * When usage information comes in we will create a request and pass it to the rate process channel
	 */
	for {
		m := <-ch
		revel.AppLog.Info(desc, m.Payload)
		models.RateChannel <- models.RateRequest{m.Payload, t}
	}
}

func RateProcess(ch chan models.RateRequest) {
	/*
	 * We will init a infinte for loop waiting for the requests to arrive
	 * We will get the current usage
	 * When requests arrive, depending upon the type of request we process.
	 * For ADD type requests we will update the current usage
	 * For RESET type we will reset the current usage
	 * At the end we will save the current usage
	 */
	for {
		req := <-ch
		val, err := models.DB.Get(req.KeyHash).Result()
		if err != nil {
			//Error while getting the key hash of the rate request
			revel.AppLog.Error("Couldn't find the key in cache for", req.KeyHash)
			revel.AppLog.Error(err.Error())
			continue
		}

		//parising the current api usage
		var curr models.ApiUsage
		dec := json.NewDecoder(strings.NewReader(val))
		err = dec.Decode(&curr)
		if err != nil {
			//error while decoing the current usage
			revel.AppLog.Error("Error while decoing the current usage from redis for", req.KeyHash)
			revel.AppLog.Error(err.Error())
			continue
		}

		//depending upon the type of request we process the request
		switch req.Type {
		case models.ADD:
			//for add tyope request we will get the current usage using the hash key
			curr.CurrentUsage++
			if curr.CurrentUsage >= curr.MaxUsage {
				go limitExceeded(req.KeyHash)
			}
		case models.RESET:
			//for reset type we will reste the current usage as 0
			curr.CurrentUsage = 0
		}

		//updating the api usage in the redis cache
		b := &strings.Builder{}
		enc := json.NewEncoder(b)
		err = enc.Encode(curr)
		if err != nil {
			//Error while encoding the curr
			revel.AppLog.Error("Error while encoding the current usage for", req.KeyHash)
			revel.AppLog.Error(err.Error())
			continue
		}
		err = models.DB.Set(req.KeyHash, b.String(), 0).Err()
		if err != nil {
			//error while writng the key hash
			revel.AppLog.Error("Error while updating  the api usage for", req.KeyHash)
			revel.AppLog.Error(err.Error())
			continue
		}
	}
}

//limitExceeded will publish the limit exceeded information through the redis message broker
func limitExceeded(keyhash string) {
	/*
	 * We will publish the blocked key to the block channel
	 */
	models.DB.Publish(models.BlockChannel, keyhash)
}
