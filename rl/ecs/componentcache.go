package ecs

import (
	"reflect"

	"github.com/bennicholls/tyumi/log"
)

type settableComponentType interface {
	componentType
	setEntity(e Entity)
	Init()
	Cleanup()
}

type componentID uint32

var componentCaches []componentContainer
var typeMap map[reflect.Type]componentID

// this defines a component cache. component caches like to return the actual component type for add and get operations,
// so we can't put those functions in the interface here.
type componentContainer interface {
	copyComponent(id Entity, new_id Entity)
	hasComponent(id Entity) bool
	removeComponent(id Entity)
}

func init() {
	componentCaches = make([]componentContainer, 0, 20)
	typeMap = make(map[reflect.Type]componentID)
}

type componentCache[T componentType] struct {
	components []T
	indices    map[Entity]uint32 //index of entity index to component index in the cache. TODO: make this not a map someday.
}

func (cc componentCache[T]) getComponent(entity Entity) *T {
	if cc.indices == nil {
		return nil
	}

	if componentIdx, ok := cc.indices[entity]; ok {
		return &cc.components[componentIdx]
	} else {
		return nil
	}
}

// adds a component for the specified entity. optionally allows you to provide an initial value for the newly created
// component
func (cc *componentCache[T]) addComponent(entity Entity, init ...T) {
	if cc.hasComponent(entity) {
		return
	}

	if cc.indices == nil {
		cc.indices = make(map[Entity]uint32)
		cc.components = make([]T, 0)
	}

	cc.indices[entity] = uint32(len(cc.components))

	var newComponent T
	if len(init) > 0 {
		newComponent = init[0]
	}

	var i any = &newComponent
	if set, ok := i.(settableComponentType); ok {
		set.setEntity(entity)
		set.Init()
	} else {
		panic("COULD NOT SET ENTITY ID ON ADDED COMPONENT??? BAD!!!")
	}

	cc.components = append(cc.components, newComponent)
}

// creates a copy of entity's component, assigned to copy
func (cc *componentCache[T]) copyComponent(entity, copy Entity) {
	cc.addComponent(copy, cc.components[cc.indices[entity]])
}

func (cc *componentCache[T]) hasComponent(entity Entity) bool {
	_, ok := cc.indices[entity]
	return ok
}

func (cc *componentCache[T]) removeComponent(entity Entity) {
	idx, ok := cc.indices[entity]
	if !ok {
		return
	}

	// covertly convert component to the settable form and run a cleanup function if it is defined.
	var i any = &cc.components[idx]
	i.(settableComponentType).Cleanup()

	delete(cc.indices, entity) // delete saved index for removed component
	endIndex := len(cc.components) - 1
	if idx != uint32(endIndex) { // if removed entity is NOT the final entity in the component:
		cc.components[idx] = cc.components[endIndex]     // overwrite removed component with component on the end
		cc.indices[cc.components[idx].GetEntity()] = idx // update index for component that was moved
	}

	var zero T
	cc.components[endIndex] = zero           // replace moved component data with zero value (just in case it was holding a pointer or something)
	cc.components = cc.components[:endIndex] // reslice component list to new len
}

func getComponentCache[T componentType]() *componentCache[T] {
	componentType := reflect.TypeFor[T]()
	componentID, ok := typeMap[componentType]
	if !ok {
		log.Error(componentType.Name(), " is not a registered component type!!!")
		panic("Unregistered component detected! (see log)")
	}

	return componentCaches[componentID].(*componentCache[T])
}
