package tyumi

// Timers allow you to run a function after a certain duration has passed. Call Process() on the timer to tick it down.
// Once the duration has lapsed the function will be run and it will be done; subsequent calls to Process() will do
// nothing. Use Done() to check for when the timer can be safely disposed of.
//
// You can use CreateTimer() on a scene object to make timers that are automatically managed and deleted.
type Timer struct {
	TimerFunction func()
	Ticks         int // ticks until timer ends
	done          bool
}

// Process ticks the timer. If the timer's duration has lapsed, does nothing.
func (t *Timer) Process() {
	if t.done {
		return
	}

	t.Ticks -= 1
	if t.Ticks <= 0 {
		if t.TimerFunction != nil {
			t.TimerFunction()
		}
		t.done = true
	}
}

// Done reports whether the timer's duration has lapsed.
func (t Timer) Done() bool {
	return t.done
}
