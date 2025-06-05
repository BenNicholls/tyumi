package ecs

import (
	"reflect"

	"github.com/bennicholls/tyumi/log"
)

type settableComponentType interface {
	componentType
	setEntity(e Entity)
}

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

type componentCache[T componentType] struct {
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

	delete(cc.indices, id) // delete saved index for removed component
	endIndex := len(cc.components) - 1
	if idx != uint32(endIndex) { // if removed entity is NOT the final entity in the component:
		cc.components[idx] = cc.components[endIndex]     // overwrite removed component with component on the end
		cc.indices[cc.components[idx].GetEntity()] = idx // update index for component that was moved
	}

	var zero T
	cc.components[endIndex] = zero           // replace moved component data with zero value (just in case it was holding a pointer or something)
	cc.components = cc.components[:endIndex] // reslice component list to new len
}

func getComponentCache[T componentType]() (*componentCache[T], bool) {
	componentType := reflect.TypeFor[T]()
	componentID, ok := typeMap[componentType]
	if !ok {
		log.Error(componentType.Name(), " is not a registered component type!!!")
		return nil, false
	}

	return componentCaches[componentID].(*componentCache[T]), true
}
