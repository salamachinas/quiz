/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	quizv1 "github.com/salamachinas/quiz/pkg/pb/quiz/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//nolint:gochecknoglobals
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "quiz application client",
	Long:  `starts quiz application client`,
	Run:   play,
}

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(playCmd)

	playCmd.PersistentFlags().String("server", "localhost:5001", "server address connect to")
}

func play(cmd *cobra.Command, args []string) {
	serverAddress, err := cmd.Flags().GetString("server")
	if err != nil {
		log.Fatal("error getting sting server flag. Error: ", err)
	}

	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error dialing server. Error: ", err)
	}

	client := quizv1.NewQuizServiceClient(conn)

	// choose quiz
	quizID := chooseQuiz(client)

	// collect answers
	quizAnswers := collectAnswers(quizID, client)

	// Participate
	participantAnswers := []*quizv1.ParticipantAnswer{}
	for question, answer := range quizAnswers {
		participantAnswers = append(participantAnswers, &quizv1.ParticipantAnswer{
			QuestionId: question,
			AnswerId:   answer,
		})
	}

	result, err := client.Participate(context.Background(), &quizv1.ParticipateRequest{
		QuizId:             quizID,
		ParticipantId:      uuid.New().String(),
		ParticipantAnswers: participantAnswers,
	})
	if err != nil {
		log.Fatal("error participate in quiz. Error: ", err)
	}

	log.Printf(
		"Your score: %d%%. You were better than %d%% of all quizzers",
		result.ParticipantScore,
		result.ParticipantOverallScore,
	)
}

// return question sting id / answer sting id map for all questions in quiz.
func collectAnswers(quizID string, client quizv1.QuizServiceClient) map[string]string {
	result, err := client.Get(context.Background(), &quizv1.GetRequest{Id: quizID})
	if err != nil {
		log.Fatal("error getting quiz. Error: ", err)
	}

	questionAnswer := map[string]string{}
	for _, question := range result.Quiz.Questions {
		// map to human readable
		textQuestionAnswers := []string{}
		for _, answer := range question.Answers {
			textQuestionAnswers = append(textQuestionAnswers, answer.Text)
		}

		prompt := promptui.Select{
			Label: question.Text,
			Items: textQuestionAnswers,
		}

		_, textQuestionAnswer, err := prompt.Run()
		if err != nil {
			log.Fatal("error prompt. Error: ", err)
		}

		// map to machine readable
		for _, answer := range question.Answers {
			if answer.Text == textQuestionAnswer {
				questionAnswer[question.Id] = answer.Id
			}
		}
	}

	return questionAnswer
}

// chooseQuiz turn selected quiz sting ID.
func chooseQuiz(client quizv1.QuizServiceClient) string {
	result, err := client.List(context.Background(), &quizv1.ListRequest{})
	if err != nil {
		log.Fatal("error listing quizzes. Error: ", err)
	}

	chooseItems := []string{}
	for _, quiz := range result.Quizzes {
		chooseItems = append(chooseItems, quiz.Id)
	}

	prompt := promptui.Select{
		Label: "Choose quiz",
		Items: chooseItems,
	}

	_, chooseResult, err := prompt.Run()
	if err != nil {
		log.Fatal("error prompt. Error: ", err)
	}

	return chooseResult
}
