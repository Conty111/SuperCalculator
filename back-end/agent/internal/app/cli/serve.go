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

			if len(args) == 0 {
				log.Fatal().Msg("agent id argument required")
			}
			num, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatal().Err(err).Msg("Invalid index argument")
			}
			ctx = context.WithValue(ctx, "index", num)

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

	return command
}
