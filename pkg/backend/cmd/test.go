package cmd

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/tool/hash"
	"github.com/spf13/cobra"
	"log"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test controller server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("0:" + args[0])
		fmt.Println("1:" + args[1])
		fmt.Println(hash.Generate(args[0] + args[1]))
		//server.Server()
		log.Println("end")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.PersistentFlags().StringP("config", "c", "", "config path")
}
