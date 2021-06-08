package seed

type QuizChoice struct {
	QuizQuestionID int
	Name           string
	Slug           string
	Order          int
}

type QuizQuestion struct {
	QuizID int
	Name   string
	Slug   string
	Answer string
	Order  int

	// `gorm:"-"` <-------- Ignores
	QuizChoice []QuizChoice `gorm:"-"`
}

type Quiz struct {
	ID           int
	Name         string
	Slug         string
	Order        int
	Duration     int
	QuizQuestion []QuizQuestion
}

type QuestionAndAnswer struct {
	Name   string
	Answer string
}
