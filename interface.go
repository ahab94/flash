package xflow

// Executable - interface for executables
type Executable interface {
	Execute() error
	IsCompleted() bool
}
