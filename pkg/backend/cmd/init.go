package cmd

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	initFunction "github.com/homenoc/dsbd-backend/pkg/api/core/tool/init"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"github.com/spf13/cobra"
	"log"
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

		store.InitDB()

		log.Println("end")
	},
}

var initRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "register database table",
	Long:  `register database table`,
	Run: func(cmd *cobra.Command, args []string) {
		confPath, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		template, err := cmd.Flags().GetString("template")
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		if config.GetConfig(confPath) != nil {
			log.Fatalf("error config process |%v", err)
		}

		if err = initFunction.RegisterTemplateConfig(template); err != nil {
			log.Fatalf("error config process |%v", err)
		}

		//api.AdminRestAPI()
		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.AddCommand(initDatabaseCmd)
	initCmd.AddCommand(initRegisterCmd)
	initCmd.PersistentFlags().StringP("config", "c", "", "config path")
	initCmd.PersistentFlags().StringP("template", "t", "", "template path")
}
