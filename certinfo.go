package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Host struct {
	URL  string `mapstructure:"url"`
	Port int    `mapstructure:"port"`
}

var (
	allHosts map[string]*Host
)

func dialTimeout(sec int) time.Duration {
	return time.Duration(sec) * time.Second
}

func initViperConfig(cmd *cobra.Command) error {
	v := viper.New()

	customConfig := cmd.Flags().Lookup("config").Value.String()

	if customConfig != "" {
		v.SetConfigFile(customConfig)
	} else {
		v.SetConfigName("config.yaml")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.config/certinfo/")
	}
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Panicf("Fatal error in config file: %s\n", err)
		} else {
		}
		log.Print("Using fallback configuration (no config file found)")
	} else {
		err = v.UnmarshalKey("hosts", &allHosts)
	}

	return nil
}

func sortKeys(hosts map[string]*Host) []string {
	keys := make([]string, 0, len(hosts))

	for k := range hosts {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func main() {

	confTLS := &tls.Config{
		InsecureSkipVerify: true,
	}

	days := 14
	port := 443
	IPversions := [2]int{4, 6}
	// set default protocol to IPv6
	protocol := IPversions[1]
	timeout := 5
	showErrors := true

	// colorful output
	blue := color.New(color.FgBlue).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	orange := color.New(color.FgYellow).SprintFunc()
	// yellow := color.New(color.FgHiYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	rootCmd := &cobra.Command{Use: "certinfo"}

	cmdListHosts := &cobra.Command{
		Use:   "list",
		Short: "show all hosts defined",
		Long:  "lists all hosts in the configuration file.",
		Args:  cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initViperConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			keys := sortKeys(allHosts)
			for _, i := range keys {
				fmt.Printf("%s: %s\n", blue(allHosts[i].URL), orange(strconv.Itoa(allHosts[i].Port)))
			}
		},
	}

	cmdExpires := &cobra.Command{
		Use:   "expires [server URL]",
		Short: "show when the certificate of [server URL] expires",
		Long:  "fetches the certificate and checks the expiration date.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := tls.DialWithDialer(
				&net.Dialer{Timeout: dialTimeout(timeout)},
				"tcp"+strconv.Itoa(protocol),
				args[0]+":"+strconv.Itoa(port),
				confTLS,
			)
			if err != nil {
				if showErrors {
					log.Printf("%s: %s", blue(args[0]+":"+strconv.Itoa(port)), orange(err))
				}
				return
			}
			defer conn.Close()
			certs := conn.ConnectionState().PeerCertificates
			fmt.Printf("expires: %s (IPv%d)\n", magenta(certs[0].NotAfter.Format("02.01.2006")), protocol)
		},
	}

	cmdBestBefore := &cobra.Command{
		Use:   "bestbefore [server URL] [days]",
		Short: "see if the certificate of [server URL] expires soon",
		Long:  "fetches the certificate and checks the difference between its expiration date and today. if that is shorter than [days], returns an error.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			loc, _ := time.LoadLocation("UTC")
			now := time.Now().In(loc)
			conn, err := tls.DialWithDialer(
				&net.Dialer{Timeout: dialTimeout(timeout)},
				"tcp"+strconv.Itoa(protocol),
				args[0]+":"+strconv.Itoa(port),
				confTLS,
			)
			if err != nil {
				if showErrors {
					log.Printf(
						"%s: %s %s",
						blue(args[0]+":"+strconv.Itoa(port)),
						red("Error during dial:"),
						orange(err),
					)
				}
				return
			}
			defer conn.Close()
			certs := conn.ConnectionState().PeerCertificates
			certExpires := certs[0].NotAfter
			daysValid := int(certExpires.Sub(now).Hours() / 24)
			if daysValid > days {
				fmt.Printf(
					"will expire in %s days (IPv%d) %s",
					green(strconv.Itoa(daysValid)),
					protocol,
					green(" -- ok!\n"),
				)
				return
			} else if daysValid < 0 {
				fmt.Printf(
					"expired %s days ago (IPv%d) %s",
					red(strconv.Itoa(daysValid)),
					protocol,
					red(" -- red alert!\n"),
				)
				os.Exit(1)
			} else {
				fmt.Printf(
					"expires in %s days (IPv%d) %s",
					orange(strconv.Itoa(daysValid)),
					protocol,
					orange(" -- please renew!\n"),
				)
				os.Exit(1)
			}
		},
	}

	cmdExpiresHosts := &cobra.Command{
		Use:   "hosts_expire",
		Short: "show when certificates of configured hosts expire",
		Long:  "fetches the certificates of hosts from the configuration file and checks the expiration date.",
		Args:  cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initViperConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			keys := sortKeys(allHosts)
			for _, i := range keys {
				for _, j := range IPversions {
					conn, err := tls.DialWithDialer(
						&net.Dialer{Timeout: dialTimeout(timeout)},
						"tcp"+strconv.Itoa(j),
						allHosts[i].URL+":"+strconv.Itoa(allHosts[i].Port),
						confTLS,
					)
					if err != nil {
						if showErrors {
							log.Printf(
								"%-35s %s",
								blue(allHosts[i].URL+":"),
								orange(err),
							)
						}
						continue
					}
					defer conn.Close()
					certs := conn.ConnectionState().PeerCertificates
					fmt.Printf(
						"%-35s expires: %s (IPv%d)\n",
						blue(allHosts[i].URL),
						magenta(certs[0].NotAfter.Format("02.01.2006")),
						j,
					)
				}
			}
		},
	}

	cmdBestBeforeHosts := &cobra.Command{
		Use:   "hosts_bestbefore",
		Short: "see if the certificates of configured hosts expire soon",
		Long:  "fetches the certificates of hosts in the configuration file and checks if their expiration date is shorter than [days].",
		Args:  cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initViperConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			loc, _ := time.LoadLocation("UTC")
			now := time.Now().In(loc)
			keys := sortKeys(allHosts)
			for _, i := range keys {
				for _, j := range IPversions {
					conn, err := tls.DialWithDialer(
						&net.Dialer{Timeout: dialTimeout(timeout)},
						"tcp"+strconv.Itoa(j),
						allHosts[i].URL+":"+strconv.Itoa(allHosts[i].Port),
						confTLS,
					)
					if err != nil {
						if showErrors {
							log.Printf(
								"%-35s %s: %s %s",
								blue(allHosts[i].URL),
								"(IPv"+strconv.Itoa(j)+")",
								red("Error during dial:"),
								orange(err),
							)
						}
						continue
					}
					defer conn.Close()
					certs := conn.ConnectionState().PeerCertificates
					certExpires := certs[0].NotAfter
					daysValid := int(certExpires.Sub(now).Hours() / 24)
					if daysValid > days {
						fmt.Printf(
							"%-35s (IPv%d): expires %-44s %s",
							blue(allHosts[i].URL),
							j,
							green(certs[0].NotAfter.Format("02.01.2006"))+", in "+green(strconv.Itoa(daysValid))+" days",
							green("-- ok!\n"),
						)
						continue
					} else if daysValid < 0 {
						fmt.Printf(
							"%-35s (IPv%d): expired %-44s %s",
							blue(allHosts[i].URL),
							j,
							red(certs[0].NotAfter.Format("02.01.2006"))+", "+red(strconv.Itoa(daysValid))+" days ago",
							red("-- red alert!\n"),
						)
						continue
					} else {
						fmt.Printf(
							"%-35s (IPv%d): expires %-44s %s",
							blue(allHosts[i].URL),
							j,
							orange(certs[0].NotAfter.Format("02.01.2006"))+", in "+orange(strconv.Itoa(daysValid))+" days",
							orange("-- please renew!\n"),
						)
						continue
					}
				}
			}
		},
	}

	rootCmd.PersistentFlags().IntVarP(&days, "days", "d", days, "number of days considered ok with bestbefore/hosts_bestbefore")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", port, "the server port to use")
	rootCmd.PersistentFlags().IntVarP(&protocol, "protocol", "i", protocol, "the IP protocol version to use")
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", timeout, "the timout for dials in seconds")
	rootCmd.PersistentFlags().BoolVarP(&showErrors, "errors", "e", false, "show all dial errors, not only resolved hosts")
	rootCmd.PersistentFlags().StringP("config", "c","", "config file (defaults to ./config.yaml or ~/.config/certinfo/config.yaml")
	rootCmd.AddCommand(cmdListHosts)
	rootCmd.AddCommand(cmdExpires)
	rootCmd.AddCommand(cmdExpiresHosts)
	rootCmd.AddCommand(cmdBestBefore)
	rootCmd.AddCommand(cmdBestBeforeHosts)

	rootCmd.Execute()
}
