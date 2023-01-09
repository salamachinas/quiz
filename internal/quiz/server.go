package quiz

import (
	"context"

	quizv1 "github.com/salamachinas/quiz/pkg/pb/quiz/v1"
)

type Server struct {
	service *Service
	quizv1.UnimplementedQuizServiceServer
}

func NewServer(s *Service) *Server {
	return &Server{
		service: s,
	}
}

func (s *Server) List(ctx context.Context, _ *quizv1.ListRequest) (*quizv1.ListResponse, error) {
	quizzes, err := s.service.Find(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	return &quizv1.ListResponse{Quizzes: mapQuizzes(quizzes)}, nil
}

func (s *Server) Get(ctx context.Context, r *quizv1.GetRequest) (*quizv1.GetResponse, error) {
	quiz, err := s.service.Get(ctx, r.Id)
	if err != nil {
		return nil, mapError(err)
	}

	return &quizv1.GetResponse{Quiz: mapQuiz(quiz)}, nil
}

func (s *Server) Participate(ctx context.Context, r *quizv1.ParticipateRequest) (*quizv1.ParticipateResponse, error) {
	participateResult, err := s.service.Participate(ctx, mapParticipateRequest(r))
	if err != nil {
		return nil, mapError(err)
	}

	return &quizv1.ParticipateResponse{
		CorrectAnswers:          int32(participateResult.CorrectAnswers),
		IncorrectAnswers:        int32(participateResult.IncorrectAnswers),
		ParticipantCount:        int32(participateResult.ParticipantCount),
		ParticipantScore:        int32(participateResult.ParticipantScore),
		ParticipantOverallScore: int32(participateResult.ParticipantOverallScore),
	}, nil
}
