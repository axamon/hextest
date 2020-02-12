package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

func redisConnection(url string) *redis.Client {
	fmt.Println("Connecting to Redis DB")
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := client.Ping().Err()

	if err != nil {
		panic(err)
	}
	return client
}

func postgresConnection(database string) *sql.DB {
	fmt.Println("Connecting to PostgreSQL DB")
	db, err := sql.Open("postgres", database)
	if err != nil {
		log.Fatalf("%s", err)
		panic(err)
	}
	return db
}
