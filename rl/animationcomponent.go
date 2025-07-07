package rl

import (
	"slices"

	"github.com/bennicholls/tyumi/anim"
	"github.com/bennicholls/tyumi/gfx"
	"github.com/bennicholls/tyumi/rl/ecs"
)

func init() {
	ecs.Register[AnimationComponent]()
}

func AddAnimation[ET ~uint32](entity ET, a gfx.VisualAnimator) {
	animComp := ecs.Get[AnimationComponent](entity)
	if animComp == nil {
		ecs.Add[AnimationComponent](entity)
		animComp = ecs.Get[AnimationComponent](entity)
	}

	animComp.AddAnimation(a)
}

// AnimationComponent is a container for animations affecting an entity in the ECS.
type AnimationComponent struct {
	ecs.Component

	animations []anim.Animator
}

func (ac *AnimationComponent) AddAnimation(animation anim.Animator) {
	if ac.animations == nil {
		ac.animations = make([]anim.Animator, 0)
	}

	animation.SetOneShot(true)
	animation.Start()

	ac.animations = append(ac.animations, animation)
}

func (ac *AnimationComponent) ApplyVisualAnimations(vis gfx.Visuals) (result gfx.Visuals) {
	result = vis
	for _, animation := range ac.animations {
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

func (as *AnimationSystem) Update() {
	as.HasBlockingAnimation = false
	emptyAnimCompEntities := make([]ecs.Entity, 0)

	for animComp := range ecs.EachComponent[AnimationComponent]() {
		updated := false
		for _, animation := range animComp.animations {
			if animation.IsPlaying() {
				animation.Update()

				if !animation.IsDone() && animation.IsBlocking() {
					as.HasBlockingAnimation = true
				}

				if animation.IsUpdated() || animation.JustStopped() {
					updated = true
				}
			}
		}

		if updated {
			if pos := ecs.Get[PositionComponent](animComp.GetEntity()); pos != nil {
				as.tileMap.SetDirty(pos.Coord)
			}
		}

		animComp.animations = slices.DeleteFunc(animComp.animations, func(a anim.Animator) bool {
			return a.IsDone() && a.IsOneShot()
		})

		if len(animComp.animations) == 0 {
			emptyAnimCompEntities = append(emptyAnimCompEntities, animComp.GetEntity())
		}
	}

	for _, entity := range emptyAnimCompEntities {
		ecs.Remove[AnimationComponent](entity)
	}
}
