package cmd

import (
	"github.com/homenoc/dsbd-backend/pkg/api"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/slack"
	"github.com/spf13/cobra"
	"github.com/stripe/stripe-go/v73"
	"log"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start controller server",
	Long:  ``,
}

var startUserCmd = &cobra.Command{
	Use:   "user",
	Short: "start user mode",
	Long:  `start user mode`,
	Run: func(cmd *cobra.Command, args []string) {
		confPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		config.IsDebug, err = cmd.Flags().GetBool("debug")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		if config.GetConfig(confPath) != nil {
			log.Fatalf("error config process |%v", err)
		}
		notification.NewNotification()
		notification.NoticeLog("good", []string{
			"Status: User側 API起動",
		})
		stripe.Key = config.Conf.Stripe.SecretKey

		api.UserRestAPI()
		log.Println("end")
	},
}

var startAdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "start admin mode",
	Long:  `start admin mode`,
	Run: func(cmd *cobra.Command, args []string) {
		confPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		config.IsDebug, err = cmd.Flags().GetBool("debug")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		if config.GetConfig(confPath) != nil {
			log.Fatalf("error config process |%v", err)
		}

		notification.NewNotification()
		go slack.StartAppSlack()
		notification.NoticeLog("good", []string{
			"Status: Admin側 API起動",
		})
		stripe.Key = config.Conf.Stripe.SecretKey

		api.AdminRestAPI()
		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.AddCommand(startAdminCmd)
	startCmd.AddCommand(startUserCmd)
	startCmd.PersistentFlags().StringP("config", "c", "", "config path")
	startCmd.PersistentFlags().BoolP("debug", "d", false, "debug command")
}
