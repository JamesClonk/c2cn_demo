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

func registerBackend() error {
	c, err := getRedisConnection()
	if err != nil {
		return err
	}
	defer c.Close()

	// add VCAP_APPLICATION.instance_id as key with a 15s TTL to redis for service discovery
	_, err = c.Do("SETEX", env.Application.InstanceID, "15", env.InstanceAddress)
	if err != nil {
		return err
	}

	// add VCAP_APPLICATION.instance_id to main set
	_, err = c.Do("SADD", redisBackendSet, env.Application.InstanceID)
	if err != nil {
		return err
	}

	return nil
}
