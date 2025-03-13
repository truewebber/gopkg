package nilcheck

// Any returns true if at least one of the given arguments is nil.
func Any(args ...interface{}) bool {
	for _, arg := range args {
		if arg == nil {
			return true
		}
	}

	return false
}
