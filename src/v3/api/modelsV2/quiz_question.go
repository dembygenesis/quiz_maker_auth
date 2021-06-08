package modelsV2

import "gorm.io/gorm"

type QuizQuestion struct {
	gorm.Model  `json:"-" gorm:"-"`
	ID          uint         `json:"id" db:"id"`
	QuizId      int          `json:"quiz_id"`
	Quiz        Quiz         `json:"-"`
	Name        string       `json:"name" db:"name" gorm:"uniqueIndex:uniq_name;type:varchar(255);"`
	Slug        string       `json:"slug" db:"slug" gorm:"uniqueIndex:uniq_slug;type:varchar(255);"`
	Answer      string       `json:"answer" db:"answer" gorm:"uniqueIndex:uniq_slug;type:varchar(255);"`
	Order       int          `json:"order" db:"order" gorm:"uniqueIndex:uniq_order;type:varchar(255);"`
	QuizChoices []QuizChoice `json:"quiz_choices"`
}
