package cache

import (
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var c *cache

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err.Error())
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err.Error())
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "redis",
		Tag:        "7.4",
		Name:       "vision-seeker-redis-test",
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err.Error())
	}

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		c, err = New(
			resource.Container.NetworkSettings.IPAddress,
			"6379",
			"",
		)
		return err
	}); err != nil {
		log.Fatalf("Could not connect to cache: %s", err.Error())
	}

	defer func() {
		if err := c.Close(); err != nil {
			log.Printf("Could not close connection to cache: %s", err.Error())
		}
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err.Error())
		}
	}()

	// run tests
	m.Run()
}

func TestExist(t *testing.T) {
	input := []string{
		"fksdjf",
		"dlfkdsjf",
		"fewoifjwe",
		"eowfjweo",
		"ewofjwoefew",
		"ewofjwoefjwe",
	}

	for _, v := range input {
		err := c.Set(v)
		if err != nil {
			t.Error(err.Error())
			return
		}
	}

	for _, v := range input {
		exist, err := c.Exist(v)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if !exist {
			t.Errorf("cache miss key: %s", v)
		}
	}
}
