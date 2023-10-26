package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"github.com/redis/go-redis/v9"

	"wallemon/pkg/database"
)

type repoInfo struct {
	name       string
	repo       string
	tag        string
	env        []string
	ports      []int
	fixedPorts []int
	checkPort  func(host string, ports []string) error
}

var (
	containers = []*dockertest.Resource{}
	repoTags   = map[string]repoInfo{
		"postgres": {
			name: "postgreSQL",
			repo: "postgres",
			tag:  "13.7",
			env: []string{
				"POSTGRES_HOST_AUTH_METHOD=trust",
				"POSTGRES_DB=portto",
			},
			ports:     []int{5432},
			checkPort: checkPostgres,
		},
		"redis": {
			name:      "redis",
			repo:      "redis",
			tag:       "4.0.2",
			env:       []string{},
			ports:     []int{6379},
			checkPort: checkRedis,
		},
		"redis-cluster": {
			name: "redis",
			repo: "grokzen/redis-cluster",
			tag:  "6.2.8",
			env: []string{
				"IP=0.0.0.0",
			},
			ports: []int{7000, 7001, 7002, 7003, 7004, 7005},
			// make sure 7000-7005 are free to use
			fixedPorts: []int{7000, 7001, 7002, 7003, 7004, 7005},
			checkPort:  checkRedisCluster,
		},
	}
)

func Run(repo string) (host, port string, err error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return "", "", err
	}

	info, ok := repoTags[repo]
	if !ok {
		return "", "", fmt.Errorf("unsupported repo")
	}

	opt := dockertest.RunOptions{
		Repository: info.repo,
		Tag:        info.tag,
		Env:        info.env,
	}

	if len(info.fixedPorts) > 0 {
		// port bindings
		exposedPorts := []string{}
		portBindings := map[dc.Port][]dc.PortBinding{}
		for i, p := range info.fixedPorts {
			exposedPorts = append(exposedPorts, fmt.Sprintf("%d/tcp", info.ports[i]))
			portBindings[dc.Port(fmt.Sprintf("%d/tcp", info.ports[i]))] = []dc.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprintf("%d/tcp", p),
				},
			}
		}
		opt.ExposedPorts = exposedPorts
		opt.PortBindings = portBindings
	}

	// ciNetwork := config.GetCIContainerNetwork()
	// if ciNetwork != "" {
	// 	// Each build step of Cloud Build is run with its container attached to a local Docker network named `cloudbuild`.
	// 	// This allows build steps to communicate with each other and share data.
	// 	// For more information: https://cloud.google.com/build/docs/overview#build_configuration_and_build_steps
	// 	opt.NetworkID = ciNetwork
	// }

	res, err := pool.RunWithOptions(&opt)
	if err != nil {
		return "", "", err
	}

	// need more time to let redis-cluster setup, otherside it will fail due to CLUSTERDOWN error
	if repo == "redis-cluster" {
		fmt.Println("redis-cluster warming up...")
		time.Sleep(10 * time.Second)
	}

	var ports []string
	if err := pool.Retry(func() error {
		var tmpPorts []string
		host = "localhost"
		for _, p := range info.ports {
			var port string
			port = res.GetPort(fmt.Sprintf("%d/tcp", p))
			// if ciNetwork != "" {
			// 	// communicate with other containers by bridge network
			// 	host = res.Container.NetworkSettings.Networks[ciNetwork].IPAddress
			// 	port = fmt.Sprintf("%d", p)
			// }
			tmpPorts = append(tmpPorts, port)
		}

		if err := info.checkPort(host, tmpPorts); err != nil {
			fmt.Println("failed to check port: ", err)
			return err
		}

		containers = append(containers, res)
		ports = tmpPorts
		return nil
	}); err != nil {
		return "", "", err
	}

	return host, ports[0], nil
}

func Remove() error {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return err
	}

	for _, res := range containers {
		if err := pool.Purge(res); err != nil {
			return err
		}
	}
	containers = []*dockertest.Resource{}

	return nil
}

func checkPostgres(host string, ports []string) error {
	database.Initialize(context.Background())
	tmpDB := database.GetSQL()
	sqlDB, err := tmpDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func checkRedis(host string, ports []string) error {
	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, ports[0]),
	})
	if err := cli.Ping(context.Background()).Err(); err != nil {
		return err
	}
	return nil
}

func checkRedisCluster(host string, ports []string) error {
	var addresses []string
	for _, p := range ports {
		addresses = append(addresses, fmt.Sprintf("%s:%s", host, p))
	}
	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addresses,
	})
	if err := cli.Ping(context.Background()).Err(); err != nil {
		return err
	}
	return nil
}
