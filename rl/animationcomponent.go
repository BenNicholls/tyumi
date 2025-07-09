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

func AddAnimation[ET ~uint32](entity ET, a gfx.VisualAnimator, oneshot bool) {
	animComp := ecs.Get[AnimationComponent](entity)
	if animComp == nil {
		ecs.Add[AnimationComponent](entity)
		animComp = ecs.Get[AnimationComponent](entity)
	}

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

	for animComp := range ecs.EachComponent[AnimationComponent]() {
		animComp.UpdateAnimations(delta)
		if animComp.HasBlockingAnimation() {
			as.HasBlockingAnimation = true
		}

		if animComp.AnimationJustUpdated || animComp.AnimationJustStopped {
			if pos := ecs.Get[PositionComponent](animComp.GetEntity()); pos != nil {
				as.tileMap.SetDirty(pos.Coord)
			}
		}

		if animComp.CountAnimations() == 0 {
			emptyAnimCompEntities = append(emptyAnimCompEntities, animComp.GetEntity())
		}
	}

	for _, entity := range emptyAnimCompEntities {
		ecs.Remove[AnimationComponent](entity)
	}
}
