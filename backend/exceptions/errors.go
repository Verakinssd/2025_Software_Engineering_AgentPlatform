package exceptions

import (
	"encoding/json"
	"errors"
)

type Error struct {
	Code   int
	Err    error
	Detail string
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Detail  string `json:"detail"`
	}{
		Code:    e.Code,
		Message: e.Error(),
		Detail:  e.Detail,
	})
}

func (e *Error) WithDetail(detail string) *Error {
	e1 := &Error{
		Code:   e.Code,
		Err:    e.Err,
		Detail: detail,
	}
	e1.Err = errors.New(e.Err.Error() + " " + detail)
	e1.Detail = detail
	return e1
}

func InternalServerError(err error) *Error {
	return &Error{Code: 500, Err: err}
}

func InvalidRequest(err error) *Error {
	return &Error{Code: 400, Err: err}
}

var (
	ErrTokenMissing       = &Error{Code: 401, Err: errors.New("user not logged in")}
	ErrTokenInvalid       = &Error{Code: 402, Err: errors.New("token invalid")}
	ErrTokenClaimsInvalid = &Error{Code: 403, Err: errors.New("token claims invalid")}
	ErrForbidden          = &Error{Code: 403, Err: errors.New("operation forbidden")}
	ErrTokenGenerate      = &Error{Code: 404, Err: errors.New("token generate failed")}
	ErrPenaltyNotFound    = &Error{Code: 404, Err: errors.New("penalty not found")}
	ErrCommentNotFound    = &Error{Code: 404, Err: errors.New("comment not found")}
	ErrPermissionDenied   = &Error{Code: 405, Err: errors.New("permission denied")}
	ErrInvalidFileTypes   = &Error{Code: 406, Err: errors.New("invalid file types")}
)

var (
	ErrPasswordIncorrect      = &Error{Code: 101, Err: errors.New("password incorrect")}
	ErrPasswordSame           = &Error{Code: 102, Err: errors.New("password same as old password")}
	ErrUserNameAlreadyExists  = &Error{Code: 102, Err: errors.New("user name already exists")}
	ErrUserEmailAlreadyExists = &Error{Code: 102, Err: errors.New("user email already exists")}
	ErrUserNotFound           = &Error{Code: 103, Err: errors.New("user not found")}
)

var (
	ErrNodeTypeInvalid       = &Error{Code: 201, Err: errors.New("node type invalid")}
	ErrNodeInputKeyNotFound  = &Error{Code: 202, Err: errors.New("node input key not found")}
	ErrNodeInputTypeInvalid  = &Error{Code: 203, Err: errors.New("node input type invalid")}
	ErrNodeInputInvalid      = &Error{Code: 204, Err: errors.New("node input invalid")}
	ErrNodeOutputKeyNotFound = &Error{Code: 205, Err: errors.New("node output key not found")}
)

var (
	ErrWorkflowFormat = &Error{Code: 501, Err: errors.New("workflow format invalid")}
)

var (
	ErrAvatarNotFound = &Error{Code: 301, Err: errors.New("avatar not found")}
)
