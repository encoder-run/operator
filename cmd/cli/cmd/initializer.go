package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the application config",
	Long: `This command helps set up the initial configuration for the application
by walking you through a series of prompts.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupConfig()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func ensureConfigExists() {
	// Check if the configuration directory exists
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	configDir := path.Join(home, ".encoder-run")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// Create the configuration directory
		err = os.Mkdir(configDir, 0755)
		if err != nil {
			fmt.Printf("Error creating config directory: %v\n", err)
			os.Exit(1)
		}

		// Create blank configuration file
		configFile := path.Join(configDir, "config")
		file, err := os.Create(configFile)
		if err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
	}

	// Ensure there is a configuration file present
	if _, err := os.Stat(path.Join(configDir, "config")); os.IsNotExist(err) {
		// Write a blank configuration file
		viper.SetConfigFile(path.Join(configDir, "config"))
		err = viper.WriteConfig()
		if err != nil {
			fmt.Printf("Error writing config: %v\n", err)
			os.Exit(1)
		}
	}
}

func setupConfig() {
	// OpenAI Token prompt
	openAITokenPrompt := promptui.Prompt{
		Label: "OpenAI Token (leave blank if not applicable)",
		Mask:  '*', // Mask input for privacy
	}

	openAIToken, err := openAITokenPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	// GitHub Personal Access Token prompt
	githubTokenPrompt := promptui.Prompt{
		Label: "GitHub Personal Access Token (leave blank if not applicable)",
		Mask:  '*', // Mask input for privacy
	}

	githubToken, err := githubTokenPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	// Process and save the inputs
	viper.Set("tokens.openai", openAIToken)
	viper.Set("tokens.github", githubToken)

	ensureConfigExists()

	// Write the configuration to a file
	err = viper.WriteConfig()
	if err != nil {
		fmt.Printf("Error writing config: %v\n", err)
	} else {
		fmt.Println("Configuration saved successfully.")
	}
}
