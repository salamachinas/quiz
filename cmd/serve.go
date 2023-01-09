package cmd

import (
	"log"
	"net"

	"github.com/salamachinas/quiz/internal/quiz"
	quizv1 "github.com/salamachinas/quiz/pkg/pb/quiz/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

//nolint:gochecknoglobals
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "quiz application server",
	Long:  `starts quiz application server on specific port`,
	Run:   serve,
}

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().String("address", "localhost:5001", "address listen on")
}

func serve(cmd *cobra.Command, args []string) {
	address, err := cmd.Flags().GetString("address")
	if err != nil {
		log.Fatal("error getting sting address flag. Error: ", err)
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("error listing. Error: ", err)
	}

	// wire
	quizStorage := quiz.NewMemoryStorage()
	quizUsecase := quiz.NewService(quizStorage)
	quizServer := quiz.NewServer(quizUsecase)
	grpcServer := grpc.NewServer()

	quizv1.RegisterQuizServiceServer(grpcServer, quizServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("error serving. Error: ", err)
	}
}
