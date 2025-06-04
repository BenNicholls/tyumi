package ecs

import (
	"testing"

	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

type testComponent struct {
	vec.Coord
}

func countAliveEntities() (count int) {
	for _, entity := range entities {
		if entity != INVALID_ID {
			count++
		}
	}

	return
}

func TestEntityAddRemove(t *testing.T) {
	totalRemoves := 0

	for range 1000 {
		CreateEntity()
	}

	if count := countAliveEntities(); count != 1000 {
		t.Errorf("Alive entities = %d, wanted 1000", count)
	}

	for range 400 {
		if entity := util.PickOne(entities); entity != INVALID_ID {
			RemoveEntity(entity)
			if entity.Alive() {
				t.Error("Entity remove failed")
			} else {
				totalRemoves++
			}
		}
	}

	if count := countAliveEntities(); count != 1000-totalRemoves {
		t.Errorf("Alive entities = %d, wanted %d", count, 1000-totalRemoves)
	}

	for range 1000 {
		CreateEntity()
	}

	if count := countAliveEntities(); count != 2000-totalRemoves {
		t.Errorf("Alive entities = %d, wanted %d", count, 2000-totalRemoves)
	}

	for range 400 {
		if entity := util.PickOne(entities); entity != INVALID_ID {
			RemoveEntity(entity)
			if entity.Alive() {
				t.Error("Entity remove failed")
			} else {
				totalRemoves++
			}
		}
	}

	if count := countAliveEntities(); count != 2000-totalRemoves {
		t.Errorf("Alive entities = %d, wanted %d", count, 2000-totalRemoves)
	}
}

func TestAddRemoveAlterComponents(t *testing.T) {
	RegisterComponent[testComponent]()
	var entities []Entity

	for range 30 {
		entities = append(entities, CreateEntity())
	}

	componentsAdded := 0

	for range 20 {
		entity := util.PickOne(entities)
		if _, ok := GetComponent[testComponent](entity); ok {
			continue
		}
		//testcomponent := entity.AddComponent(testComponentID).(*TestComponent)
		testcomponent := AddComponent[testComponent](entity)
		testcomponent.Coord = vec.Coord{int(entity), int(entity)}
		componentsAdded++
	}

	componentsFound := 0

	for _, entity := range entities {
		if test, ok := GetComponent[testComponent](entity); ok {
			componentsFound++
			if coord := (vec.Coord{int(entity), int(entity)}); test.Coord != coord {
				t.Errorf("Improperly set position for entity, position is %v, wanted %v", test.Coord, coord)
			}
		}
	}

	if componentsAdded != componentsFound {
		t.Errorf("Added %d components, but only found %d", componentsAdded, componentsFound)
	}

	componentsRemoved := 0

	for range 15 {
		entity := util.PickOne(entities)
		if _, ok := GetComponent[testComponent](entity); ok {
			RemoveComponent[testComponent](entity)
			componentsRemoved++
		}
	}

	componentsFound = 0

	for _, entity := range entities {
		if test, ok := GetComponent[testComponent](entity); ok {
			componentsFound++
			if coord := (vec.Coord{int(entity), int(entity)}); test.Coord != coord {
				t.Errorf("Improperly set position for entity, position is %v, wanted %v", test.Coord, coord)
			}
		}
	}

	if componentsFound != componentsAdded-componentsRemoved {
		t.Errorf("Found %d components, expected %d", componentsFound, componentsAdded-componentsRemoved)
	}
}
