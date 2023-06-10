package poleerror

type PoleError struct {
	err    error
	ErrFuc ErrorFuc
}

func Default() *PoleError {
	return &PoleError{}
}
func (e *PoleError) Error() string {
	return e.err.Error()
}

func (e *PoleError) Put(err error) {
	e.check(err)
}

func (e *PoleError) check(err error) {
	if err != nil {
		e.err = err
		panic(e)
	}
}

type ErrorFuc func(msError *PoleError)

//暴露一个方法 让用户自定义

func (e *PoleError) Result(errFuc ErrorFuc) {
	e.ErrFuc = errFuc
}
func (e *PoleError) ExecResult() {
	e.ErrFuc(e)
}
