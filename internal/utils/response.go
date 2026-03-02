package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorCode string

const (
	ErrCodeBadRequest ErrorCode = "BAD_REQUEST"
	ErrCodeNotFound   ErrorCode = "NOT_FOUND"
	ErrCodeConflict   ErrorCode = "CONFLICT"
	ErrCodeInternal   ErrorCode = "INTERNAL_SERVER_ERROR"
)

type AppError struct {
	Message string
	Code    ErrorCode
	Err     error
}

func (ae *AppError) Error() string {
	if ae.Err != nil {
		return ae.Message + ": " + ae.Err.Error()
	}
	return ae.Message
}

// NewError tạo AppError mới với message và mã lỗi nghiệp vụ
func NewError(message string, code ErrorCode) error {
	return &AppError{
		Message: message,
		Code:    code,
	}
}

func WrapError(err error, message string, code ErrorCode) error {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// ResponseError trả response lỗi chuẩn cho client dựa trên AppError hoặc lỗi thường
func ResponseError(ctx *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		status := httpStatusFromCode(appErr.Code)
		response := gin.H{
			"error": appErr.Message,
			"code":  appErr.Code,
		}

		if appErr.Err != nil {
			response["detail"] = appErr.Err.Error()
		}

		ctx.JSON(status, response)
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"code":  ErrCodeInternal,
	})
}

// ResponseSuccess trả response thành công
func ResponseSuccess(ctx *gin.Context, status int, data any) {
	ctx.JSON(status, gin.H{
		"status": "success",
		"data":   data,
	})
}

// ResponseStatusCode trả response chỉ với HTTP status
func ResponseStatusCode(ctx *gin.Context, status int) {
	ctx.Status(status)
}

// ResponseValidator trả response lỗi validate
func ResponseValidator(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusBadRequest, data)
}

// httpStatusFromCode map ErrorCode sang HTTP status tương ứng
func httpStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func ToGRPCStatus(err error) error {
	if err == nil {
		return nil
	}

	if ae, ok := err.(*AppError); ok {
		switch ae.Code {
		case ErrCodeBadRequest:
			return status.Error(codes.InvalidArgument, ae.Message)
		case ErrCodeNotFound:
			return status.Error(codes.NotFound, ae.Message)
		case ErrCodeConflict:
			return status.Error(codes.AlreadyExists, ae.Message)
		default:
			return status.Error(codes.Internal, ae.Message)
		}
	}

	return status.Error(codes.Internal, err.Error())
}
