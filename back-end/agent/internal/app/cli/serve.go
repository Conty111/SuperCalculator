package cli

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Conty111/SuperCalculator/back-end/agent/internal/app"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// NewServeCmd starts new application instance
func NewServeCmd() *cobra.Command {
	command := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"s"},
		Short:   "starts agent",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("Starting")

			sigchan := make(chan os.Signal, 1)
			signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			http_port, err := cmd.Flags().GetUint("http_port")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to get http_port flags value")
			}
			ctx = context.WithValue(ctx, "http_port", strconv.Itoa(int(http_port)))
			num, err := cmd.Flags().GetUint("agent_id")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to get http_port flags value")
			}
			ctx = context.WithValue(ctx, "partition", int32(num))

			application, err := app.InitializeApplication(ctx)

			if err != nil {
				log.Fatal().Err(err).Msg("can not initialize application")
			}

			cliMode := false
			application.Start(ctx, cliMode)

			log.Info().Msg("Started")
			<-sigchan

			log.Error().Err(application.Stop()).Msg("stop application")

			time.Sleep(time.Second * cliCmdExecFinishDelaySeconds)
			log.Info().Msg("Finished")
		},
	}
	command.PersistentFlags().Uint("agent_id", 0, "equal to number of partition which should be used by consumer")
	command.PersistentFlags().Uint("http_port", 8000, "http server port")
	return command
}
