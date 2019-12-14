package flash

type Work struct {
	Executable
	done chan struct{}
}
