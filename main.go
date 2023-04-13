package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/assistent-ai/client/chat"
	"github.com/assistent-ai/client/cli"
	"github.com/assistent-ai/client/db"
	"github.com/assistent-ai/client/model"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "jessica",
		Short: "Jessica is an AI assistent.",
	}

	apiKeyFilePath := ""
	defaultFilePath := filepath.Join(os.Getenv("HOME"), ".open-ai.key")
	rootCmd.PersistentFlags().StringVar(&apiKeyFilePath, "key-file", defaultFilePath, "Path to the text file containing the API key")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all conversations",
		Run: func(cmd *cobra.Command, args []string) {
			dialogIds, err := db.GetDialogIDs()
			if err != nil {
				cli.PrintErrorAndExit(err)
			} else {
				cli.PrintDialogIDs(dialogIds)
			}
		},
	}

	startCmd := &cobra.Command{
		Use:   "continue [id]",
		Short: "Continue dialog with id. If id does not exist it will be created",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cli.PrintErrorAndExit(errors.New("please provide dialog id"))
			} else {
				id := args[0]
				// Replace with your actual logic to start a new conversation
				fmt.Println("Starting a new conversation...")
				err := chat.StartChat(id)
				if err != nil {
					cli.PrintErrorAndExit(err)
				}
			}
		},
	}

	showCmd := &cobra.Command{
		Use:   "show [id]",
		Short: "Show conversation by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cli.PrintErrorAndExit(errors.New("please provide dialog id"))
			} else {
				id := args[0]
				messages, err := db.GetMessagesByDialogID(id)
				if err != nil {
					cli.PrintErrorAndExit(err)
				} else {
					chat.ShowMessages(messages)
				}
			}
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete conversation by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cli.PrintErrorAndExit(errors.New("please provide dialog id"))
			} else {
				id := args[0]
				err := db.RemoveMessagesByDialogId(id)
				if err != nil {
					cli.PrintErrorAndExit(err)
				}
			}
		},
	}

	deleteDefaultCmd := &cobra.Command{
		Use:   "delete-default",
		Short: "Delete default conversation",
		Run: func(cmd *cobra.Command, args []string) {
			err := db.RemoveMessagesByDialogId(model.DefaultDialogId)
			if err != nil {
				cli.PrintErrorAndExit(err)
			}
		},
	}

	restartDefaultCmd := &cobra.Command{
		Use:   "restart-default",
		Short: "Removes all default messages and starts new default from scratch",
		Run: func(cmd *cobra.Command, args []string) {
			err := db.RemoveMessagesByDialogId(model.DefaultDialogId)
			if err != nil {
				cli.PrintErrorAndExit(err)
			} else {
				err := chat.StartChat(model.DefaultDialogId)
				if err != nil {
					cli.PrintErrorAndExit(err)
				}
			}
		},
	}

	rootCmd.AddCommand(listCmd, startCmd, showCmd, deleteCmd, deleteDefaultCmd, restartDefaultCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
