package cmd

import (
	"log"

	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	initFunction "github.com/homenoc/dsbd-backend/pkg/api/core/tool/init"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/seed"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init command",
	Long:  ``,
}

var initDatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: "init database",
	Long:  `init database`,
	Run: func(cmd *cobra.Command, args []string) {
		confPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		if config.GetConfig(confPath) != nil {
			log.Fatalf("error config process |%v", err)
		}
		notification.NewNotification()

		store.InitDB()

		log.Println("[Init] Database initialization completed")
	},
}

var initRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "register database table",
	Long:  `register database table`,
	Run: func(cmd *cobra.Command, args []string) {
		confPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("[Config] Failed to get config path: %v", err)
		}
		template, err := cmd.Flags().GetString("template")
		if err != nil {
			log.Fatalf("[Config] Failed to get template path: %v", err)
		}

		if config.GetConfig(confPath) != nil {
			log.Fatalf("[Config] Failed to load config: %v", err)
		}
		notification.NewNotification()

		if err = initFunction.RegisterTemplateConfig(template); err != nil {
			log.Fatalf("[Register] Failed to register template: %v", err)
		}

		log.Println("[Register] Template registration completed")
	},
}

var initSeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "create seed data for development",
	Long:  `create seed data for development (test users, groups, NOC, etc.)`,
	Run: func(cmd *cobra.Command, args []string) {
		confPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("[Config] Failed to get config path: %v", err)
		}

		log.Printf("[Config] Loading config from %s", confPath)
		if config.GetConfig(confPath) != nil {
			log.Fatalf("[Config] Failed to load config: %v", err)
		}
		log.Println("[Config] Config loaded successfully")

		notification.NewNotification()

		if err := seed.Run(); err != nil {
			log.Fatalf("[Seed] Failed to create seed data: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.AddCommand(initDatabaseCmd)
	initCmd.AddCommand(initRegisterCmd)
	initCmd.AddCommand(initSeedCmd)
	initCmd.PersistentFlags().StringP("config", "c", "", "config path")
	initCmd.PersistentFlags().StringP("template", "t", "", "template path")
}
