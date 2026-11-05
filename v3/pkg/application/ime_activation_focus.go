package application

// imeActivationFocusState is owned by the window thread. DefWindowProc handles
// an activating WM_ACTIVATE synchronously and can re-enter WndProc with
// WM_SETFOCUS. The depth limits suppression to that exact native call stack;
// later clicks and explicit focus requests must retain their immediate path.
type imeActivationFocusState struct {
	defaultActivationDepth uint
}

func (s *imeActivationFocusState) beginDefaultActivation() {
	s.defaultActivationDepth++
}

func (s *imeActivationFocusState) endDefaultActivation() {
	if s.defaultActivationDepth > 0 {
		s.defaultActivationDepth--
	}
}

func (s *imeActivationFocusState) suppressImmediateControllerFocus() bool {
	return s.defaultActivationDepth > 0
}
