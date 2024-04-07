package structures

type Banner struct {
	ID        int                    `json:"banner_id"`  // Идентификатор баннера
	TagIDs    []int                  `json:"tag_ids"`    // Идентификаторы тэгов
	FeatureID int                    `json:"feature_id"` // Идентификатор фичи
	Content   map[string]interface{} `json:"content"`    // Содержимое баннера
	IsActive  bool                   `json:"is_active"`  // Флаг активности баннера
	CreatedAt string                 `json:"created_at"` // Дата создания баннера
	UpdatedAt string                 `json:"updated_at"` // Дата обновления баннера
}

type RequestBanner struct {
	FeatureID int `json:"feature_id"`
	TagID     int `json:"tag_id"`
	Limit     int `json:"limit"`
	Offset    int `json:"offset"`
}
