package discovery

import (
	"sync"

	"github.com/Axway/agent-sdk/pkg/cache"
	"github.com/Axway/agents-layer7/pkg/models/service"
)

type validator struct {
	cache cache.Cache
	lock  sync.RWMutex
}

func newValidator() *validator {
	return &validator{
		cache: cache.New(),
	}
}

// SetAPIs takes a list of apis and sets the ids to the cache. Removes all old api ids before updating.
func (v *validator) SetAPIs(svcs []service.Item) {
	v.lock.Lock()
	defer v.lock.Unlock()

	oldKeys := v.cache.GetKeys()
	for _, k := range oldKeys {
		v.cache.Delete(k)
	}

	for _, svc := range svcs {
		v.cache.Set(svc.ID, svc)
	}
}

// Validate returns true if the api is found in the cache, and false if not
func (v *validator) Validate(apiID, _ string) bool {
	v.lock.RLock()
	defer v.lock.RUnlock()

	_, err := v.cache.Get(apiID)
	return err == nil
}
