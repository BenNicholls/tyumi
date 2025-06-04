package ecs

import (
	"reflect"

	"github.com/bennicholls/tyumi/log"
)

type Component any

type componentID uint32

var componentCaches []componentContainer
var typeMap map[reflect.Type]componentID

// this defines a component cache. component caches like to return the actual component type for add and get operations,
// so we can't put those functions in the interface here.
type componentContainer interface {
	removeComponent(id Entity)
}

func init() {
	componentCaches = make([]componentContainer, 0, 20)
	typeMap = make(map[reflect.Type]componentID)
}

type componentCache[T Component] struct {
	components []T
	indices    map[Entity]uint32 //index of entity index to component index in the cache. TODO: make this not a map someday.
}

func (cc componentCache[T]) getComponent(id Entity) (*T, bool) {
	if !id.Alive() {
		log.Error("Removed entity cannot get component!")
		return nil, false
	}

	if cc.indices == nil || cc.components == nil {
		return nil, false
	}

	if componentIdx, ok := cc.indices[id]; ok {
		return &cc.components[componentIdx], true
	} else {
		return nil, false
	}
}

// adds a component for the specified entity and returns a pointer to it so you can edit it. if the component already
// exists just returns a pointer to it
func (cc *componentCache[T]) addComponent(id Entity) *T {
	if comp, ok := cc.getComponent(id); ok {
		return comp
	}

	if cc.indices == nil {
		cc.indices = make(map[Entity]uint32)
		cc.components = make([]T, 0)
	}

	cc.indices[id] = uint32(len(cc.components))
	var newComponent T
	cc.components = append(cc.components, newComponent)
	return &cc.components[len(cc.components)-1]
}

func (cc *componentCache[T]) removeComponent(id Entity) {
	idx, ok := cc.indices[id]
	if !ok {
		return
	}

	endIndex := len(cc.components) - 1
	cc.components[idx] = cc.components[endIndex]
	delete(cc.indices, id)
	for k, v := range cc.indices {
		if v == uint32(endIndex) {
			cc.indices[k] = idx
			break
		}
	}
	var zero T
	cc.components[endIndex] = zero
	cc.components = cc.components[:endIndex]
}

func getComponentCache[T Component]() (*componentCache[T], bool) {
	componentType := reflect.TypeFor[T]()
	componentID, ok := typeMap[componentType]
	if !ok {
		log.Error(componentType.Name(), " is not a registered component type!!!")
		return nil, false
	}

	return componentCaches[componentID].(*componentCache[T]), true
}

// RegisterComponent registers a type to be used as a component for entities. Types MUST be registered before being
// added to entities. Trying to add an unregistered component to an entity results in a panic.
func RegisterComponent[T Component]() {
	var newCache componentCache[T]
	componentCaches = append(componentCaches, &newCache)
	typeMap[reflect.TypeFor[T]()] = componentID(len(componentCaches) - 1)
}

// AddComponent adds a new component of type T to an entity. The component type must be registered; if not, a panic
// occurs. The added component is returned, just in case you want to immediately set some values. If the entity already
// has a component of this type, nothing is added and the already-there component is returned.
func AddComponent[T Component](entity_id Entity) (component *T) {
	if !entity_id.Alive() {
		log.Error("Cannot add component to dead/invalid entity")
		return
	}

	if cache, ok := getComponentCache[T](); ok {
		return cache.addComponent(entity_id)
	} else {
		panic("Could not add component! (see log)")
	}
}

// GetComponent retrieves the component of type T from an entity. If component is unregistered of if the entity does
// not have the requested component, nil is returned and ok will be false.
func GetComponent[T Component](entity_id Entity) (*T, bool) {
	if !entity_id.Alive() {
		log.Error("Cannot get component from dead/invalid entity")
		return nil, false
	}

	if cache, ok := getComponentCache[T](); !ok {
		panic("Could not get component! (see log)")
	} else {
		if component, ok := cache.getComponent(entity_id); ok {
			return component, true
		} else {
			return nil, false
		}
	}
}

// RemoveComponents removes the component of type T from the entity. If the entity does not have the requested component,
// does nothing.
func RemoveComponent[T Component](entity_id Entity) {
	if cache, ok := getComponentCache[T](); ok {
		cache.removeComponent(entity_id)
	} else {
		panic("Could not remove component! (see log)")
	}
}
