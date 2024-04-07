package handlers

import (
	"avito_tech/internal/structures"
	"avito_tech/pkg/logger"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"sync/atomic"
)

var BannerCounter = atomic.Int64{}

func GetAllBanners(c *gin.Context) {
	c.Status(http.StatusOK)
}

func CreateBanner(c *gin.Context, db *sql.DB) {
	var banner structures.Banner

	if c.Request.ContentLength == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Некорректные данные"})
		return
	}

	err := c.ShouldBindJSON(&banner)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bannerID := BannerCounter.Load()
	BannerCounter.Add(1)
	tagIDsArray := pq.Array(banner.TagIDs)
	bannerContentJSON, err := json.Marshal(banner.Content)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO BannerMatrix (tag_id, feature_id, banner_id, content, created_at, updated_at, is_active)
			  SELECT
					unnest($1::integer[]) AS tag_id,
					$2 AS feature_id,
					$3 AS banner_id,
					$4::json AS content,
					CURRENT_TIMESTAMP AS created_at,
					CURRENT_TIMESTAMP AS updated_at,
					$5 AS is_active; `

	_, err = db.Exec(query, tagIDsArray, banner.FeatureID, bannerID, bannerContentJSON, banner.IsActive)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"banner_id": bannerID})
}
