package cmd

import (
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/defektive/arsenic/arsenic/lib/util"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "arsenic",
	Short: "Arsenic - Pentest Conventions",
	Long: `Arsenic - Pentest Conventions


`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "the arsenic.yaml config file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	defaultDiscoverScripts := make(map[string]interface{})
	defaultReconScripts := make(map[string]interface{})
	defaultHuntScripts := make(map[string]interface{})
	defaultInitScripts := make(map[string]interface{})

	defaultInitScripts["as-init-op"] = util.NewScriptConfig("as-init-op", 0, true)
	defaultInitScripts["as-setup-hugo"] = util.NewScriptConfig("as-setup-hugo", 100, true)
	defaultInitScripts["as-init-hooks"] = util.NewScriptConfig("as-init-hooks", 200, true)
	defaultInitScripts["as-init-cleanup"] = util.NewScriptConfig("as-init-cleanup", 300, true)

	defaultDiscoverScripts["as-subdomain-discovery"] = util.NewScriptConfig("as-subdomain-discovery", 0, true)
	defaultDiscoverScripts["as-subdomain-enumeration"] = util.NewScriptConfig("as-subdomain-enumeration", 100, true)
	defaultDiscoverScripts["as-domains-from-domain-ssl-certs"] = util.NewScriptConfig("as-domains-from-domain-ssl-certs", 200, true)
	defaultDiscoverScripts["as-dns-resolution"] = util.NewScriptConfig("as-dns-resolution", 300, true)
	defaultDiscoverScripts["as-ip-recon"] = util.NewScriptConfig("as-ip-recon", 400, true)
	defaultDiscoverScripts["as-domains-from-ip-ssl-certs"] = util.NewScriptConfig("as-domains-from-ip-ssl-certs", 500, true)
	defaultDiscoverScripts["as-ip-resolution"] = util.NewScriptConfig("as-ip-resolution", 600, true)
	defaultDiscoverScripts["as-http-screenshot-domains"] = util.NewScriptConfig("as-http-screenshot-domains", 700, true)

	defaultReconScripts["as-port-scan-tcp"] = util.NewScriptConfig("as-port-scan-tcp", 0, true)
	defaultReconScripts["as-content-discovery"] = util.NewScriptConfig("as-content-discovery", 100, true)
	defaultReconScripts["as-http-screenshot-hosts"] = util.NewScriptConfig("as-http-screenshot-hosts", 200, true)
	defaultReconScripts["as-port-scan-udp"] = util.NewScriptConfig("as-port-scan-udp", 300, true)

	defaultHuntScripts["as-takeover-aquatone"] = util.NewScriptConfig("as-takeover-aquatone", 0, true)
	defaultHuntScripts["as-searchsploit"] = util.NewScriptConfig("as-searchsploit", 100, true)

	defaultScripts := make(map[string]interface{})
	defaultScripts["init"] = defaultInitScripts
	defaultScripts["discover"] = defaultDiscoverScripts
	defaultScripts["recon"] = defaultReconScripts
	defaultScripts["hunt"] = defaultHuntScripts

	wordlists := make(map[string][]string)
	wordlists["web-content"] = []string{
		"Discovery/Web-Content/AdobeCQ-AEM.txt",
		"Discovery/Web-Content/apache.txt",
		"Discovery/Web-Content/Common-DB-Backups.txt",
		"Discovery/Web-Content/Common-PHP-Filenames.txt",
		"Discovery/Web-Content/common.txt",
		"Discovery/Web-Content/confluence-administration.txt",
		"Discovery/Web-Content/default-web-root-directory-linux.txt",
		"Discovery/Web-Content/default-web-root-directory-windows.txt",
		"Discovery/Web-Content/frontpage.txt",
		"Discovery/Web-Content/graphql.txt",
		"Discovery/Web-Content/jboss.txt",
		"Discovery/Web-Content/Jenkins-Hudson.txt",
		"Discovery/Web-Content/nginx.txt",
		"Discovery/Web-Content/oracle.txt",
		"Discovery/Web-Content/quickhits.txt",
		"Discovery/Web-Content/raft-large-directories.txt",
		"Discovery/Web-Content/raft-medium-words.txt",
		"Discovery/Web-Content/reverse-proxy-inconsistencies.txt",
		"Discovery/Web-Content/RobotsDisallowed-Top1000.txt",
		"Discovery/Web-Content/websphere.txt",
	}

	setConfigDefault("scripts", defaultScripts)
	setConfigDefault("wordlists", wordlists)
	setConfigDefault("wordlist-paths", []string{"/opt/SecLists"})

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in home directory with name "arsenic" (without extension).
		viper.AddConfigPath(cwd)
		viper.AddConfigPath(home)
		viper.SetConfigName("arsenic")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.ReadInConfig()
}

// If no config file exists, all possible keys in the defaults
// need to be registered with viper otherwise viper will only think
// scripts and sec-lists-path are valid keys
func setConfigDefault(key string, value interface{}) {
	if valueMap, ok := value.(map[string]interface{}); ok {
		for k, v := range valueMap {
			setConfigDefault(fmt.Sprintf("%s.%s", key, k), v)
		}
	} else if mappable, ok := value.(util.Mappable); ok {
		valueMap := mappable.ToMap()
		for k, v := range valueMap {
			setConfigDefault(fmt.Sprintf("%s.%s", key, k), v)
		}
	} else {
		viper.SetDefault(key, value)
	}
}
