package ecs

import (
	"testing"

	"github.com/bennicholls/tyumi/util"
	"github.com/bennicholls/tyumi/vec"
)

type testComponent struct {
	Component

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
	batchAdds := 1000
	totalRemoves := 0

	for range batchAdds {
		CreateEntity()
	}

	if count := countAliveEntities(); count != batchAdds {
		t.Errorf("Alive entities = %d, wanted 1000", count)
	}

	for range batchAdds / 2 {
		if entity := util.PickOne(entities); entity != INVALID_ID {
			RemoveEntity(entity)
			if entity.Alive() {
				t.Error("Entity remove failed")
			} else {
				totalRemoves++
			}
		}
	}

	if count := countAliveEntities(); count != batchAdds-totalRemoves {
		t.Errorf("Alive entities = %d, wanted %d", count, batchAdds-totalRemoves)
	}

	for range batchAdds {
		CreateEntity()
	}

	if count := countAliveEntities(); count != 2*batchAdds-totalRemoves {
		t.Errorf("Alive entities = %d, wanted %d", count, 2*batchAdds-totalRemoves)
	}

	for range batchAdds / 2 {
		if entity := util.PickOne(entities); entity != INVALID_ID {
			RemoveEntity(entity)
			if entity.Alive() {
				t.Error("Entity remove failed")
			} else {
				totalRemoves++
			}
		}
	}

	if count := countAliveEntities(); count != 2*batchAdds-totalRemoves {
		t.Errorf("Alive entities = %d, wanted %d", count, 2*batchAdds-totalRemoves)
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
		if HasComponent[testComponent](entity) {
			continue
		}

		AddComponent(entity, testComponent{Coord: vec.Coord{int(entity), int(entity)}})
		componentsAdded++
	}

	componentsFound := 0

	for _, entity := range entities {
		if test := GetComponent[testComponent](entity); test != nil {
			componentsFound++
			if coord := (vec.Coord{int(entity), int(entity)}); test.Coord != coord {
				t.Errorf("Improperly set position for entity, position is %v, wanted %v", test.Coord, coord)
			}

			if test.entity != entity {
				t.Errorf("Improperly set entityID component. Found %v, wanted %v", test.entity, entity)
			}
		}
	}

	if componentsAdded != componentsFound {
		t.Errorf("Added %d components, but only found %d", componentsAdded, componentsFound)
	}

	componentsRemoved := 0

	for range 15 {
		entity := util.PickOne(entities)
		if HasComponent[testComponent](entity) {
			RemoveComponent[testComponent](entity)
			componentsRemoved++
		}
	}

	componentsFound = 0

	for _, entity := range entities {
		if test := GetComponent[testComponent](entity); test != nil {
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
