package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"media-content-api/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestGetMovies(t *testing.T) {
	router := setupTestRouter()
	router.GET("/api/movies", GetMovies)

	req, _ := http.NewRequest("GET", "/api/movies", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateMovie(t *testing.T) {
	router := setupTestRouter()
	router.POST("/api/movies", CreateMovie)

	movie := models.Movie{
		Title:       "Test Movie",
		Description: "Test Description",
		ReleaseYear: 2024,
		Duration:    120,
		Rating:      8.5,
		Genre:       "Action",
		Director:    "Test Director",
	}

	jsonValue, _ := json.Marshal(movie)
	req, _ := http.NewRequest("POST", "/api/movies", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetMovieByID(t *testing.T) {
	router := setupTestRouter()
	router.GET("/api/movies/:id", GetMovie)

	req, _ := http.NewRequest("GET", "/api/movies/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateMovie(t *testing.T) {
	router := setupTestRouter()
	router.PUT("/api/movies/:id", UpdateMovie)

	updateData := map[string]interface{}{
		"rating": 9.0,
	}
	jsonValue, _ := json.Marshal(updateData)

	req, _ := http.NewRequest("PUT", "/api/movies/1", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteMovie(t *testing.T) {
	router := setupTestRouter()
	router.DELETE("/api/movies/:id", DeleteMovie)

	req, _ := http.NewRequest("DELETE", "/api/movies/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetSeries(t *testing.T) {
	router := setupTestRouter()
	router.GET("/api/series", GetSeries)

	req, _ := http.NewRequest("GET", "/api/series", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateSeries(t *testing.T) {
	router := setupTestRouter()
	router.POST("/api/series", CreateSeries)

	series := models.Series{
		Title:       "Test Series",
		Description: "Test Description",
		ReleaseYear: 2024,
		Seasons:     1,
		Episodes:    10,
		Rating:      8.0,
		Genre:       "Drama",
		Creator:     "Test Creator",
	}

	jsonValue, _ := json.Marshal(series)
	req, _ := http.NewRequest("POST", "/api/series", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRegisterUser(t *testing.T) {
	router := setupTestRouter()
	router.POST("/api/auth/register", Register)

	user := map[string]string{
		"username": "testuser",
		"email":    "test@test.com",
		"password": "test123",
	}
	jsonValue, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLoginUser(t *testing.T) {
	router := setupTestRouter()
	router.POST("/api/auth/login", Login)

	user := map[string]string{
		"email":    "test@test.com",
		"password": "test123",
	}
	jsonValue, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetExternalMovie(t *testing.T) {
	router := setupTestRouter()
	router.GET("/api/external/movie/:title", GetExternalMovie)

	req, _ := http.NewRequest("GET", "/api/external/movie/Inception", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
