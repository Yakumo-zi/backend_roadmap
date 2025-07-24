package improve

func When[T any](cond bool, tv T, fv T) T {
	if cond {
		return tv
	}
	return fv
}
