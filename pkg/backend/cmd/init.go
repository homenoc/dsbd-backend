package cmd

import (
	"github.com/homenoc/dsbd-backend/pkg/config"
	"github.com/homenoc/dsbd-backend/pkg/store"
	"github.com/spf13/cobra"
	"log"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init database",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		confPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		if config.GetConfig(confPath) != nil {
			log.Fatalf("error config process |%v", err)
		}

		if err := store.InitDB(); err != nil {
			log.Println("Success!!")
		} else {
			log.Println("error: " + err.Error())
		}

		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
