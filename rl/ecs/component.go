package ecs

import (
	"reflect"

	"github.com/bennicholls/tyumi/log"
)

type componentType interface {
	GetEntity() Entity
}

// The base for all components! Embed this into your custom components.
type Component struct {
	entity Entity
}

// Returns the entity this component is attached to.
func (c Component) GetEntity() Entity {
	return c.entity
}

// Used by the content caches when adding a component to set the component's entity id. Should never be used for anything
// else!!
func (c *Component) setEntity(e Entity) {
	c.entity = e
}

// RegisterComponent registers a type to be used as a component for entities. Types MUST be registered before being
// added to entities. Trying to add, get, or remove an unregistered component to an entity results in a panic.
func RegisterComponent[T componentType]() {
	var newCache componentCache[T]
	componentCaches = append(componentCaches, &newCache)
	typeMap[reflect.TypeFor[T]()] = componentID(len(componentCaches) - 1)
}

// AddComponent adds a new component of type T to an entity. The component type must be registered; if not, a panic
// occurs. Optionally, you can provide an already initialized component to be added. If the entity already has a
// component of this type, nothing is added and the initValue, if present, is ignored.
func AddComponent[T componentType, EntityType ~uint32](entity EntityType, init_value ...T) {
	if !Alive(entity) {
		log.Error("Cannot add component to dead/invalid entity")
		return
	}

	getComponentCache[T]().addComponent(Entity(entity), init_value...)
}

// GetComponent retrieves the component of type T from an entity. If the entity does not have the requested component,
// returns nil.
func GetComponent[T componentType, EntityType ~uint32](entity EntityType) (component *T) {
	if !Alive(entity) {
		log.Error("Cannot get component from dead/invalid entity")
		return nil
	}

	return getComponentCache[T]().getComponent(Entity(entity))
}

// HasComponent returns true if the entity contains the requested component.
func HasComponent[T componentType, EntityType ~uint32](entity EntityType) bool {
	if !Alive(entity) {
		log.Error("Cannot query component of dead/invalid entity")
		return false
	}

	return getComponentCache[T]().hasComponent(Entity(entity))
}

// RemoveComponent removes the component of type T from the entity. If the entity does not have the requested component,
// does nothing.
func RemoveComponent[T componentType, EntityType ~uint32](entity EntityType) {
	if !Alive(entity) {
		log.Error("Cannot remove component from dead/invalid entity.")
	}

	getComponentCache[T]().removeComponent(Entity(entity))
}
