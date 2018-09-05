package configuration

type Environment int

const (
	Default     Environment = 0
	Development Environment = 1
	Testing     Environment = 2
	Production  Environment = 3
)
