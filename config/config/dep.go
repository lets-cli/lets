package config

type Dep struct {
	Name string
	Args []string
	Env  map[string]string
}
