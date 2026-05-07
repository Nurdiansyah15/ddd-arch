package apperror

// Kind adalah kategori error aplikasi — satu-satunya yang diketahui HTTP layer.
type Kind int

const (
	KindNotFound     Kind = iota // resource tidak ditemukan
	KindConflict                 // data sudah ada (email duplicate, dll)
	KindUnauthorized             // tidak punya akses / kredensial salah
	KindForbidden                // punya akses, tapi tidak boleh
	KindValidation               // input tidak valid
	KindInternal                 // unexpected / infra error — jangan expose Cause ke user
)

// AppError adalah error terstruktur yang dibawa dari usecase ke handler.
// Cause disimpan untuk keperluan logging, tidak pernah dikirim ke client.
type AppError struct {
	Kind    Kind
	Message string
	Cause   error
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Cause }

// ── Constructor helpers ────────────────────────────────────────────────────────

func NotFound(msg string, cause error) *AppError {
	return &AppError{Kind: KindNotFound, Message: msg, Cause: cause}
}

func Conflict(msg string) *AppError {
	return &AppError{Kind: KindConflict, Message: msg}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Kind: KindUnauthorized, Message: msg}
}

func Forbidden(msg string) *AppError {
	return &AppError{Kind: KindForbidden, Message: msg}
}

func Validation(msg string) *AppError {
	return &AppError{Kind: KindValidation, Message: msg}
}

func Internal(cause error) *AppError {
	return &AppError{Kind: KindInternal, Message: "internal server error", Cause: cause}
}
