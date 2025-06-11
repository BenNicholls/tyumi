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

// Init is run when the component is added to an entity. Use this to initialize any slices or maps or whatever.
func (c *Component) Init() {
	return
}

// Cleanup is run when the component is removed from an entity. Use this to... I dunno, send events?
func (c *Component) Cleanup() {
	return
}

// Register registers a type to be used as a component for entities. Types MUST be registered before being
// added to entities. Trying to add, get, or remove an unregistered component to/from an entity results in a panic.
func Register[T componentType]() {
	t := reflect.TypeFor[T]()
	if _, ok := typeMap[t]; ok { // duplicate register
		log.Debug("ECS: Duplicate component register! " + t.Name() + " already registered.")
		return
	}

	var newCache componentCache[T]
	componentCaches = append(componentCaches, &newCache)
	typeMap[t] = componentID(len(componentCaches) - 1)
}

// Add adds a new component of type T to an entity. The component type must be registered; if not, a panic
// occurs. Optionally, you can provide an already initialized component to be added. If the entity already has a
// component of this type, nothing is added and the initValue, if present, is ignored.
func Add[T componentType, ET ~uint32](entity ET, init_value ...T) {
	if !Alive(entity) {
		log.Error("ECS: Cannot add " + reflect.TypeFor[T]().Name() + " component to dead/invalid entity")
		return
	}

	getComponentCache[T]().addComponent(Entity(entity), init_value...)
}

// Get retrieves the component of type T from an entity. If the entity does not have the requested component,
// returns nil.
func Get[T componentType, ET ~uint32](entity ET) (component *T) {
	if !Alive(entity) {
		log.Error("Cannot get " + reflect.TypeFor[T]().Name() + " component from dead/invalid entity")
		return nil
	}

	return getComponentCache[T]().getComponent(Entity(entity))
}

// Has returns true if the entity contains the requested component.
func Has[T componentType, ET ~uint32](entity ET) bool {
	if !Alive(entity) {
		log.Error("Cannot query " + reflect.TypeFor[T]().Name() + " component of dead/invalid entity")
		return false
	}

	return getComponentCache[T]().hasComponent(Entity(entity))
}

// Remove removes the component of type T from the entity. If the entity does not have the requested component,
// does nothing.
func Remove[T componentType, ET ~uint32](entity ET) {
	if !Alive(entity) {
		log.Error("Cannot remove " + reflect.TypeFor[T]().Name() + " component from dead/invalid entity.")
	}

	getComponentCache[T]().removeComponent(Entity(entity))
}

// Toggle will add a component to an entity if it does not have one (optionally using the providing init value),
// otherwise it removes the component.
func Toggle[T componentType, ET ~uint32](entity ET, init ...T) {
	if !Alive(entity) {
		log.Error("Cannot toggle " + reflect.TypeFor[T]().Name() + " component from dead/invalid entity.")
	}

	if cache := getComponentCache[T](); cache.hasComponent(Entity(entity)) {
		cache.removeComponent(Entity(entity))
	} else {
		cache.addComponent(Entity(entity), init...)
	}
}
