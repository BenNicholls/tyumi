package anim

import (
	"testing"
	"time"
)

// startAndPrime calls Start() and then does one zero-length Update. This is needed because Start() sets an internal
// "reset" flag, and the very next call to Update() only consumes that flag (resetting elapsed to 0) rather than
// advancing by the passed-in delta - so the delta passed to the first Update() after Start() is effectively
// dropped. This priming step gets the animation past that reset frame so subsequent Update() calls behave as
// expected.
func startAndPrime(a *Animation) {
	a.Start()
	a.Update(0)
}

func TestAnimationStartPlayPauseStop(t *testing.T) {
	var a Animation
	a.Duration = time.Second

	if a.IsPlaying() {
		t.Error("fresh animation should not be playing")
	}

	a.Start()
	if !a.IsPlaying() {
		t.Error("expected animation to be playing after Start()")
	}

	a.Pause()
	if a.IsPlaying() {
		t.Error("expected animation to not be playing after Pause()")
	}

	a.Play()
	if !a.IsPlaying() {
		t.Error("expected animation to be playing after Play()")
	}

	a.Stop()
	if a.IsPlaying() {
		t.Error("expected animation to not be playing after Stop()")
	}
}

func TestAnimationUpdateProgress(t *testing.T) {
	var a Animation
	a.Duration = 10 * time.Second
	startAndPrime(&a)

	a.Update(5 * time.Second)

	if got := a.GetProgress(); got != 0.5 {
		t.Errorf("GetProgress() = %v, want 0.5", got)
	}
	if a.IsDone() {
		t.Error("animation should not be done at 50% progress")
	}
}

func TestAnimationIsDone(t *testing.T) {
	var a Animation
	a.Duration = 10 * time.Second
	startAndPrime(&a)

	a.Update(10 * time.Second)

	if !a.IsDone() {
		t.Error("expected animation to be done once elapsed >= duration")
	}
}

func TestAnimationUpdateCallsOnDone(t *testing.T) {
	var a Animation
	a.Duration = time.Second
	called := false
	a.OnDone = func() { called = true }
	startAndPrime(&a)

	a.Update(time.Second)

	if !called {
		t.Error("expected OnDone callback to be called when animation finishes")
	}
	if a.IsPlaying() {
		t.Error("expected animation to stop playing once done")
	}
}

func TestAnimationRepeat(t *testing.T) {
	var a Animation
	a.Duration = 10 * time.Second
	a.Repeat = true
	startAndPrime(&a)

	a.Update(15 * time.Second) // should wrap around

	if a.IsDone() {
		t.Error("a repeating animation should never report IsDone")
	}
	if !a.JustLooped() {
		t.Error("expected JustLooped() to be true after wrapping past duration")
	}
}

func TestAnimationSetOneShotClearsRepeat(t *testing.T) {
	var a Animation
	a.Repeat = true

	a.SetOneShot(true)

	if !a.IsOneShot() {
		t.Error("expected IsOneShot() to be true")
	}
	if a.Repeat {
		t.Error("expected Repeat to be cleared when setting OneShot")
	}
}

func TestAnimationSetBlockingClearsRepeat(t *testing.T) {
	var a Animation
	a.Repeat = true

	a.SetBlocking(true)

	if !a.IsBlocking() {
		t.Error("expected IsBlocking() to be true")
	}
	if a.Repeat {
		t.Error("expected Repeat to be cleared when setting Blocking")
	}
}

func TestAnimationBlockingRepeatPreventedDuringUpdate(t *testing.T) {
	var a Animation
	a.Duration = 10 * time.Second
	a.Repeat = true
	a.Blocking = true
	startAndPrime(&a)

	a.Update(5 * time.Second) // Update() should defuse Repeat since Blocking is set

	if a.Repeat {
		t.Error("expected Repeat to be forced off during Update when Blocking is true")
	}
}

func TestAnimationGetTicksBackwards(t *testing.T) {
	var a Animation
	a.Duration = 10 * time.Second
	a.Backwards = true
	startAndPrime(&a)

	a.Update(3 * time.Second)

	got := a.GetTicks()
	want := a.Duration - 3*time.Second - 1
	if got != want {
		t.Errorf("GetTicks() backwards = %v, want %v", got, want)
	}
}

func TestAnimationIsUpdated(t *testing.T) {
	var a Animation
	if a.IsUpdated() {
		t.Error("fresh animation should not be marked as updated")
	}

	a.Updated = true
	if !a.IsUpdated() {
		t.Error("expected IsUpdated() to be true when Updated is set")
	}

	var alwaysUpdating Animation
	alwaysUpdating.AlwaysUpdates = true
	if !alwaysUpdating.IsUpdated() {
		t.Error("expected IsUpdated() to be true when AlwaysUpdates is set")
	}
}
