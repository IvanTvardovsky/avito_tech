package handlers

import (
	"avito_tech/internal/structures"
	"avito_tech/pkg/logger"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

var BannerCounter = atomic.Int64{}

func GetUserBanner(c *gin.Context, db *sql.DB, cfg *structures.Config, matrix *map[int]map[int]int, bannerMap *map[int]*structures.Banner) {
	var tagInt, featureInt int
	var useLastRevision bool
	var err error

	userTokenType, _ := c.Get("token_type")

	tagString := c.Query("tag_id")
	if tagString != "" {
		tagInt, err = strconv.Atoi(tagString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	featureString := c.Query("feature_id")
	if featureString != "" {
		featureInt, err = strconv.Atoi(featureString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	useLastRevisionString := c.Query("use_last_revision")
	if useLastRevisionString != "" {
		useLastRevision, err = strconv.ParseBool(useLastRevisionString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if useLastRevision {
		var query string
		if cfg.AppMode.IsTest {
			query = `SELECT content, is_active FROM test_bannermatrix WHERE tag_id = $1 AND feature_id = $2`
		} else {
			query = `SELECT content, is_active FROM bannermatrix WHERE tag_id = $1 AND feature_id = $2`
		}

		var contentJSON string
		var isActive bool
		err := db.QueryRow(query, tagInt, featureInt).Scan(&contentJSON, &isActive)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Баннер с указанным ID не найден"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !isActive && userTokenType == "user" {
			c.JSON(http.StatusOK, gin.H{})
			return
		}

		var content map[string]interface{}
		if err := json.Unmarshal([]byte(contentJSON), &content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, content)
	} else {
		bannerID, ok := (*matrix)[tagInt][featureInt]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Баннер с указанным ID не найден"})
			return
		}
		banner := *((*bannerMap)[bannerID])
		if !banner.IsActive && userTokenType == "user" {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		c.JSON(http.StatusOK, banner.Content)
	}
}

func fillBannersSlice(tagString, featureString string, tagInt, featureInt int, matrix *map[int]map[int]int,
	bannerMap *map[int]*structures.Banner, wasPlacedBanner *map[int]struct{}, banners *[]structures.Banner) {
	if matrix == nil {
		return
	}

	switch {
	case tagString != "" && featureString != "":
		if tagMap, ok := (*matrix)[tagInt]; ok {
			if bannerID, ok := tagMap[featureInt]; ok {
				if _, ok := (*wasPlacedBanner)[bannerID]; !ok {
					(*wasPlacedBanner)[bannerID] = struct{}{}
					*banners = append(*banners, *((*bannerMap)[bannerID]))
				}
			}
		}
	case tagString != "" && featureString == "":
		if tagMap, ok := (*matrix)[tagInt]; ok {
			for _, v := range tagMap {
				if _, ok := (*wasPlacedBanner)[v]; !ok {
					(*wasPlacedBanner)[v] = struct{}{}
					*banners = append(*banners, *((*bannerMap)[v]))
				}
			}
		}
	case tagString == "" && featureString != "":
		for _, tagValue := range *matrix {
			for featureKey, bannerID := range tagValue {
				if featureKey == featureInt {
					if _, ok := (*wasPlacedBanner)[bannerID]; !ok {
						(*wasPlacedBanner)[bannerID] = struct{}{}
						*banners = append(*banners, *((*bannerMap)[bannerID]))
					}
				}
			}
		}
	default:
		for _, tagValue := range *matrix {
			for _, bannerID := range tagValue {
				if _, ok := (*wasPlacedBanner)[bannerID]; !ok {
					(*wasPlacedBanner)[bannerID] = struct{}{}
					*banners = append(*banners, *((*bannerMap)[bannerID]))
				}
			}
		}
	}
}

func GetAllBanners(c *gin.Context, matrix *map[int]map[int]int, bannerMap *map[int]*structures.Banner) {
	var tagInt, featureInt, limitInt, offsetInt int
	var err error

	tagString := c.Query("tag_id")
	if tagString != "" {
		tagInt, err = strconv.Atoi(tagString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	featureString := c.Query("feature_id")
	if featureString != "" {
		featureInt, err = strconv.Atoi(featureString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	limitString := c.DefaultQuery("limit", "-1")
	limitInt, err = strconv.Atoi(limitString)
	if err != nil || limitInt < -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректное значение параметра 'limit'"})
		return
	}

	offsetString := c.DefaultQuery("offset", "0")
	offsetInt, err = strconv.Atoi(offsetString)
	if err != nil || offsetInt < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректное значение параметра 'offset'"})
		return
	}

	var wasPlacedBanner = make(map[int]struct{})
	var banners []structures.Banner

	fillBannersSlice(tagString, featureString, tagInt, featureInt, matrix, bannerMap, &wasPlacedBanner, &banners)

	if offsetInt >= len(banners) {
		c.JSON(http.StatusOK, []structures.Banner{})
		return
	}

	end := offsetInt + limitInt
	if limitInt == -1 || end > len(banners) {
		end = len(banners)
	}

	c.JSON(http.StatusOK, banners[offsetInt:end])
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

func PatchBanner(c *gin.Context, db *sql.DB) {
	idString := c.Param("id")
	if idString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Нет ID в параметрах"})
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := `SELECT EXISTS(SELECT 1 FROM bannermatrix WHERE banner_id = $1)`
	var exists bool
	err = db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выполнении запроса в базу данных"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Баннер с указанным ID не найден"})
		return
	}

	var request structures.Banner
	err = c.ShouldBindJSON(&request)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var columns []string
	var values []interface{}
	counter := 1

	if len(request.TagIDs) > 0 {
		columns = append(columns, fmt.Sprintf("tag_ids = $%d", counter))
		values = append(values, pq.Array(request.TagIDs))
		counter++
	}
	if request.FeatureID != 0 {
		columns = append(columns, fmt.Sprintf("feature_id = $%d", counter))
		values = append(values, request.FeatureID)
		counter++
	}
	if len(request.Content) > 0 {
		contentJSON, er := json.Marshal(request.Content)
		if er != nil {
			logger.Log.Error(er)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при кодировании содержимого баннера в JSON"})
			return
		}

		columns = append(columns, fmt.Sprintf("content = $%d", counter))
		values = append(values, contentJSON)
		counter++
	}
	if request.IsActive {
		columns = append(columns, fmt.Sprintf("is_active = $%d", counter))
		values = append(values, request.IsActive)
		counter++
	}

	if len(columns) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"description": "Нет обновляемых полей"})
		return
	}

	query = fmt.Sprintf("UPDATE bannermatrix SET updated_at = CURRENT_TIMESTAMP, %s WHERE banner_id = $%d", strings.Join(columns, ", "), counter)
	values = append(values, id)

	_, err = db.Exec(query, values...)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Баннер успешно обновлен"})
}

func DeleteBanner(c *gin.Context, db *sql.DB) {
	idString := c.Param("id")
	if idString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет ID в параметрах"})
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM bannermatrix WHERE banner_id = $1)", id).Scan(&exists)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при выполнении запроса в базу данных"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Баннер с указанным ID не найден"})
		return
	}

	_, err = db.Exec("DELETE FROM bannermatrix WHERE banner_id = $1", id)
	if err != nil {
		logger.Log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении баннера из базы данных"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func GetBannerVersions(c *gin.Context, db *sql.DB) {

}
