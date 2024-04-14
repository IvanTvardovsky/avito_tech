CREATE TABLE bannermatrix (
    tag_id INTEGER,
    feature_id INTEGER,
    banner_id INTEGER,
    PRIMARY KEY (tag_id, feature_id),
    content json,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN
);
CREATE TABLE test_bannermatrix (
    tag_id INTEGER,
    feature_id INTEGER,
    banner_id INTEGER,
    PRIMARY KEY (tag_id, feature_id),
    content json,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN
);
