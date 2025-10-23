package api

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	NotFound()(ctx)

	assert.Equal(t, 404, w.Code)
	assert.JSONEq(t, `{"errorCode":404, "errorDescription":"route not found", "error":"Not Found"}`, w.Body.String())
}
