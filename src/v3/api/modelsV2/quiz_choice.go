package modelsV2

type QuizChoice struct {
	ID             int          `json:"id"`
	QuizQuestionID int          `json:"quiz_question_id" gorm:"uniqueIndex:uniq_composite;type:varchar(255);"`
	QuizQuestion   QuizQuestion `json:"-"`
	Name           string       `json:"name" gorm:"uniqueIndex:uniq_composite;type:varchar(255);"`
	Slug           string       `json:"slug" gorm:"uniqueIndex:uniq_composite;type:varchar(255);"`
	Order          int          `json:"order" gorm:"uniqueIndex:uniq_composite;type:varchar(255);"`
}
