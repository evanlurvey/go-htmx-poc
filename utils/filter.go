package utils

func Filter[S ~[]E, E comparable](s S, f E) S {
	var out S
	for i := range s {
		if s[i] == f {
			out = append(out, s[i])
		}
	}
	return out
}

func FilterFunc[S ~[]E, E any](s S, f func(E) bool) S {
	var out S
	for i := range s {
		if f(s[i]) {
			out = append(out, s[i])
		}
	}
	return out
}
