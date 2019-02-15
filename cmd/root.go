package cmd

import (
	"fmt"
	"os"

	"github.com/denouche/plex-watcher/handlers"
	"github.com/denouche/plex-watcher/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config  = &handlers.Config{}
	cfgFile string
)

const (
	parameterConfigurationFile = "config"
	parameterLogLevel          = "loglevel"
	parameterLogFormat         = "logformat"
	parameterPort              = "port"
	parameterSlackURI          = "slackuri"
	parameterDBInMemoryFile    = "dbinmemoryfile"
)

var (
	defaultLogLevel       = logrus.InfoLevel.String()
	defaultLogFormat      = utils.LogFormatText
	defaultPort           = 80
	defaultSlackURI       = ""
	defaultDBInMemoryFile = ""
)

var rootCmd = &cobra.Command{
	Use:   "plex-watcher",
	Short: "plex-watcher",
	Run: func(cmd *cobra.Command, args []string) {
		utils.InitLogger(config.LogLevel, config.LogFormat)

		logrus.
			WithField(parameterConfigurationFile, cfgFile).
			WithField(parameterDBInMemoryFile, config.DBInMemoryFile).
			WithField(parameterLogLevel, config.LogLevel).
			WithField(parameterLogFormat, config.LogFormat).
			WithField(parameterPort, config.Port).
			Warn("Configuration")

		router := handlers.NewRouter(config)
		router.Run(fmt.Sprintf(":%d", config.Port))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, parameterConfigurationFile, "", "Config file. All flags given in command line will override the values from this file.")

	rootCmd.Flags().String(parameterLogLevel, defaultLogLevel, "Use this flag to set the logging level")
	viper.BindPFlag(parameterLogLevel, rootCmd.Flags().Lookup(parameterLogLevel))

	rootCmd.Flags().String(parameterLogFormat, defaultLogFormat, "Use this flag to set the logging format")
	viper.BindPFlag(parameterLogFormat, rootCmd.Flags().Lookup(parameterLogFormat))

	rootCmd.Flags().Int(parameterPort, defaultPort, "Use this flag to set the listening port of the api")
	viper.BindPFlag(parameterPort, rootCmd.Flags().Lookup(parameterPort))

	rootCmd.Flags().String(parameterSlackURI, defaultSlackURI, "Use this flag to set the slack notification URI")
	viper.BindPFlag(parameterSlackURI, rootCmd.Flags().Lookup(parameterSlackURI))

	rootCmd.Flags().String(parameterDBInMemoryFile, defaultDBInMemoryFile, "Use this flag to set file to load and save db in memory datas")
	viper.BindPFlag(parameterDBInMemoryFile, rootCmd.Flags().Lookup(parameterDBInMemoryFile))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	config.LogLevel = viper.GetString(parameterLogLevel)
	config.LogFormat = viper.GetString(parameterLogFormat)
	config.Port = viper.GetInt(parameterPort)
	config.SlackURI = viper.GetString(parameterSlackURI)
	config.DBInMemoryFile = viper.GetString(parameterDBInMemoryFile)
}
