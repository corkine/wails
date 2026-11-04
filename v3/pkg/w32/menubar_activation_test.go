//go:build windows

package w32

import "testing"

// Model the actual caller: if the menu layer says "not handled", the outer
// window procedure runs its activation logic, then calls DefWindowProc once.
// The old themed-menu path executed default handling twice in this sequence.
func TestMenuBarActivationDefaultRunsOnlyAtOuterWindow(t *testing.T) {
	originalDefault := menuBarDefaultProc
	t.Cleanup(func() { menuBarDefaultProc = originalDefault })
	for _, theme := range []struct {
		name  string
		value *MenuBarTheme
	}{{"unthemed", nil}, {"themed", &MenuBarTheme{}}} {
		for _, activation := range []WPARAM{WA_INACTIVE, WA_ACTIVE, WA_CLICKACTIVE} {
			t.Run(theme.name+"/"+string(rune('0'+activation)), func(t *testing.T) {
				calls := 0
				menuBarDefaultProc = func(HWND, uint32, WPARAM, LPARAM) LRESULT { calls++; return 0 }
				handled, _ := MenuBarWndProc(0, WM_ACTIVATE, activation, 0, theme.value)
				if handled {
					t.Fatal("menu layer must not swallow application activation handling")
				}
				if calls != 0 {
					t.Fatalf("menu layer activated host prematurely %d times", calls)
				}
				// Simulate the caller's final DefWindowProc; no real window is activated.
				menuBarDefaultProc(0, WM_ACTIVATE, activation, 0)
				if calls != 1 {
					t.Fatalf("default processing ran %d times, want once", calls)
				}
			})
		}
	}
}
