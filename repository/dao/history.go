package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

const (
	HistoryStatusHidden  = 1
	HistoryStatusVisible = 0
)

type SearchHistory struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Uid      int64  `gorm:"uniqueIndex:uid_location_keyword"`
	Location int32  `gorm:"uniqueIndex:uid_location_keyword"`
	Keyword  string `gorm:"uniqueIndex:uid_location_keyword; type:varchar(30)"`
	Status   int32
	Ctime    int64
	Utime    int64 `gorm:"index"`
}

type SearchHistoryGORMDAO struct {
	db *gorm.DB
}

func NewSearchHistoryGORMDAO(db *gorm.DB) HistoryDAO {
	return &SearchHistoryGORMDAO{db: db}
}

func (dao *SearchHistoryGORMDAO) Insert(ctx context.Context, history SearchHistory) error {
	now := time.Now().UnixMilli()
	history.Ctime = now
	history.Utime = now
	return dao.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"utime":  now,
				"status": HistoryStatusVisible,
			})}).
		Create(&history).Error
}

func (dao *SearchHistoryGORMDAO) GetUserSearchHistories(ctx context.Context, uid int64, location int32, offset int64, limit int64) ([]SearchHistory, error) {
	var hs []SearchHistory
	err := dao.db.WithContext(ctx).
		Where("uid = ? and location = ?", uid, location).
		Order("utime desc").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&hs).Error
	return hs, err
}

func (dao *SearchHistoryGORMDAO) HideAllUserSearchHistories(ctx context.Context, uid int64, location int32) error {
	return dao.db.WithContext(ctx).
		Model(&SearchHistory{}).
		Where("uid = ? and location = ?", uid, location).
		Updates(map[string]any{
			"utime":  time.Now().UnixMilli(),
			"status": HistoryStatusHidden,
		}).Error

}

func (dao *SearchHistoryGORMDAO) HideUserSearchHistoriesByIds(ctx context.Context, uid int64, location int32, historyIds []int64) error {
	return dao.db.WithContext(ctx).
		Model(&SearchHistory{}).
		Where("uid = ? and location = ? and id in ?", uid, location, historyIds).
		Updates(map[string]any{
			"utime":  time.Now().UnixMilli(),
			"status": HistoryStatusHidden,
		}).Error
}
