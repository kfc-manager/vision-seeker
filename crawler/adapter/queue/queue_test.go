package queue

import (
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var q *queue

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
		Repository: "rabbitmq",
		Tag:        "4.0.4",
		Name:       "vision-seeker-queue-test",
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
		q, err = New(
			resource.Container.NetworkSettings.IPAddress,
			"5672",
			"test-queue",
			true,
			0,
		)
		return err
	}); err != nil {
		log.Fatalf("Could not connect to queue: %s", err.Error())
	}

	defer func() {
		err := q.Close()
		if err != nil {
			log.Printf("Could not close connection to queue: %s", err.Error())
		}
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err.Error())
		}
	}()

	// run tests
	m.Run()
}

func TestQueue(t *testing.T) {
	input := map[string]*struct{}{
		"eowfjoweifew": nil,
		"owejvweovwev": nil,
		"weoijfewoffl": nil,
		"woeijvewivvv": nil,
		"weocbgmlkadf": nil,
		"dsvkdvmdfskl": nil,
		"lsdfjwerowel": nil,
		"ovjojvsvdsvd": nil,
		"ljklfsdfsdkw": nil,
	}

	for k := range input {
		err := q.Push([]byte(k))
		if err != nil {
			t.Error(err.Error())
			return
		}
	}

	for i := 0; i < len(input); i++ {
		msg, err := q.Pull()
		if err != nil {
			t.Error(err.Error())
			return
		}
		input[string(msg)] = &struct{}{}
	}

	for k := range input {
		if input[k] == nil {
			t.Errorf("missing message: %s", k)
		}
	}
}
