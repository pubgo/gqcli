package errors

// ErrTag err tags
var ErrTag = struct {
	UnknownTypeCode string
}{
	"errors.unknown_type",
}

func init() {
	ErrTagRegistry(ErrTag)
}
