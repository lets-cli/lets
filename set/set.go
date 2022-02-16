package set

type StringSet struct {
	entryMap map[string]struct{}
}

func (s *StringSet) Add(value string) {
	s.entryMap[value] = struct{}{}
}

func (s *StringSet) AddMany(values []string) {
	for _, value := range values {
		s.entryMap[value] = struct{}{}
	}
}

func (s *StringSet) ToList() []string {
	values := make([]string, 0, len(s.entryMap))
	for k := range s.entryMap {
		values = append(values, k)
	}

	return values
}

func (s *StringSet) Remove(value string) {
	delete(s.entryMap, value)
}

func (s *StringSet) Contains(value string) bool {
	_, c := s.entryMap[value]

	return c
}

func NewStringSet() *StringSet {
	return &StringSet{
		entryMap: make(map[string]struct{}),
	}
}

func NewStringSetWithValues(values []string) *StringSet {
	set := &StringSet{
		entryMap: make(map[string]struct{}),
	}
	set.AddMany(values)

	return set
}
