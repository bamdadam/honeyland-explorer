package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "honeyland-explorer",
	Short: "Honeyland explorer is a nft exploration app",
	Long:  `An Nft exploration app written in go for The HoneyVerse`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// Do Stuff Here
	// },
}

func Execute() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	// log.AddHook(logruseq.NewSeqHook("http://localhost:5341"))
	Register(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("error executing command")
	}
}
