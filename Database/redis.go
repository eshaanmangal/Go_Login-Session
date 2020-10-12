//Package Database Golbal Pkg
package Database

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// SaveKeyValue to save key in Redis Server
func SaveKeyValue(key, value string, duration time.Duration) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:15000",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	testping, err := rdb.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Redis Server Connected with %s & Error: %e \n", testping, err)
	}

	err = rdb.Set(key, value, duration*(time.Second)).Err()
	if err != nil {
		panic(err)
	}
}

// GetValue to save key in Redis Server
func GetValue(key string) (string,bool) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:15000",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	testping, err := rdb.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Redis Server Connected with %s & Error: %s", testping, err)
	}

	result, err := rdb.Get(key).Result()
	if err == redis.Nil {
		fmt.Printf("\n %s", "Key doesn't exist in DB")
		return "",false
	} else if err != nil {
		panic(err)
	}

	return result,true
}
