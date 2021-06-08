package modelsV2

import "gorm.io/gorm"

type Quiz struct {
	gorm.Model    `json:"-" gorm:"-"`
	ID            int            `json:"id"`
	Name          string         `json:"name" db:"name" gorm:"uniqueIndex:uniq_name;type:varchar(255);"`
	Slug          string         `json:"slug" db:"slug" gorm:"uniqueIndex:uniq_slug;type:varchar(255);"`
	Order         int            `json:"order" db:"order" gorm:"uniqueIndex:uniq_order;type:varchar(255);"`
	Duration      int            `json:"duration"`
	QuizQuestions []QuizQuestion `json:"quiz_questions"`
}
