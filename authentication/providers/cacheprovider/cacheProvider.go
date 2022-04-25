package cacheprovider

import (
	"Assignment/models"
	"Assignment/providers"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type cacheProvider struct {
	cache map[string]interface{}
	mutex *sync.Mutex
}

// NewCacheProvider - this is the function to initialise the cache provider
// this method is taking a map as argument which is the in-memory storage
// To handle the multiple instances of server we can have a service like redis or dynamo db which can be plugged in with this provider
func NewCacheProvider(Cache map[string]interface{}) providers.CacheProvider {
	return &cacheProvider{
		cache: Cache,
		mutex: &sync.Mutex{},
	}
}

// Get - this function is getting the value from the cache based on the key
func (r cacheProvider) Get(key string) (string, error) {
	var value string
	val, ok := r.cache[key]
	if ok {
		var tokenDetails models.TokenDetail
		value = string(val.([]byte))
		err := json.Unmarshal(val.([]byte), &tokenDetails)
		if err != nil {
			logrus.Errorf("Set: cache err %v", errors.New("key not exist"))
			return "", errors.New("key not exist")
		}

		if tokenDetails.ExpirationTime.Before(time.Now()) {
			logrus.Errorf("Set: cache err %v", errors.New("key not exist"))
			return "", errors.New("key not exist")
		}

		return value, nil
	}

	return value, errors.New("key not exist")
}

// Set - this function is setting the value from the cache based on the key
func (r cacheProvider) Set(key string, value interface{}) error {
	r.mutex.Lock()

	byteData, err := json.Marshal(value)
	if err != nil {
		logrus.Errorf("Set: cache err %v", err)
		return err
	}
	r.cache[key] = byteData
	r.mutex.Unlock()
	_, ok := r.cache[key]
	if ok {
		return nil
	}

	return errors.New("not able to set")
}

// Delete - this function is deleting the key from the cache
func (r cacheProvider) Delete(key string) int64 {
	delete(r.cache, key)
	_, ok := r.cache[key]
	if ok {
		return 1
	}

	return 0
}
