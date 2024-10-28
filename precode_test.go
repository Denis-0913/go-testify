package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое.
func TestMainHandlerWhenCorrectRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=4&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	body := responseRecorder.Body.String()
	statusCode := responseRecorder.Code

	require.Equalf(t, http.StatusOK, statusCode, "expected status code: %d, got %d", http.StatusOK, statusCode) // сервис возвращает код ответа 200
	require.NotEmpty(t, body)                                                                                   //тело не пустое
}

// Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestMainHandlerWhenIncorrectCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=4&city=piter", nil)
	//city := req.URL.Query().Get("city")

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	body := responseRecorder.Body.String()
	statusCode := responseRecorder.Code
	expected := "wrong city value"

	assert.Equalf(t, http.StatusBadRequest, statusCode, "expected status code: %d, got %d", http.StatusBadRequest, statusCode) // сервис возвращает код ответа 400
	require.Equalf(t, expected, body, "expected body: %s, got %s", expected, body)                                             //wrong city value в теле ответа.
}

// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе.
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")

	assert.Lenf(t, list, totalCount, "Expected cafe count: %d, got %d", totalCount, len(list))
}
