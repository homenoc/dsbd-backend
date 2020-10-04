package cmd

import (
	"github.com/homenoc/dsbd-backend/pkg/api"
	"github.com/homenoc/dsbd-backend/pkg/config"
	"github.com/spf13/cobra"
	"log"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start controller server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		confPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		if config.GetConfig(confPath) != nil {
			log.Fatalf("error config process |%v", err)
		}

		api.UserRestAPI()
		//server.Server()
		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringP("config", "c", "", "config path")
}
