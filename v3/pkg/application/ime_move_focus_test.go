package application

import "testing"

func TestIMEMoveFocusTimerOrdering(t *testing.T) {
	for _, tc := range []struct {
		name   string
		during bool
		after  bool
	}{
		{"activation-before-drag/timer-during", false, false},
		{"activation-before-drag/timer-after", false, true},
		{"activation-during-drag/timer-during", true, false},
		{"activation-during-drag/timer-after", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var state imeMoveFocusState
			token := state.epoch
			state.begin()
			if tc.during {
				token = state.epoch
			}
			if !state.moving {
				t.Fatal("controller focus must be suppressed during drag")
			}
			if tc.after && !state.end() {
				t.Fatal("missing move-end restoration")
			}
			if state.acceptsTimer(token) {
				t.Fatal("timer must not focus during drag or after move-end restoration")
			}
			if !tc.after && !state.end() {
				t.Fatal("missing move-end restoration")
			}
			if state.end() {
				t.Fatal("duplicate exit must not restore twice")
			}
			if state.moving || !state.acceptsTimer(state.epoch) {
				t.Fatal("ordinary activation after drag must still work")
			}
		})
	}
}

func TestIMEMoveFocusOrdinaryActivationUnchanged(t *testing.T) {
	var state imeMoveFocusState
	token := state.epoch
	if !state.acceptsTimer(token) || !state.acceptsTimer(token) {
		t.Fatal("ordinary activation compensation was suppressed")
	}
	if state.end() {
		t.Fatal("unmatched exit must not request focus")
	}
	state.begin()
	state.end()
	newToken := state.epoch
	state.begin()
	state.end()
	if state.acceptsTimer(token) || state.acceptsTimer(newToken) {
		t.Fatal("a later drag must invalidate earlier timers")
	}
}
