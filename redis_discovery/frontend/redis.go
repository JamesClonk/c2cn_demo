package main

import (
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func getRedisConnection() (redis.Conn, error) {
	service := env.GetService(redisServiceInstance)
	if service == nil {
		return nil, errors.New("Service is nil")
	}

	address := fmt.Sprintf("%v:%v", service.Credentials["hostname"], service.Credentials["port"])
	conn, err := redis.Dial("tcp", address)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Cannot connect to Redis[%v]: %v", address, err))
	}

	// login to redis with credentials
	if _, err := conn.Do("AUTH", fmt.Sprintf("%v", service.Credentials["password"])); err != nil {
		conn.Close()
		return nil, errors.New(fmt.Sprintf("Redis authentication error: %v", err))
	}

	return conn, nil
}

func increaseHitCounter() error {
	c, err := getRedisConnection()
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.Do("INCR", "hit-counter")
	if err != nil {
		return err
	}
	return nil
}

func getHitCounter() (int64, error) {
	c, err := getRedisConnection()
	if err != nil {
		return 0, err
	}
	defer c.Close()

	counter, err := redis.Int64(c.Do("GET", "hit-counter"))
	if err != nil {
		return 0, err
	}
	return counter, nil
}

func discoverBackends() ([]string, error) {
	c, err := getRedisConnection()
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// check if main set exists, if not return an empty list of backends
	exists, err := redis.Bool(c.Do("EXISTS", redisBackendSet))
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, nil
	}

	// get all enries from main set
	var backends []string
	entries, err := redis.Strings(c.Do("SMEMBERS", redisBackendSet))
	if err != nil {
		return backends, err
	}

	// go through all entries of main set
	for _, entry := range entries {
		exists, err := redis.Bool(c.Do("EXISTS", entry))
		if err != nil {
			return backends, err
		}

		if !exists {
			// if VCAP_APPLICATION.instance_id does not exists as key (because TTL expired)
			// then remove the entry from main set
			_, err := c.Do("SREM", redisBackendSet, entry)
			if err != nil {
				return backends, err
			}
		} else {
			// get value of VCAP_APPLICATION.instance_id key, and add to list of available backends
			backend, err := redis.String(c.Do("GET", entry))
			if err != nil {
				return backends, err
			}
			backends = append(backends, backend)
		}
	}
	return backends, nil
}
