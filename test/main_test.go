package test

import (
	"avito_tech/internal/structures"
	"avito_tech/pkg/logger"
	"avito_tech/test/config_test"
	"avito_tech/test/database_test"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"
)

type GetUserBannerTestCase struct {
	ID           string
	ExpectedCode int
	Request      structures.RequestBanner
	Token        string
	Result       map[string]interface{}
}

func TestIntegrationGetUserBanner(t *testing.T) {
	cfg := config_test.GetTestConfig()
	db := database_test.InitAndFillTestDB(cfg)
	time.Sleep(10 * time.Second)

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Log.Fatal(err)
		}
	}(db)

	httpClient := &http.Client{}

	cases := []GetUserBannerTestCase{
		{
			ID:           "1",
			ExpectedCode: http.StatusOK,
			Request: structures.RequestBanner{
				TagID:           1,
				FeatureID:       111,
				UseLastRevision: false,
			},
			Result: map[string]interface{}{
				"title": "hello im test1",
				"url":   "google.com",
			},
			Token: "user_token",
		},
		{
			ID:           "2",
			ExpectedCode: http.StatusOK,
			Request: structures.RequestBanner{
				TagID:           1,
				FeatureID:       111,
				UseLastRevision: true,
			},
			Result: map[string]interface{}{
				"title": "hello im test1",
				"url":   "google.com",
			},
			Token: "user_token",
		},
		{
			ID:           "3",
			ExpectedCode: http.StatusOK,
			Request: structures.RequestBanner{
				TagID:           2,
				FeatureID:       222,
				UseLastRevision: false,
			},
			Result: map[string]interface{}{},
			Token:  "user_token",
		},
		{
			ID:           "4",
			ExpectedCode: http.StatusUnauthorized,
			Request: structures.RequestBanner{
				TagID:           2,
				FeatureID:       222,
				UseLastRevision: false,
			},
			Result: map[string]interface{}{},
			Token:  "not_user_token",
		},
		{
			ID:           "5",
			ExpectedCode: http.StatusOK,
			Request: structures.RequestBanner{
				TagID:           2,
				FeatureID:       222,
				UseLastRevision: true,
			},
			Result: map[string]interface{}{
				"title": "hello im inactive test2",
				"url":   "google2.com",
			},
			Token: "admin_token",
		},
		{
			ID:           "6",
			ExpectedCode: http.StatusUnauthorized,
			Request: structures.RequestBanner{
				TagID:           2,
				FeatureID:       222,
				UseLastRevision: true,
			},
			Result: map[string]interface{}{},
			Token:  "",
		},
		{
			ID:           "7",
			ExpectedCode: http.StatusNotFound,
			Request: structures.RequestBanner{
				TagID:           11111111,
				FeatureID:       22222222,
				UseLastRevision: false,
			},
			Result: map[string]interface{}{},
			Token:  "admin_token",
		},
		{
			ID:           "8",
			ExpectedCode: http.StatusNotFound,
			Request: structures.RequestBanner{
				TagID:           11111111,
				FeatureID:       22222222,
				UseLastRevision: true,
			},
			Result: map[string]interface{}{},
			Token:  "user_token",
		},
	}

	for _, tc := range cases {
		t.Run(tc.ID, func(t *testing.T) {
			path := "/user_banner?tag_id=" + strconv.Itoa(tc.Request.TagID) + "&feature_id=" + strconv.Itoa(tc.Request.FeatureID) +
				"&use_last_revision=" + strconv.FormatBool(tc.Request.UseLastRevision)

			req, err := http.NewRequest("GET", "http://"+cfg.DockerServer.Host+":"+cfg.DockerServer.Port+path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("token", tc.Token)

			resp, err := httpClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					logger.Log.Fatal(err)
				}
			}(resp.Body)

			if resp.StatusCode != tc.ExpectedCode {
				t.Errorf("Expected status code %d, got %d", tc.ExpectedCode, resp.StatusCode)
			}

			var body map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			for key, value := range tc.Result {
				if bodyValue, ok := body[key]; !ok || bodyValue != value {
					t.Errorf("Expected value %v for key %s, got %v", value, key, bodyValue)
				}
			}
		})
	}
	time.Sleep(10 * time.Second)
	database_test.ClearDatabase(db)
}
