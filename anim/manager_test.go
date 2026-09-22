package anim

import (
	"testing"
	"time"
)

func TestAnimationManagerAddCountRemove(t *testing.T) {
	var am AnimationManager
	a1 := &Animation{Duration: time.Second}
	a2 := &Animation{Duration: time.Second}

	am.AddAnimation(a1)
	am.AddAnimation(a2)

	if got := am.CountAnimations(); got != 2 {
		t.Errorf("CountAnimations() = %d, want 2", got)
	}

	am.AddAnimation(a1) // duplicate, should be no-op
	if got := am.CountAnimations(); got != 2 {
		t.Errorf("CountAnimations() after duplicate add = %d, want 2", got)
	}

	am.RemoveAnimation(a1)
	if got := am.CountAnimations(); got != 1 {
		t.Errorf("CountAnimations() after remove = %d, want 1", got)
	}
}

func TestAnimationManagerUpdateAnimations(t *testing.T) {
	var am AnimationManager
	a := &Animation{Duration: 10 * time.Second}
	am.AddAnimation(a)
	a.Start()

	am.UpdateAnimations(0) // consume the reset frame (see Animation.Start/Update quirk)
	am.UpdateAnimations(5 * time.Second)

	if got := a.GetProgress(); got != 0.5 {
		t.Errorf("animation progress after manager update = %v, want 0.5", got)
	}
}

func TestAnimationManagerOnlyUpdatesPlayingAnimations(t *testing.T) {
	var am AnimationManager
	playing := &Animation{Duration: 10 * time.Second}
	notPlaying := &Animation{Duration: 10 * time.Second}
	am.AddAnimation(playing)
	am.AddAnimation(notPlaying)

	playing.Start()
	am.UpdateAnimations(0)
	am.UpdateAnimations(5 * time.Second)

	if playing.GetProgress() != 0.5 {
		t.Errorf("playing animation progress = %v, want 0.5", playing.GetProgress())
	}
	if notPlaying.GetProgress() != 0 {
		t.Errorf("non-playing animation progress = %v, want 0 (should not update)", notPlaying.GetProgress())
	}
}

func TestAnimationManagerOneShotRemovedWhenDone(t *testing.T) {
	var am AnimationManager
	a := &Animation{Duration: time.Second}
	am.AddOneShotAnimation(a)

	if got := am.CountAnimations(); got != 1 {
		t.Fatalf("CountAnimations() after AddOneShotAnimation = %d, want 1", got)
	}
	if !a.IsOneShot() {
		t.Error("expected AddOneShotAnimation to set OneShot")
	}
	if !a.IsPlaying() {
		t.Error("expected AddOneShotAnimation to start the animation")
	}

	am.UpdateAnimations(0) // consume reset frame
	am.UpdateAnimations(time.Second)

	if got := am.CountAnimations(); got != 0 {
		t.Errorf("CountAnimations() after one-shot animation finished = %d, want 0 (should be auto-removed)", got)
	}
}

func TestAnimationManagerHasBlockingAnimation(t *testing.T) {
	var am AnimationManager
	a := &Animation{Duration: time.Second, Blocking: true}
	am.AddAnimation(a)

	if am.HasBlockingAnimation() {
		t.Error("expected no blocking animation before it starts playing")
	}

	a.Start()
	if !am.HasBlockingAnimation() {
		t.Error("expected HasBlockingAnimation() to be true once the blocking animation is playing")
	}
}

func TestAnimationManagerEachAnimation(t *testing.T) {
	var am AnimationManager
	a1 := &Animation{Duration: time.Second}
	a2 := &Animation{Duration: time.Second}
	am.AddAnimation(a1)
	am.AddAnimation(a2)

	count := 0
	for range am.EachAnimation() {
		count++
	}
	if count != 2 {
		t.Errorf("EachAnimation() visited %d animations, want 2", count)
	}
}

func TestAnimationManagerEachPlayingAnimation(t *testing.T) {
	var am AnimationManager
	playing := &Animation{Duration: time.Second}
	notPlaying := &Animation{Duration: time.Second}
	am.AddAnimation(playing)
	am.AddAnimation(notPlaying)
	playing.Start()

	count := 0
	for a := range am.EachPlayingAnimation() {
		count++
		if a != Animator(playing) {
			t.Errorf("EachPlayingAnimation() yielded %v, want the playing animation", a)
		}
	}
	if count != 1 {
		t.Errorf("EachPlayingAnimation() visited %d animations, want 1", count)
	}
}
