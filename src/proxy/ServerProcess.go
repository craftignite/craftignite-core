package proxy

type ServerProcess struct {
	starting bool
	FileName string
}

func (process *ServerProcess) Start() {
	if process.starting {
		return
	}

	process.starting = true

	// TODO ...
}
