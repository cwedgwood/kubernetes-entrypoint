package entrypoint

import (
	"sync"
	"time"

	cli "github.com/stackanetes/kubernetes-entrypoint/client"
	"github.com/stackanetes/kubernetes-entrypoint/logger"
	"k8s.io/client-go/1.5/rest"
)

var dependencies []Resolver // List containing all dependencies to be resolved
const (
	//DependencyPrefix is a prefix for env variables
	DependencyPrefix      = "DEPENDENCY_"
	JsonSuffix            = "_JSON"
	resolverSleepInterval = 2
)

//Resolver is an interface which all dependencies should implement
type Resolver interface {
	IsResolved(entrypoint EntrypointInterface) (bool, error)
}

type EntrypointInterface interface {
	Resolve()
	Client() cli.ClientInterface
}

// Entrypoint is a main struct which checks dependencies
type Entrypoint struct {
	client    cli.ClientInterface
	namespace string
}

//Register is a function which registers new dependencies
func Register(res Resolver) {
	if res == nil {
		panic("Entrypoint: could not register nil Resolver")
	}
	dependencies = append(dependencies, res)
}

//New is a constructor for entrypoint
func New(config *rest.Config) (entry *Entrypoint, err error) {
	entry = new(Entrypoint)
	client, err := cli.New(config)
	if err != nil {
		return nil, err
	}
	entry.client = client
	return entry, err
}

func (e Entrypoint) Client() (client cli.ClientInterface) {
	return e.client
}

//Resolve is a main loop which iterates through all dependencies and resolves them
func (e Entrypoint) Resolve() {
	var wg sync.WaitGroup
	for _, dep := range dependencies {
		wg.Add(1)
		go func(dep Resolver) {
			defer wg.Done()
			logger.Info.Printf("Resolving %v", dep)
			var err error
			status := false
			for status == false {
				if status, err = dep.IsResolved(e); err != nil {
					logger.Warning.Printf("Resolving dependency %s failed: %v .", dep, err)
				}
				time.Sleep(resolverSleepInterval * time.Second)
			}
			logger.Info.Printf("Dependency %v is resolved.", dep)

		}(dep)
	}
	wg.Wait()

}
