package application

import "testing"

func TestIMEActivationFocusSuppressesOnlyDefaultActivationReentry(t *testing.T) {
	var state imeActivationFocusState
	if state.suppressImmediateControllerFocus() {
		t.Fatal("ordinary focus must remain immediate")
	}

	state.beginDefaultActivation()
	if !state.suppressImmediateControllerFocus() {
		t.Fatal("activation default reentry must wait for delayed recovery")
	}

	state.endDefaultActivation()
	if state.suppressImmediateControllerFocus() {
		t.Fatal("focus after activation default returns must remain immediate")
	}
}

func TestIMEActivationFocusBalancesNestedAndUnmatchedDefaults(t *testing.T) {
	var state imeActivationFocusState
	state.endDefaultActivation()
	if state.suppressImmediateControllerFocus() {
		t.Fatal("unmatched end must not underflow into suppression")
	}

	state.beginDefaultActivation()
	state.beginDefaultActivation()
	state.endDefaultActivation()
	if !state.suppressImmediateControllerFocus() {
		t.Fatal("nested activation default must remain suppressed until the outer return")
	}

	state.endDefaultActivation()
	if state.suppressImmediateControllerFocus() {
		t.Fatal("balanced nested defaults must restore ordinary focus")
	}
}
