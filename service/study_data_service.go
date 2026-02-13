package service

import (
	"2026-FM247-BackEnd/models"
	"time"
)

type StudyDataRepository interface {
	GenerateDailyKey(userID uint, date time.Time) string
	GenerateMonthlyKey(userID uint, date time.Time) string
	IncrementDailyStudyTime(userID uint, date time.Time, studyTime int) error
	IncrementDailyTomatoes(userID uint, date time.Time, tomatoes int) error
	SyncDailyDataToMySQL(userID uint, date time.Time, addStudyTime int, addTomatoes int) error
	GetDailyStudyData(userID uint, date time.Time) (*models.DailyStudyData, error)
	GetMonthlyStudyData(userID uint, date time.Time) (*models.MonthlyStudyData, error)
	GetTotalStudyData(userID uint) (*models.TotalStudyData, error)
	GetStudyDataSummary(userID uint, startDate, endDate time.Time) ([]models.DailyStudyData, error)
}

type StudyDataService struct {
	repo StudyDataRepository
}

func NewStudyDataService(repo StudyDataRepository) *StudyDataService {
	return &StudyDataService{repo: repo}
}

// 增加学习时长、番茄钟次数
func (s *StudyDataService) AddStudyData(userID uint, date time.Time, studyTime int, tomatoes int) (bool, string) {
	err := s.repo.IncrementDailyStudyTime(userID, date, studyTime)
	if err != nil {
		return false, "记录学习时长失败" + err.Error()
	}

	err = s.repo.IncrementDailyTomatoes(userID, date, tomatoes)
	if err != nil {
		return false, "记录番茄钟次数失败" + err.Error()
	}

	err = s.repo.SyncDailyDataToMySQL(userID, date, studyTime, tomatoes)
	if err != nil {
		return false, "同步数据到MySQL失败" + err.Error()
	}

	return true, "学习数据记录成功"
}

// 获取每日学习数据
func (s *StudyDataService) GetDailyStudyData(userID uint, date time.Time) (DailyStudyDataInfo, string) {
	data, err := s.repo.GetDailyStudyData(userID, date)
	if err != nil {
		return DailyStudyDataInfo{}, "查询每日学习数据失败" + err.Error()
	}
	return DailyStudyDataInfo{
		Date:      data.Date,
		StudyTime: data.StudyTime,
		Tomatoes:  data.Tomatoes,
	}, ""
}

// 获取每月学习数据
func (s *StudyDataService) GetMonthlyStudyData(userID uint, date time.Time) (MonthlyStudyDataInfo, string) {
	data, err := s.repo.GetMonthlyStudyData(userID, date)
	if err != nil {
		return MonthlyStudyDataInfo{}, "查询每月学习数据失败" + err.Error()
	}
	return MonthlyStudyDataInfo{
		Month:     data.Month,
		StudyTime: data.StudyTime,
		Tomatoes:  data.Tomatoes,
	}, ""
}

// 获取总学习数据
func (s *StudyDataService) GetTotalStudyData(userID uint) (TotalStudyDataInfo, string) {
	data, err := s.repo.GetTotalStudyData(userID)
	if err != nil {
		return TotalStudyDataInfo{}, "查询总学习数据失败" + err.Error()
	}
	return TotalStudyDataInfo{
		StudyTime: data.StudyTime,
		Tomatoes:  data.Tomatoes,
	}, ""
}

// 获取本周每日学习数据
func (s *StudyDataService) GetWeekStudyData(userID uint, date time.Time) ([]DailyStudyDataInfo, string) {
	today := int(date.Weekday())
	if today == 0 {
		today = 7
	}
	monday := date.AddDate(0, 0, -today+1)

	var result []DailyStudyDataInfo
	for i := 0; i < 7; i++ {
		date := monday.AddDate(0, 0, i)
		data, msg := s.GetDailyStudyData(userID, date)
		if msg != "" {
			data = DailyStudyDataInfo{
				Date:      date,
				StudyTime: 0,
				Tomatoes:  0,
			}
		}
		result = append(result, data)
	}

	return result, ""
}

// 获取本月每日学习数据
func (s *StudyDataService) GetMonthStudyData(userID uint, date time.Time) ([]DailyStudyDataInfo, string) {
	firstDay := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	daysInMonth := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
	var result []DailyStudyDataInfo
	for i := 0; i < daysInMonth; i++ {
		date := firstDay.AddDate(0, 0, i)
		data, msg := s.GetDailyStudyData(userID, date)
		if msg != "" {
			data = DailyStudyDataInfo{
				Date:      date,
				StudyTime: 0,
				Tomatoes:  0,
			}
		}
		result = append(result, data)
	}

	return result, ""
}

// 获取本年每月学习数据
func (s *StudyDataService) GetYearStudyData(userID uint, date time.Time) ([]MonthlyStudyDataInfo, string) {
	var result []MonthlyStudyDataInfo
	for i := 1; i <= 12; i++ {
		date := time.Date(date.Year(), time.Month(i), 1, 0, 0, 0, 0, date.Location())
		data, msg := s.GetMonthlyStudyData(userID, date)
		if msg != "" {
			data = MonthlyStudyDataInfo{
				Month:     date,
				StudyTime: 0,
				Tomatoes:  0,
			}
		}
		result = append(result, data)
	}
	return result, ""
}
