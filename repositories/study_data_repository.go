package repository

import (
	"2026-FM247-BackEnd/models"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StudyDataRepository struct {
	db    *gorm.DB
	redis *redis.Client
	ctx   context.Context
}

func NewStudyDataRepository(db *gorm.DB, redis *redis.Client) *StudyDataRepository {
	return &StudyDataRepository{
		db:    db,
		redis: redis,
		ctx:   context.Background(),
	}
}

// 生成redis的每日key
func (r *StudyDataRepository) generateDailyKey(userID uint, date time.Time) string {
	return fmt.Sprintf("user:%d:studydata:date:%s", userID, date.Format("2006-01-02"))
}

// 生成redis的每月key
func (r *StudyDataRepository) generateMonthlyKey(userID uint, date time.Time) string {
	return fmt.Sprintf("user:%d:studydata:month:%s", userID, date.Format("2006-01"))
}

//redis中仅包含每日学习时长数据，每月将在mysql中聚合
//数据结构： key：user:{userID}:studydata:date:{YYYY-MM-DD}
// 			study_time：学习时长 单位分钟		tomatoes：番茄钟数量

// 增加每日和每月学习时长
func (r *StudyDataRepository) IncrementDailyStudyTime(userID uint, date time.Time, studyTime int) error {
	dailykey := r.generateDailyKey(userID, date)
	monthlykey := r.generateMonthlyKey(userID, date)

	pipe := r.redis.Pipeline()
	pipe.HIncrBy(r.ctx, dailykey, "study_time", int64(studyTime))
	pipe.HIncrBy(r.ctx, monthlykey, "study_time", int64(studyTime))

	//设置并只设置一次过期时间(NX)
	pipe.ExpireNX(r.ctx, dailykey, 25*time.Hour)
	pipe.ExpireNX(r.ctx, monthlykey, 32*24*time.Hour)

	_, err := pipe.Exec(r.ctx)
	return err
}

// 增加每日和每月番茄钟次数，随后同步mysql
func (r *StudyDataRepository) IncrementDailyTomatoes(userID uint, date time.Time, tomatoes int) error {
	dailykey := r.generateDailyKey(userID, date)
	monthlykey := r.generateMonthlyKey(userID, date)

	pipe := r.redis.Pipeline()
	pipe.HIncrBy(r.ctx, dailykey, "tomatoes", int64(tomatoes))
	pipe.HIncrBy(r.ctx, monthlykey, "tomatoes", int64(tomatoes))

	//其实可以不写来着，以防万一
	pipe.ExpireNX(r.ctx, dailykey, 25*time.Hour)
	pipe.ExpireNX(r.ctx, monthlykey, 32*24*time.Hour)

	_, err := pipe.Exec(r.ctx)
	if err != nil {
		return err
	}

	return r.SyncDailyDataToMySQL(userID, date)
}

// 同步数据到mysql
func (r *StudyDataRepository) SyncDailyDataToMySQL(userID uint, date time.Time) error {
	dailykey := r.generateDailyKey(userID, date)
	monthlykey := r.generateMonthlyKey(userID, date)

	daydata, err := r.redis.HGetAll(r.ctx, dailykey).Result()
	if err != nil {
		return err
	}
	if len(daydata) == 0 {
		return nil
	}
	monthdata, err := r.redis.HGetAll(r.ctx, monthlykey).Result()
	if err != nil {
		return err
	}

	//每日的日期格式xxxx-xx-xx-00-00-00
	daystudytime, _ := strconv.Atoi(daydata["study_time"])
	daytomatoes, _ := strconv.Atoi(daydata["tomatoes"])
	dailyData := models.DailyStudyData{
		UserID:    userID,
		Date:      time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()),
		StudyTime: daystudytime,
		Tomatoes:  daytomatoes,
	}
	//每月的日期格式xxxx-xx-01-00-00-00
	monthstudytime, _ := strconv.Atoi(monthdata["study_time"])
	monthtomatoes, _ := strconv.Atoi(monthdata["tomatoes"])
	monthlyData := models.MonthlyStudyData{
		UserID:    userID,
		Month:     time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location()),
		StudyTime: monthstudytime,
		Tomatoes:  monthtomatoes,
	}

	tx := r.db.Begin()
	//更新每日数据
	err = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"study_time", "tomatoes", "end_at", "updated_at"}),
	}).Create(&dailyData).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	//更新每月数据
	err = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "month"}},
		DoUpdates: clause.AssignmentColumns([]string{"study_time", "tomatoes", "end_at", "updated_at"}),
	}).Create(&monthlyData).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
