package database

import (
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var db *database

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

	// Build a custom image from the Dockerfile
	err = pool.Client.BuildImage(docker.BuildImageOptions{
		Name:         "vision-seeker-database", // Name for the custom image
		ContextDir:   "../../../database",      // Path to the Dockerfile directory
		Dockerfile:   "Dockerfile",             // Name of the Dockerfile
		OutputStream: log.Writer(),             // Log output for the build process
	})
	if err != nil {
		log.Fatalf("Could not build Docker image: %s", err.Error())
	}

	// Run the container using the custom image
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "vision-seeker-database",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_PASSWORD=password123",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=postgres",
		},
		Name: "vision-seeker-postgres-test",
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start container: %s", err.Error())
	}

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = New(
			resource.Container.NetworkSettings.IPAddress,
			"5432",
			"postgres",
			"postgres",
			"password123",
		)
		return err
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err.Error())
	}

	defer func() {
		db.Close()
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err.Error())
		}
	}()

	// run tests
	m.Run()
}

func TestExistUrl(t *testing.T) {
	input := []string{
		"dflkjsdofz",
		"fowifjewfw",
		"goahgphwge",
		"gowejgwoeg",
		"weojfowefu",
		"oewgiwggwe",
		"woegoewgwh",
		"owegtowele",
	}

	for _, v := range input {
		ok, err := db.InsertUrl(v)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if !ok {
			t.Errorf("insert not ok for key: %s", v)
			return
		}
	}

	for _, v := range input {
		exist, err := db.ExistUrl(v)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if !exist {
			t.Errorf("missing key: %s", v)
		}
	}
}
