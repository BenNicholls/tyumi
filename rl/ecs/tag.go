package ecs

import (
	"github.com/bennicholls/tyumi/log"
	"github.com/bennicholls/tyumi/util"
)

type Tag uint32

func (t Tag) isValid() bool {
	return int(t) < len(tagCaches)
}

var tagCaches []util.Set[Entity]

func init() {
	tagCaches = make([]util.Set[Entity], 0)
}

func RegisterTag() Tag {
	newTagSet := util.Set[Entity]{}
	tagCaches = append(tagCaches, newTagSet)
	return Tag(len(tagCaches) - 1)
}

func AddTag[ET ~uint32](entity ET, tag Tag) {
	if !tag.isValid() {
		log.Debug("ECS: Invalid Tag with ID: ", tag)

		return
	}

	tagCaches[int(tag)].Add(Entity(entity))
}

func HasTag[ET ~uint32](entity ET, tag Tag) bool {
	if !tag.isValid() {
		log.Debug("ECS: Invalid Tag with ID: ", tag)

		return false
	}

	return tagCaches[int(tag)].Contains(Entity(entity))
}

func RemoveTag[ET ~uint32](entity ET, tag Tag) {
	if !tag.isValid() {
		log.Debug("ECS: Invalid Tag with ID: ", tag)

		return
	}

	tagCaches[int(tag)].Remove(Entity(entity))
}
