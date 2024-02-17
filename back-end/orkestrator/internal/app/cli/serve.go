package cli

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// NewServeCmd starts new application instance
func NewServeCmd() *cobra.Command {
	command := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"s"},
		Short:   "Start server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("Starting")

			sigchan := make(chan os.Signal, 1)
			signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			l, err := cmd.Flags().GetString("local")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to get local flag")
			}
			ctx = context.WithValue(ctx, "local", l)
			http_port, err := cmd.Flags().GetUint("http_base_port")
			if err != nil {
				log.Fatal().Err(err).Msg("failed to get http_base_port flags value")
			}
			ctx = context.WithValue(ctx, "http_base_port", strconv.Itoa(int(http_port)))

			application, err := app.InitializeApplication(ctx)

			if err != nil {
				log.Fatal().Err(err).Msg("can not initialize application")
			}

			cliMode := false
			application.Start(ctx, cliMode)

			log.Info().Msg("Started")
			<-sigchan
			ctx.Done()
			log.Error().Err(application.Stop()).Msg("stop application")

			time.Sleep(time.Second * cliCmdExecFinishDelaySeconds)
			log.Info().Msg("Finished")
		},
	}
	command.Flags().String("local", "", "runs orchestrator with local addresses")
	command.Flags().Uint("http_base_port", 8000, "http base server port for local running")
	return command
}
