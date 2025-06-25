package rl

import (
	"github.com/bennicholls/tyumi/event"
	"github.com/bennicholls/tyumi/rl/ecs"
)

func init() {
	ecs.Register[HealthComponent]()
}

type HealthComponent struct {
	ecs.Component

	HP Stat[int]
}

func (hc *HealthComponent) ChangeHealth(delta int) {
	if delta == 0 {
		return
	}

	oldHealth := hc.HP.Get()
	hc.HP.Mod(delta)

	if hc.HP.Get() != oldHealth {
		e := Entity(hc.GetEntity())
		event.Fire(EV_ENTITYHEALTHCHANGED, &EntityHealthChangedEvent{
			Entity: e, OldHP: oldHealth, NewHP: hc.HP.Get()})
		if hc.HP.Get() == 0 {
			event.Fire(EV_ENTITYDIED, &EntityEvent{Entity: e})
		}
	}
}
