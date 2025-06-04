package ecs

import (
	"reflect"

	"github.com/bennicholls/tyumi/log"
)

type componentType interface {
	GetEntity() Entity
}

type Component struct {
	entity Entity
}

func (c Component) GetEntity() Entity {
	return c.entity
}

func (c *Component) setEntity(e Entity) {
	c.entity = e
}

// RegisterComponent registers a type to be used as a component for entities. Types MUST be registered before being
// added to entities. Trying to add an unregistered component to an entity results in a panic.
func RegisterComponent[T componentType]() {
	var newCache componentCache[T]
	componentCaches = append(componentCaches, &newCache)
	typeMap[reflect.TypeFor[T]()] = componentID(len(componentCaches) - 1)
}

// AddComponent adds a new component of type T to an entity. The component type must be registered; if not, a panic
// occurs. The added component is returned, just in case you want to immediately set some values. If the entity already
// has a component of this type, nothing is added and the already-there component is returned.
func AddComponent[T componentType](entity_id Entity) (component *T) {
	if !entity_id.Alive() {
		log.Error("Cannot add component to dead/invalid entity")
		return
	}

	if cache, ok := getComponentCache[T](); ok {
		component = cache.addComponent(entity_id)
		var i any = component
		if set, ok := i.(settableComponentType); ok {
			set.setEntity(entity_id)
		} else {
			panic("BAD!!!")
		}
		return
	} else {
		panic("Could not add component! (see log)")
	}
}

// GetComponent retrieves the component of type T from an entity. If component is unregistered of if the entity does
// not have the requested component, nil is returned and ok will be false.
func GetComponent[T componentType](entity_id Entity) (*T, bool) {
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
func RemoveComponent[T componentType](entity_id Entity) {
	if cache, ok := getComponentCache[T](); ok {
		cache.removeComponent(entity_id)
	} else {
		panic("Could not remove component! (see log)")
	}
}
