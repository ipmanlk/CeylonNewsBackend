package sqldb

import (
	"errors"
	"ipmanlk/cnapi/common"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB() error {
	dsn := os.Getenv("MYSQL_DSN")

	if dsn == "" {
		dsn = "root:@tcp(0.0.0.0:3306)/cn?charset=utf8mb4&parseTime=True&loc=Local"
	}

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		 Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		return err
	}

	err = db.AutoMigrate(&common.NewsItem{})
	if err != nil {
		return err
	}

	// Setup full text search index
	// Check if the index already exists
	var indexCount int64
	result := db.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.statistics 
		WHERE table_schema = DATABASE() 
		AND table_name = 'news_items' 
		AND index_name = 'idx_fulltext_search'
	`).Count(&indexCount)

	if result.Error != nil {
		return result.Error
	}

	// If the index doesn't exist, create it
	if indexCount == 0 {
		if err := db.Exec(`
			CREATE FULLTEXT INDEX idx_fulltext_search 
			ON news_items (title, content_text)
			`).Error; err != nil {
			return err
		}
	}

	return nil
}

func InsertItem(item common.NewsItem) error {
	err := db.Create(&item).Error
	if err != nil {
		return err
	}

	return nil
}

// PlanetScale has issues with this.
// 
// func InsertItems(items []common.NewsItem) error {
// 	tx := db.Begin()
// 	for _, item := range items {
// 		// Use InsertIgnore to skip duplicate entry errors
// 		if err := tx.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&item).Error; err != nil {
// 			// Other errors, rollback the transaction and return the error
// 			tx.Rollback()
// 			return err
// 		}
// 	}
// 	return tx.Commit().Error
// }
// 
func InsertItems(items []common.NewsItem) error {
    for _, item := range items {
        if err := db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&item).Error; err != nil {
            if !errors.Is(err, gorm.ErrRecordNotFound) {
                return err
            }
        }
    }
    return nil
}

func GetItemByID(id uint) (common.NewsItem, error) {
	var item common.NewsItem
	err := db.First(&item, id).Error
	if err != nil {
		return item, err
	}

	return item, nil
}

func SearchItemsByLangSources(langs []common.Lang, sources []string, query string, page, pageSize int) ([]common.NewsItem, error) {
	var items []common.NewsItem

	dbQuery := db.
		Select("id", "title", "created_at", "language", "source_name", "url", "thumbnail_url").
		Where("language IN ?", langs).
		Where("source_name IN ?", sources).
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize)

	if query != "" {
		dbQuery = dbQuery.Where("MATCH(title, content_text) AGAINST (? IN BOOLEAN MODE)", query)
	}

	err := dbQuery.Find(&items).Error

	if err != nil {
		return items, err
	}

	return items, nil
}
