package flash

// Executable - interface for executables
type Executable interface {
	Execute() error
	IsCompleted() bool
	OnSuccess()
	OnFailure(err error)
}
