package main

import (
	"crypto/tls"

	"github.com/C3S/certinfo/internal/commands"
	"github.com/C3S/certinfo/internal/config"
	. "github.com/C3S/certinfo/internal/globals"
	"github.com/spf13/cobra"
)

func main() {

	confTLS := &tls.Config{
		InsecureSkipVerify: true,
	}

	// set default protocol to IPv6
	protocol := IPversions[1]

	rootCmd := &cobra.Command{Use: "certinfo"}

	cmdListHosts := &cobra.Command{
		Use:   "list",
		Short: "show all hosts defined",
		Long:  "lists all hosts in the configuration file.",
		Args:  cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return config.InitViperConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			commands.ListHosts(AllHosts)
		},
	}

	cmdExpires := &cobra.Command{
		Use:   "expires [server URL]",
		Short: "show when the certificate of [server URL] expires",
		Long:  "fetches the certificate and checks the expiration date.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			commands.Expires(args, protocol, confTLS)
		},
	}

	cmdBestBefore := &cobra.Command{
		Use:   "bestbefore [server URL] [days]",
		Short: "see if the certificate of [server URL] expires soon",
		Long:  "fetches the certificate and checks the difference between its expiration date and today. if that is shorter than [days], returns an error.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			commands.BestBefore(args, protocol, confTLS)
		},
	}

	cmdExpiresHosts := &cobra.Command{
		Use:   "hosts_expire",
		Short: "show when certificates of configured hosts expire",
		Long:  "fetches the certificates of hosts from the configuration file and checks the expiration date.",
		Args:  cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return config.InitViperConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			commands.ExpireHosts(confTLS)
		},
	}

	cmdBestBeforeHosts := &cobra.Command{
		Use:   "hosts_bestbefore",
		Short: "see if the certificates of configured hosts expire soon",
		Long:  "fetches the certificates of hosts in the configuration file and checks if their expiration date is shorter than [days].",
		Args:  cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return config.InitViperConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			commands.BestBeforeHosts(confTLS)
		},
	}

	rootCmd.PersistentFlags().IntVarP(&Days, "days", "d", Days, "number of days considered ok with bestbefore/hosts_bestbefore")
	rootCmd.PersistentFlags().IntVarP(&Port, "port", "p", Port, "the server port to use")
	rootCmd.PersistentFlags().IntVarP(&protocol, "protocol", "i", protocol, "the IP protocol version to use")
	rootCmd.PersistentFlags().IntVarP(&Timeout, "timeout", "t", Timeout, "the timout for dials in seconds")
	rootCmd.PersistentFlags().BoolVarP(&ShowErrors, "errors", "e", false, "show all dial errors, not only resolved hosts")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (defaults to ./config.yaml or ~/.config/certinfo/config.yaml")
	rootCmd.AddCommand(cmdListHosts)
	rootCmd.AddCommand(cmdExpires)
	rootCmd.AddCommand(cmdExpiresHosts)
	rootCmd.AddCommand(cmdBestBefore)
	rootCmd.AddCommand(cmdBestBeforeHosts)

	rootCmd.Execute()
}
