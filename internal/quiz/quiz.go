package quiz

import "errors"

var (
	ErrQuizNotFound    = errors.New("QuizNotFound")
	ErrInvalidArgument = errors.New("InvalidRequest")
)

type (
	Quiz struct {
		ID                  string
		Text                string
		Questions           map[string]Question
		ParticipantCount    int
		ParticipantsResults []ParticipateResult
	}

	Question struct {
		ID              string
		Text            string
		CorrectAnswerID string
		Answers         map[string]Answer
	}

	Answer struct {
		ID   string
		Text string
	}
)

type (
	ParticipantAnswer struct {
		QuestionID string
		AnswerID   string
	}

	ParticipateArgs struct {
		QuizID             string
		ParticipantID      string
		ParticipantAnswers []ParticipantAnswer
	}

	ParticipateResult struct {
		ParticipantID           string
		CorrectAnswers          int
		IncorrectAnswers        int
		ParticipantCount        int
		ParticipantScore        int
		ParticipantOverallScore int
	}
)
