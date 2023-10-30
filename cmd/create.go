package cmd

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/melkeydev/go-blueprint/cmd/program"
	"github.com/melkeydev/go-blueprint/cmd/steps"
	"github.com/melkeydev/go-blueprint/cmd/ui/multiInput"
	"github.com/melkeydev/go-blueprint/cmd/ui/textinput"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Go project and don't worry about the structre",
	Long:  "Go Blueprint is a CLI tool that allows you to focus on the actual Go code, and not the project structure. Perfect for someone new to the Go language",

	Run: func(cmd *cobra.Command, args []string) {

		options := steps.Options{
			ProjectName: &textinput.Output{},
		}

		project := &program.Project{
			FrameworkMap: make(map[string]program.Framework),
		}
		steps := steps.InitSteps(&options)

		tprogram := tea.NewProgram(textinput.InitialTextInputModel(options.ProjectName, "What is the name of your project?", project))
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(err)
		}
		project.ExitCLI(tprogram)

		for _, step := range steps.Steps {
			s := &multiInput.Selection{}
			tprogram = tea.NewProgram(multiInput.InitialModelMulti(step.Options, s, step.Headers, project))
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(err)
			}
			project.ExitCLI(tprogram)

			*step.Field = s.Choice
		}

		project.ProjectName = options.ProjectName.Output
		project.ProjectType = strings.ToLower(options.ProjectType)
		currentWorkingDir, err := os.Getwd()
		if err != nil {
			cobra.CheckErr(err)
		}

		project.AbsolutePath = currentWorkingDir

		// This calls the templates
		err = project.CreateMainFile()
		if err != nil {
			cobra.CheckErr(err)
		}

		// Display the message to the user in bullet form
		fmt.Println("\nNext steps:")
		fmt.Printf("• cd %s\n", project.ProjectName)
	},
}
