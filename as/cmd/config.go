package cmd

import (
	"fmt"

	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display config information",
	Long: `Display config information.

Helpful to see what scripts would be executed.`,
	Run: func(cmd *cobra.Command, args []string) {
		writeCfg, _ := cmd.Flags().GetBool("write")

		t, err := toml.TreeFromMap(viper.AllSettings())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Configuration")
		s := t.String()
		fmt.Printf(s)


		fmt.Println("Discovery files to be run")
		for _, scriptFile := range getScripts("discovery") {
			fmt.Println(scriptFile)
		}

		fmt.Println("Recon files to be run")
		for _, scriptFile := range getScripts("recon") {
			fmt.Println(scriptFile)
		}


		if writeCfg {
			fmt.Println("Writing Config")
			viper.WriteConfig()
		}

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolP("write", "w", false, "write config")
}