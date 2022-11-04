package set

// TODOL maybe drop

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(values ...T) {
	for _, value := range values {
		s[value] = struct{}{}
	}
}

func (s Set[T]) ToList() []T {
	values := make([]T, 0, len(s))
	for k := range s {
		values = append(values, k)
	}

	return values
}

func (s Set[T]) Remove(value T) {
	delete(s, value)
}

func (s Set[T]) Contains(value T) bool {
	_, c := s[value]

	return c
}

func NewSet[T comparable](values ...T) Set[T] {
	set := make(Set[T])
	set.Add(values...)

	return set
}
