package rl

import (
	"time"

	"github.com/bennicholls/tyumi/anim"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/rl/ecs"
)

func init() {
	ecs.Register[AnimationComponent]()
}

func AddAnimation[ET ~uint32](entity ET, a anim.Animator, oneshot bool) {
	animComp := ecs.GetOrAdd[AnimationComponent](entity)

	if oneshot {
		animComp.AddOneShotAnimation(a)
	} else {
		animComp.AddAnimation(a)
	}
}

// AnimationComponent is a container for animations affecting an entity in the ECS.
type AnimationComponent struct {
	ecs.Component
	anim.AnimationManager
}

func (ac *AnimationComponent) ApplyVisualAnimations(vis gfx.Visuals) (result gfx.Visuals) {
	result = vis
	for animation := range ac.EachPlayingAnimation() {
		if visualAnimation, ok := animation.(gfx.VisualAnimator); ok {
			result = visualAnimation.ApplyToVisuals(result)
		}
	}

	return
}

type AnimationSystem struct {
	System

	HasBlockingAnimation bool // a blocking animation is playing

	tileMap *TileMap
}

func (as *AnimationSystem) Init(tm *TileMap) {
	as.tileMap = tm
}

func (as *AnimationSystem) Update(delta time.Duration) {
	as.HasBlockingAnimation = false
	emptyAnimCompEntities := make([]ecs.Entity, 0)

	for animComp, entity := range ecs.EachComponent[AnimationComponent]() {
		animComp.UpdateAnimations(delta)
		if animComp.HasBlockingAnimation() {
			as.HasBlockingAnimation = true
		}

		if animComp.AnimationJustUpdated || animComp.AnimationJustStopped {
			// If something has updated or stopped, we check all of the animations to see what parts of the tilemap
			// to set as dirty. For now we're just doing VisualAnimators, but once support for CanvasAnimators goes in
			// we'll have to think this through my thoroughly. Ideally this loop shouldn't be necessary at all, we
			// should have a way of directly extracting the updated/stopped animations from the animation manager on
			// the component. TODO.
			for animation := range animComp.EachAnimation() {
				if _, ok := animation.(gfx.VisualAnimator); ok {
					if pos := ecs.Get[PositionComponent](entity); pos != nil {
						as.tileMap.SetDirty(pos.Coord)
					}
					break
				}
			}
		}

		if animComp.CountAnimations() == 0 {
			emptyAnimCompEntities = append(emptyAnimCompEntities, entity)
		}
	}

	for _, entity := range emptyAnimCompEntities {
		if !ecs.Alive(entity) {
			continue
		}

		ecs.Remove[AnimationComponent](entity)
	}
}
