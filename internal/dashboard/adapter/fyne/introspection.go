package fyne

// RenderedObjectCount returns the number of live canvas image objects currently
// owned by the adapter. Tests use this to verify object reuse and part counts
// without reaching into the container internals.
func (a *Adapter) RenderedObjectCount() int {
	if a == nil || a.root == nil {
		return 0
	}
	return len(a.root.Objects)
}
