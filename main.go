package main

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "nostrtact",
		Usage: "This is a TUI-based Nostr client.",
		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "start nostr-tact client",
				Action: func(*cli.Context) error {
					tuiApp := tview.NewApplication()
					commands := tview.NewList().
						AddItem("List item 1", "Some explanatory text", 'a', nil).
						AddItem("List item 2", "Some explanatory text", 'b', nil).
						AddItem("List item 3", "Some explanatory text", 'c', nil).
						AddItem("List item 4", "Some explanatory text", 'd', nil).
						AddItem("Quit", "Press to exit", 'q', func() {
							tuiApp.Stop()
						})
					commands.SetBorder(true)
					commands.SetTitle("Commands")
					// list.SetBorder(true)
					// list.SetTitle("Commands")
					timeline := tview.NewFlex().
						AddItem(tview.NewBox().SetBorder(true).SetTitle("Timeline"), 0, 1, false)
					footer := tview.NewTextView().
						SetDynamicColors(true).
						SetRegions(true).
						SetTextAlign(tview.AlignCenter)
					fmt.Fprintf(footer, "footer")
					center := tview.NewFlex().SetDirection(tview.FlexColumn).
						AddItem(commands, 0, 1, true).
						AddItem(timeline, 0, 2, false)
					footer.SetBorder(true)
					flex := tview.NewFlex().SetDirection(tview.FlexRow).
						AddItem(center, 0, 1, true).
						AddItem(footer, 3, 1, false)
					if err := tuiApp.SetRoot(flex, true).SetFocus(commands).Run(); err != nil {
						return err
					}
					return nil
				},
			},
		},
		DefaultCommand: "start",
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
