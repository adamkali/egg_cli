package templates_test

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/adamkali/egg_cli/pkg/templates"
)

const ResultSERVICES_RedisServiceTemplate = `
/* Generated by egg v0.0.1 */

package services

import (
	"context"
	"fmt"
	"time"

	"github.com/adamkali/egg/cmd/configuration"
	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	ctx    context.Context
	client *redis.Client
}

// Returns a refrence to a new UserService to be used in the controller
func CreateRedisService(ctx context.Context, config *configuration.Configuration) *RedisService {
	url := config.Cache.URL

	fmt.Printf("[INFO] RedisService.CreateRedisServiceurl{ url: %v }", url)
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	// connect to redis so that we can use it
	client := redis.NewClient(opts)
	return &RedisService{ctx, client}
}

// SetWithExpiration
//
// params:
//   key: string
//   value: string
//   expiration: time.Duration
// returns:
//   error
// 
// Sets a value in the redis cache with an expiration
// time
func (r *RedisService)SetWithExpiration(
	key string,
	value string,
	expiration time.Duration,
) error {
	err := r.client.Set(r.ctx, key, value, expiration).Err()
	return err
}

// Set
//
// params:
//   key: string
//   value: string
// returns:
//   error
//
// Sets a value in the redis cache. Uses SetWithExpiration with an expiration of 0
// seconds so that the value never expires
func (r *RedisService) Set(key string, value string) error {
	err := r.SetWithExpiration(key, value, 0)
	return err
}

// GetWithExpiration
//
// params:
//   key: string
//   expiration: time.Duration
// returns:
//   string
//   error
//
// Gets a value from the redis cache with an expiration
func (r *RedisService) GetWithExpiration(key string, expiration time.Duration) (string, error) {
	value, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// Get
//
// params:
//   key: string
// returns:
//   string
//   error
//
// Gets a value from the redis cache. Uses GetWithExpiration with an expiration of 0
// seconds.
func (r *RedisService) Get(key string) (string, error) {
	value, err := r.GetWithExpiration(key, 0)
	if err != nil {
		return "", err
	}
	return value, nil
}

// Delete
//
// params:
//   key: string
// returns:
//   error
//
// Deletes a value from the redis cache
func (r *RedisService) Delete(key string) error {
	err := r.client.Del(r.ctx, key).Err()
	return err
}

`

func TestSERVICES_RedisServiceTemplate(t *testing.T) {
	// load the template
	temp := templates.SERVICES_RedisServiceTemplate
	templateTest := template.Must(template.New("redis_service.go").Parse(temp))

	// execute the template
	stringWriter := new(bytes.Buffer)
	err := templateTest.ExecuteTemplate(stringWriter, "redis_service.go", createConfiguration())
	if err != nil {
		t.Error(err)
	}

	// check the result
	if stringWriter.String() != ResultSERVICES_RedisServiceTemplate {
		diff := Diff(stringWriter.String(), ResultSERVICES_RedisServiceTemplate)
		for i, v := range diff {
			t.Errorf("line %d: expected %s, got %s", i, v.Expected, v.Actual)
		}
	}
}
