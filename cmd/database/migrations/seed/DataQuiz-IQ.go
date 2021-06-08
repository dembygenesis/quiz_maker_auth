package seed

import (
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/string_utils"
)

func GetDataQuiz1() Quiz {
	var quiz Quiz

	// Quiz
	quiz.Name = "IQ Test"
	quiz.Slug = string_utils.Slugify("IQ Test")
	quiz.Order = 0
	quiz.Duration = 10 // 10 minutes only for "IQ Test"

	// Quiz Questions
	quiz.QuizQuestion = []QuizQuestion{
		{
			Name:       "Pick the correct sequence",
			Slug:       "Pick the correct sequence",
			Answer:     string_utils.Slugify("Pick the correct sequence"),
			// Quiz Choices
			QuizChoice: []QuizChoice{
				{
					Name:           "A",
					Slug:           string_utils.Slugify("A"),
					Order:          0,
				},
				{
					Name:           "B",
					Slug:           string_utils.Slugify("B"),
					Order:          1,
				},
				{
					Name:           "C",
					Slug:           string_utils.Slugify("C"),
					Order:          2,
				},
			},
		},
		{
			Name:       "Pick the correct sequence 2",
			Slug:       "Pick the correct sequence 2",
			Answer:     string_utils.Slugify("Pick the correct sequence 2"),
			// Quiz Choices
			QuizChoice: []QuizChoice{
				{
					Name:           "A",
					Slug:           string_utils.Slugify("A"),
					Order:          0,
				},
				{
					Name:           "B",
					Slug:           string_utils.Slugify("B"),
					Order:          1,
				},
				{
					Name:           "C",
					Slug:           string_utils.Slugify("C"),
					Order:          2,
				},
			},
		},
		{
			Name:       "Pick the correct sequence 3",
			Slug:       "Pick the correct sequence 3",
			Answer:     string_utils.Slugify("Pick the correct sequence 3"),
			// Quiz Choices
			QuizChoice: []QuizChoice{
				{
					Name:           "A",
					Slug:           string_utils.Slugify("A"),
					Order:          0,
				},
				{
					Name:           "B",
					Slug:           string_utils.Slugify("B"),
					Order:          1,
				},
				{
					Name:           "C",
					Slug:           string_utils.Slugify("C"),
					Order:          2,
				},
			},
		},
	}

	return quiz
}
