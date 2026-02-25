package main

import (
	"fmt"
	"os"

	"github.com/EwanClarke/postfixcalc/internal/tui/graph_ui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "plot [expression]",
	Short: "A novel graph plotting utility",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if err := graphui.Start(); err != nil {
				fmt.Printf("Error: %v", err)
				os.Exit(1)
			}
			return
		}

		err := showGraph(args[0])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func showGraph(expression string) error {
	width, height := 40, 20
	c := graphui.NewCanvas()
	c.Resize(width, height)
	c.SetBounds(-10, 10, -5, 5)
	c.ChangeExpression(expression)

	result, err := c.Render()
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
