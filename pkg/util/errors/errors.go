package errors

var (
	ErrCodeNone          = New(0, "")
	ParameterTransferErr = New(1000, "parameter transfer error")
)
