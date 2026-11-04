package application

// imeMoveFocusState is owned by the window thread. Epochs only change across
// move/size loops, so ordinary activation compensation retains its behavior.
type imeMoveFocusState struct {
	moving bool
	epoch  uintptr
}

func (s *imeMoveFocusState) begin() {
	s.moving = true
	s.epoch++
}

func (s *imeMoveFocusState) end() bool {
	if !s.moving {
		return false
	}
	s.moving = false
	// Invalidate queued timers even when a short drag finishes before they fire.
	s.epoch++
	return true
}

func (s *imeMoveFocusState) acceptsTimer(epoch uintptr) bool {
	return !s.moving && epoch == s.epoch
}
