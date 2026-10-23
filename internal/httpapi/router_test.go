package httpapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"recipe-sharing/internal/auth"
	"recipe-sharing/internal/config"
	"recipe-sharing/internal/handler"
	"recipe-sharing/internal/httpapi"
	"recipe-sharing/internal/repository"
	"recipe-sharing/internal/service"
	"strings"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func performRequest(
	handler http.Handler,
	method string,
	path string,
	body string,
	token string,
) *httptest.ResponseRecorder {
	req := httptest.NewRequest(
		method,
		path,
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	return rec
}

func setupTestRouter(t *testing.T) http.Handler {
	t.Helper()

	os.Setenv("JWT_KEY", "test-jwt-secret-key")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	auth.InitKey(cfg.JWTKey)
	db, _ := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})

	recipeRepository := repository.NewRecipeRepository(db)
	recipeService := service.NewRecipeService(recipeRepository)
	recipeHandler := handler.NewRecipeHandler(recipeService)

	ratingRepo := repository.NewRatingRepository(db)
	ratingService := service.NewRatingService(ratingRepo, recipeRepository)
	ratingHandler := handler.NewRatingHandler(ratingService)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	authHandler := handler.NewAuthHandler(userService)

	favoriteRepo := repository.NewFavoriteRepository(db)
	favoriteService := service.NewFavoriteService(favoriteRepo)
	favoriteHandler := handler.NewFavoriteHandler(favoriteService)

	return httpapi.NewRouter(recipeHandler, authHandler, favoriteHandler, ratingHandler)
}

func TestGetRecipes(t *testing.T) {
	router := setupTestRouter(t)

	res := performRequest(
		router,
		http.MethodGet,
		"/api/recipes",
		"",
		"",
	)

	if res.Code != http.StatusOK {
		t.Fatalf(
			"expect status %d, got %d; body: %s",
			http.StatusOK,
			res.Code,
			res.Body.String(),
		)
	}
}

func TestGetNonexistentRecipe(t *testing.T) {
	router := setupTestRouter(t)

	res := performRequest(
		router,
		http.MethodGet,
		"/api/recipes/999999",
		"",
		"",
	)

	if res.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusNotFound,
			res.Code,
			res.Body.String(),
		)
	}
}

func TestCreateRecipeWithoutToken(t *testing.T) {
	router := setupTestRouter(t)

	res := performRequest(
		router,
		http.MethodPost,
		"/api/recipes",
		`{"title":"Test recipe","description":"Test","ingredients":"Flour"}`,
		"",
	)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusUnauthorized,
			res.Code,
			res.Body.String(),
		)
	}
}

func TestDeleteRecipeWithoutToken(t *testing.T) {
	router := setupTestRouter(t)

	res := performRequest(
		router,
		http.MethodDelete,
		"/api/recipes/1",
		"",
		"",
	)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusUnauthorized,
			res.Code,
			res.Body.String(),
		)
	}
}

func TestRegisterLoginAndCreateRecipe(t *testing.T) {
	router := setupTestRouter(t)

	email := fmt.Sprintf(
		"test-%d@example.com",
		time.Now().UnixNano(),
	)

	registerBody := fmt.Sprintf(`{
        "username": "testuser",
        "email": %q,
        "password": "TestPassword123!"
    }`, email)

	registerRes := performRequest(
		router, http.MethodPost,
		"/api/auth/register", registerBody, "",
	)

	if registerRes.Code != http.StatusCreated {
		t.Fatalf("registration failed: %d: %s",
			registerRes.Code, registerRes.Body.String())
	}

	loginBody := fmt.Sprintf(`{
        "username": "testuser",
        "email": %q,
        "password": "TestPassword123!"
    }`, email)

	loginRes := performRequest(
		router, http.MethodPost,
		"/api/auth/login", loginBody, "",
	)

	if loginRes.Code != http.StatusOK {
		t.Fatalf("login failed: %d: %s",
			loginRes.Code, loginRes.Body.String())
	}

	var loginData struct {
		Token string `json:"token"`
	}

	if err := json.Unmarshal(loginRes.Body.Bytes(), &loginData); err != nil {
		t.Fatalf("invalid login response: %v", err)
	}

	if loginData.Token == "" {
		t.Fatal("login response contains an empty token")
	}

	recipeBody := `{
		"title": "Test Pancakes",
		"description": "A simple recipe",
		"ingredients": "Flour, milk, eggs"
	}`

	createRes := performRequest(
		router, http.MethodPost,
		"/api/recipes", recipeBody, loginData.Token,
	)

	if createRes.Code != http.StatusCreated {
		t.Fatalf("recipe creation failed: %d: %s",
			createRes.Code, createRes.Body.String())
	}
}
