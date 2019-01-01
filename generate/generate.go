package generate

import (
	"fmt"
	"io"
	"sort"
	"sync"

	"github.com/shanna/phrasebook"
)

type Driver func() (Generator, error)

type Generator interface {
	Generate(phrasebook.Exports, io.Writer) error
}

var generatorsMutex sync.RWMutex
var generators = make(map[string]Driver)

func Register(name string, driver Driver) {
	generatorsMutex.Lock()
	defer generatorsMutex.Unlock()

	if driver == nil {
		panic("register generator driver is nil")
	}
	if _, ok := generators[name]; ok {
		panic(fmt.Sprintf("generator driver named '%s' already registered", name))
	}
	generators[name] = driver
}

// Generators registered.
func Generators() []string {
	generatorsMutex.RLock()
	defer generatorsMutex.RUnlock()
	var list []string
	for name := range generators {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

// New generator instance by name.
func New(name string) (Generator, error) {
	generatorsMutex.RLock()
	generator, ok := generators[name]
	generatorsMutex.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown generator %q (forgotten import?)", name)
	}
	return generator()
}
