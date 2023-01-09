package quiz

import (
	"errors"

	quizv1 "github.com/salamachinas/quiz/pkg/pb/quiz/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapQuizzes(q []Quiz) []*quizv1.Quiz {
	result := []*quizv1.Quiz{}
	for _, quiz := range q {
		result = append(result, mapQuiz(quiz))
	}

	return result
}

func mapQuiz(q Quiz) *quizv1.Quiz {
	result := &quizv1.Quiz{
		Id:               q.ID,
		Text:             q.Text,
		ParticipantCount: int32(q.ParticipantCount),
		Questions:        make(map[string]*quizv1.Question),
	}

	for _, question := range q.Questions {
		result.Questions[question.ID] = mapQuestion(question)
	}

	return result
}

func mapParticipateRequest(r *quizv1.ParticipateRequest) ParticipateArgs {
	result := ParticipateArgs{
		QuizID:        r.QuizId,
		ParticipantID: r.ParticipantId,
	}

	for _, answer := range r.ParticipantAnswers {
		result.ParticipantAnswers = append(result.ParticipantAnswers, ParticipantAnswer{
			QuestionID: answer.QuestionId,
			AnswerID:   answer.AnswerId,
		})
	}

	return result
}

func mapQuestion(q Question) *quizv1.Question {
	result := &quizv1.Question{
		Id:      q.ID,
		Text:    q.Text,
		Answers: make(map[string]*quizv1.Answer),
	}

	for _, answer := range q.Answers {
		result.Answers[answer.ID] = mapAnswer(answer)
	}

	return result
}

func mapAnswer(a Answer) *quizv1.Answer {
	return &quizv1.Answer{
		Id:   a.ID,
		Text: a.Text,
	}
}

func mapError(e error) error {
	if errors.Is(e, ErrQuizNotFound) {
		return status.Errorf(codes.NotFound, "quiz not found")
	}

	if errors.Is(e, ErrInvalidArgument) {
		return status.Errorf(codes.InvalidArgument, "provided arguments are not valid")
	}

	return status.Errorf(codes.Internal, "internal server error, please contact support")
}
