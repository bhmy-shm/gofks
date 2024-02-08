package redisx

import (
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestAddToBlackList(t *testing.T) {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "10.35.149.23:30501",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Initialize BlackWhiteUser
	blackWhiteUser := NewBlackWhite()

	// Set up test data
	user := "testUser-black"

	// Add user to blacklist
	err := blackWhiteUser.AddToBlackList(user)
	assert.NoError(t, err)

	// Check if user is in the blacklist
	isBlacklisted, err := blackWhiteUser.CheckBlackList(user)
	assert.NoError(t, err)
	assert.True(t, isBlacklisted)

	result := blackWhiteUser.FindList(blacklist)
	log.Println("black list:", result)

	// Clean up test data
	blackWhiteUser.RemoveBlackList(user)
}

func TestAddToWhiteList(t *testing.T) {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Initialize BlackWhiteUser
	blackWhiteUser := NewBlackWhite()

	// Set up test data
	user := "testUser"

	// Add user to whitelist
	err := blackWhiteUser.AddToWhiteList(user)
	assert.NoError(t, err)

	// Check if user is in the whitelist
	isWhitelisted, err := blackWhiteUser.CheckWhiteList(user)
	assert.NoError(t, err)
	assert.True(t, isWhitelisted)

	// Clean up test data
	blackWhiteUser.RemoveWhiteList(user)
}

func TestRemoveBlackList(t *testing.T) {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Initialize BlackWhiteUser
	blackWhiteUser := NewBlackWhite()

	// Set up test data
	user := "testUser"
	blackWhiteUser.AddToBlackList(user)

	// Remove user from blacklist
	err := blackWhiteUser.RemoveBlackList(user)
	assert.NoError(t, err)

	// Check if user is in the blacklist
	isBlacklisted, err := blackWhiteUser.CheckBlackList(user)
	assert.NoError(t, err)
	assert.False(t, isBlacklisted)
}

func TestRemoveWhiteList(t *testing.T) {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Initialize BlackWhiteUser
	blackWhiteUser := NewBlackWhite()

	// Set up test data
	user := "testUser"
	blackWhiteUser.AddToWhiteList(user)

	// Remove user from whitelist
	err := blackWhiteUser.RemoveWhiteList(user)
	assert.NoError(t, err)

	// Check if user is in the whitelist
	isWhitelisted, err := blackWhiteUser.CheckWhiteList(user)
	assert.NoError(t, err)
	assert.False(t, isWhitelisted)
}

func TestCheckBlackList(t *testing.T) {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Initialize BlackWhiteUser
	blackWhiteUser := NewBlackWhite()

	// Set up test data
	user := "testUser"
	blackWhiteUser.AddToBlackList(user)

	// Check if user is in the blacklist
	isBlacklisted, err := blackWhiteUser.CheckBlackList(user)
	assert.NoError(t, err)
	assert.True(t, isBlacklisted)

	// Clean up test data
	blackWhiteUser.RemoveBlackList(user)
}
