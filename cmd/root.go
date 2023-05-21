/*
Copyright © 2023 Nextlinux <dev@next-linux.systems>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/nextlinux/fpmanagement/fpmanagement"
	"github.com/nextlinux/fpmanagement/internal/client"
	"github.com/nextlinux/fpmanagement/internal/config"
	"github.com/spf13/cobra"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fpmanagement",
	Short: "Loads a bunch of package corrections into Nextlinux",
	Run: func(cmd *cobra.Command, args []string) {
		fpmanagement.AddCorrections(enterpriseAPIClient)
	},
}

var appConfig *config.Application
var enterpriseAPIClient *client.EnterpriseClient

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig,
		initClient)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfg, err := config.LoadApplicationConfig(viper.GetViper())
	if err != nil {
		if _, printErr := fmt.Fprintln(os.Stderr, fmt.Errorf("unable to load application config: %w", err)); printErr != nil {
			panic(err)
		}
		os.Exit(1)
	}
	appConfig = cfg
	fmt.Println(appConfig.String())
}

func initClient() {
	var scheme string
	var hostname = appConfig.Nextlinux.URL
	urlFields := strings.Split(hostname, "://")

	if len(urlFields) > 1 {
		scheme = urlFields[0]
		hostname = urlFields[1]
	}
	enterpriseClient, err := client.NewEnterpriseClient(client.Configuration{
		Hostname:       hostname,
		Username:       appConfig.Nextlinux.User,
		Password:       appConfig.Nextlinux.Password,
		Scheme:         scheme,
		TimeoutSeconds: appConfig.Nextlinux.HTTP.TimeoutSeconds,
		Insecure:       appConfig.Nextlinux.HTTP.Insecure,
		NextlinuxAccount: appConfig.Nextlinux.Account,
	})
	if err != nil {
		fmt.Printf("failed to initialize enterprise client: %s\n", err.Error())
		os.Exit(1)
	}
	enterpriseAPIClient = enterpriseClient
}
