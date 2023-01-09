package quiz

import (
	"context"
)

type storage interface {
	Find(ctx context.Context) ([]Quiz, error)
	FindOne(ctx context.Context, id string) (Quiz, error)
	PushParticipantsResults(ctx context.Context, id string, p ParticipateResult) error
	IncParticipantCount(ctx context.Context, id string) error
}

type Service struct {
	storage storage
}

func NewService(s storage) *Service {
	return &Service{
		storage: s,
	}
}

func (s *Service) Find(ctx context.Context) ([]Quiz, error) {
	return s.storage.Find(ctx)
}

func (s *Service) Get(ctx context.Context, id string) (Quiz, error) {
	return s.storage.FindOne(ctx, id)
}

func (s *Service) Participate(ctx context.Context, args ParticipateArgs) (ParticipateResult, error) {
	// load quiz
	quiz, err := s.storage.FindOne(ctx, args.QuizID)
	if err != nil {
		return ParticipateResult{}, err
	}

	// validate args
	if len(args.ParticipantAnswers) != len(quiz.Questions) {
		return ParticipateResult{}, ErrInvalidArgument
	}

	// participate
	correctAnswers, incurredAnswers := checkAnswers(quiz, args.ParticipantAnswers)

	// overall score
	scored := 0
	for _, result := range quiz.ParticipantsResults {
		if result.CorrectAnswers < correctAnswers {
			scored++
		}
	}

	// store participate participateResult
	participateResult := ParticipateResult{
		ParticipantID:           args.ParticipantID,
		CorrectAnswers:          correctAnswers,
		IncorrectAnswers:        incurredAnswers,
		ParticipantCount:        quiz.ParticipantCount + 1,
		ParticipantScore:        calculateScore(correctAnswers, len(quiz.Questions)),
		ParticipantOverallScore: calculateScore(scored, len(quiz.ParticipantsResults)),
	}

	if err := s.storage.PushParticipantsResults(ctx, args.QuizID, participateResult); err != nil {
		return ParticipateResult{}, err
	}

	// inc participant count
	if err := s.storage.IncParticipantCount(ctx, args.QuizID); err != nil {
		return ParticipateResult{}, err
	}

	return participateResult, nil
}

// checkAnswers check quiz answers and return correct answers count int, incurred answers count int.
func checkAnswers(q Quiz, a []ParticipantAnswer) (int, int) {
	correctAnswers := 0
	incurredAnswers := 0

	for _, answer := range a {
		if answer.AnswerID == q.Questions[answer.QuestionID].CorrectAnswerID {
			correctAnswers++
		} else {
			incurredAnswers++
		}
	}

	return correctAnswers, incurredAnswers
}

func calculateScore(scored, total int) int {
	score := float64(scored) / float64(total) * 100 //nolint:gomnd
	return int(score)
}
