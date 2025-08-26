package generic

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// ResourceManager is a generic thread-safe manager for resources indexed by UUID
type ResourceManager[T any] struct {
	resources map[uuid.UUID]T
	mu        sync.RWMutex
}

// NewResourceManager creates a new generic resource manager
func NewResourceManager[T any]() *ResourceManager[T] {
	return &ResourceManager[T]{
		resources: make(map[uuid.UUID]T),
	}
}

// Get retrieves a resource by UUID
func (rm *ResourceManager[T]) Get(id uuid.UUID) (T, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	resource, exists := rm.resources[id]
	return resource, exists
}

// Set stores a resource with the given UUID
func (rm *ResourceManager[T]) Set(id uuid.UUID, resource T) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.resources[id] = resource
}

// Delete removes a resource by UUID
func (rm *ResourceManager[T]) Delete(id uuid.UUID) bool {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	_, exists := rm.resources[id]
	if exists {
		delete(rm.resources, id)
	}
	return exists
}

// Exists checks if a resource exists for the given UUID
func (rm *ResourceManager[T]) Exists(id uuid.UUID) bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	_, exists := rm.resources[id]
	return exists
}

// GetAll returns a copy of all resources
func (rm *ResourceManager[T]) GetAll() map[uuid.UUID]T {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	result := make(map[uuid.UUID]T, len(rm.resources))
	for k, v := range rm.resources {
		result[k] = v
	}
	return result
}

// Count returns the number of resources
func (rm *ResourceManager[T]) Count() int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return len(rm.resources)
}

// Clear removes all resources
func (rm *ResourceManager[T]) Clear() {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.resources = make(map[uuid.UUID]T)
}

// GetOrSet retrieves a resource or sets it if it doesn't exist
func (rm *ResourceManager[T]) GetOrSet(id uuid.UUID, factory func() T) T {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	if resource, exists := rm.resources[id]; exists {
		return resource
	}
	
	resource := factory()
	rm.resources[id] = resource
	return resource
}

// GetOrSetWithError retrieves a resource or sets it if it doesn't exist, with error handling
func (rm *ResourceManager[T]) GetOrSetWithError(id uuid.UUID, factory func() (T, error)) (T, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	if resource, exists := rm.resources[id]; exists {
		return resource, nil
	}
	
	resource, err := factory()
	if err != nil {
		var zero T
		return zero, err
	}
	
	rm.resources[id] = resource
	return resource, nil
}

// Update atomically updates a resource if it exists
func (rm *ResourceManager[T]) Update(id uuid.UUID, updater func(T) T) bool {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	resource, exists := rm.resources[id]
	if !exists {
		return false
	}
	
	rm.resources[id] = updater(resource)
	return true
}

// UpdateWithError atomically updates a resource if it exists, with error handling
func (rm *ResourceManager[T]) UpdateWithError(id uuid.UUID, updater func(T) (T, error)) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	resource, exists := rm.resources[id]
	if !exists {
		return fmt.Errorf("resource with id %s not found", id)
	}
	
	updatedResource, err := updater(resource)
	if err != nil {
		return err
	}
	
	rm.resources[id] = updatedResource
	return nil
}