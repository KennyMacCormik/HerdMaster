package middleware

import (
	"errors"
	"github.com/KennyMacCormik/HerdMaster/pkg/val"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddleware_GenerateUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/", func(c *gin.Context) {
		requestID := c.Writer.Header().Get(RequestIDKey)
		c.String(http.StatusOK, requestID)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Response status should be 200")
	requestID := w.Body.String()
	assert.NotEmpty(t, requestID, "Request ID should be generated and returned")
	assert.NoError(t, val.GetValidator().ValidateWithTag(requestID, "uuid4_rfc4122"),
		"Request ID should be valid UUID version 4 based on RFC 4122 ")
}

func TestRequestIDMiddleware_UseExistingRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/", func(c *gin.Context) {
		requestID := c.Writer.Header().Get(RequestIDKey)
		c.String(http.StatusOK, requestID)
	})

	existingRequestID := genUuid()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(RequestIDKey, existingRequestID)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Response status should be 200")
	assert.Equal(t, existingRequestID, w.Body.String(), "Existing Request ID should be used")
}

func TestGetRequestIDFromCtx_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestIDMiddleware())
	router.GET("/", func(c *gin.Context) {
		requestID, err := GetRequestIDFromCtx(c)
		assert.NoError(t, err, "GetRequestIDFromCtx should not return an error")
		assert.NotEmpty(t, requestID, "Request ID should be retrieved successfully")
		c.String(http.StatusOK, requestID)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Response status should be 200")
	assert.NotEmpty(t, w.Body.String(), "Request ID should be present in the response")
}

func TestGetRequestIDFromCtx_Fallback(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Next()
	})
	router.GET("/", func(c *gin.Context) {
		requestID, err := GetRequestIDFromCtx(c)
		assert.Error(t, err, "GetRequestIDFromCtx should return an error when no Request ID is present")
		var fallbackErr *ErrFallbackUuidUsed
		assert.True(t, errors.As(err, &fallbackErr), "Error should be of type ErrFallbackUuidUsed")
		c.String(http.StatusOK, requestID)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)

	requestID := w.Body.String()
	assert.Equal(t, http.StatusOK, w.Code, "Response status should be 200")
	assert.NotEmpty(t, requestID, "Fallback Request ID should be present in the response")
	assert.NoError(t, val.GetValidator().ValidateWithTag(requestID, "uuid4_rfc4122"),
		"Request ID should be valid UUID version 4 based on RFC 4122 ")
}

func TestGetRequestIDFromCtx_TypeCastFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(RequestIDKey, 12345) // Simulate a type mismatch
		c.Next()
	})
	router.GET("/", func(c *gin.Context) {
		_, err := GetRequestIDFromCtx(c)
		assert.Error(t, err, "GetRequestIDFromCtx should return an error when type cast fails")
		var typeCastErr *ErrTypeCastFailed
		assert.True(t, errors.As(err, &typeCastErr), "Error should be of type ErrTypeCastFailed")
		c.String(http.StatusInternalServerError, "type cast error")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code, "Response status should be 500 on type cast failure")
}

func TestFallbackUuidHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		requestID := fallbackUuidHandler(c)
		assert.NotEmpty(t, requestID, "Fallback UUID should be generated")
		assert.NoError(t, val.GetValidator().ValidateWithTag(requestID, "uuid4_rfc4122"),
			"Request ID should be valid UUID version 4 based on RFC 4122 ")
		reqIdFromCtx, ok := c.Get(RequestIDKey)
		assert.True(t, ok, "Request ID should be retrieved from context")
		assert.NotEmpty(t, reqIdFromCtx, "Request ID should be retrieved from context")
		requestID2, ok := reqIdFromCtx.(string)
		assert.True(t, ok, "Request ID successfully cast to string")
		assert.Equal(t, requestID, requestID2, "Generated request ID should be stored in context")
		requestID3 := c.Writer.Header().Get(RequestIDKey)
		assert.NotEmpty(t, requestID3, "Request ID should be present writer header")
		assert.Equal(t, requestID, requestID3, "Generated request ID should be present writer header")

		c.String(http.StatusOK, requestID)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Response status should be 200")
	assert.Len(t, w.Body.String(), 36, "Fallback UUID in response should be valid")
}

func TestGenUuid(t *testing.T) {
	requestID := genUuid()
	assert.NotEmpty(t, requestID, "UUID should be generated")
	assert.NoError(t, val.GetValidator().ValidateWithTag(requestID, "uuid4_rfc4122"),
		"Request ID should be valid UUID version 4 based on RFC 4122 ")
}

func TestErrTypeCastFailed(t *testing.T) {
	err := NewErrTypeCastFailed("test msg")
	assert.NotEmpty(t, err, "Should be not empty")
	assert.Error(t, err, "Should be error")
	assert.IsType(t, &ErrTypeCastFailed{}, err, "Error should be of type ErrTypeCastFailed")
	assert.Equal(t, "test msg", err.Error(), "Should contain test msg")
}

func TestErrFallbackUuidUsed(t *testing.T) {
	err := NewErrFallbackUuidUsed("test msg")
	assert.NotEmpty(t, err, "Should be not empty")
	assert.Error(t, err, "Should be error")
	assert.IsType(t, &ErrFallbackUuidUsed{}, err, "Error should be of type ErrFallbackUuidUsed")
	assert.Equal(t, "test msg", err.Error(), "Should contain test msg")
}
