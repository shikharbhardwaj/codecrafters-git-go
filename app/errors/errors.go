package commands

type PathError struct {
	Op   string
	Path string
	Err  error
}

func (e *PathError) Error() string {
	if e.Err != nil {
		return e.Op + " " + e.Path + ": " + e.Err.Error()
	} else {
		return e.Op + " " + e.Path
	}
}

type GitError struct {
	Message string
}

func (e GitError) Error() string {
	return e.Message
}
