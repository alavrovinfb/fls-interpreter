package messages

var (
	ErrBodyEmpty     = "%s function body is empty"
	ErrCast          = "%s can't cast cmd params"
	ErrCmdMissed     = "command %d in function %s has incorrect format 'cmd' key is missed"
	ErrInitMissed    = "init function is missed"
	ErrIncorrectType = "incorrect object %s type %T"
	ErrExecution     = "error during command executing %v"
	ErrFuncMissed    = "function %s is missed"
)
