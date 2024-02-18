package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// NewServeCmd starts new application instance
func NewServeCmd() *cobra.Command {
	var local bool
	var count_agents uint
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

			ctx = context.WithValue(ctx, "local", local)
			ctx = context.WithValue(ctx, "agents_count", count_agents)

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
	command.Flags().BoolVar(&local, "local", false, "runs orchestrator with local agent addresses")
	command.Flags().UintVar(&count_agents, "count_agents", 2, "set count of agents for local running")
	return command
}
