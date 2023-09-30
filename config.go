package main

type Config struct {
	Round  Round  `ini:"round"`
	Remote Remote `ini:"remote"`
	Checks []Check
}

type Round struct {
	Title string `ini:"title"`
	Os    string `ini:"os"`
	User  string `ini:"user"`
	Local string `ini:"local"`
}

type Remote struct {
	Enable   bool `ini:"enable"`
	Name     bool `ini:"name"`
	Server   bool `ini:"server"`
	Password bool `ini:"password"`
}

type Check struct {
	Message string
	Points  int
	Cond    Condition
}

type PathExists struct {
	Path string
}

type FileContains struct {
	File  string
	Value string
}
