package quiz

import (
	"context"
	"sync"

	"github.com/go-faker/faker/v3"
)

type MemoryStorage struct {
	quizzes             map[string]Quiz
	ParticipantCount    int32
	ParticipantsResults []ParticipateResult
	mutex               *sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		mutex:               &sync.RWMutex{},
		ParticipantCount:    0,
		ParticipantsResults: []ParticipateResult{},
		quizzes:             generateQuizzes(),
	}
}

func (s *MemoryStorage) Find(_ context.Context) ([]Quiz, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := []Quiz{}

	for _, q := range s.quizzes {
		result = append(result, q)
	}

	return result, nil
}

func (s *MemoryStorage) FindOne(_ context.Context, id string) (Quiz, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	quiz, ok := s.quizzes[id]
	if !ok {
		return Quiz{}, ErrQuizNotFound
	}

	return quiz, nil
}

func (s *MemoryStorage) PushParticipantsResults(_ context.Context, id string, p ParticipateResult) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	quiz, ok := s.quizzes[id]
	if !ok {
		return ErrQuizNotFound
	}

	quiz.ParticipantsResults = append(quiz.ParticipantsResults, p)

	s.quizzes[id] = quiz

	return nil
}

func (s *MemoryStorage) IncParticipantCount(_ context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	quiz, ok := s.quizzes[id]
	if !ok {
		return ErrQuizNotFound
	}

	quiz.ParticipantCount++

	s.quizzes[id] = quiz

	return nil
}

func generateQuizzes() map[string]Quiz {
	quizzes := map[string]Quiz{}

	for i := 0; i < 3; i++ {
		quiz := Quiz{
			ID:                  faker.UUIDDigit(),
			Text:                faker.Sentence(),
			ParticipantCount:    0,
			Questions:           generateQuestions(),
			ParticipantsResults: []ParticipateResult{},
		}

		quizzes[quiz.ID] = quiz
	}

	return quizzes
}

func generateQuestions() map[string]Question {
	questions := map[string]Question{}

	for i := 0; i < 5; i++ {
		answers := generateAnswers()

		question := Question{
			ID:      faker.UUIDDigit(),
			Text:    faker.Sentence(),
			Answers: answers,
		}

		// finger cross hope for go random order.
		for _, answer := range answers {
			question.CorrectAnswerID = answer.ID
			break
		}

		questions[question.ID] = question
	}

	return questions
}

func generateAnswers() map[string]Answer {
	answers := map[string]Answer{}

	for i := 0; i < 3; i++ {
		answer := Answer{
			ID:   faker.UUIDDigit(),
			Text: faker.Word(),
		}

		answers[answer.ID] = answer
	}

	return answers
}
