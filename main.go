package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// 名前、メッセージ、時刻用のTextViewを作成
	nameView := tview.NewTextView().SetDynamicColors(true)
	messageView := tview.NewTextView().SetDynamicColors(true)
	timeView := tview.NewTextView().SetDynamicColors(true)

	// レイアウト用のFlexを作成し、3つのTextViewを追加
	flex := tview.NewFlex().
		AddItem(nameView, 0, 1, false).
		AddItem(messageView, 0, 2, false).
		AddItem(timeView, 0, 1, false)

	go func(ctx context.Context) {
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(1 * time.Second)
				app.QueueUpdateDraw(func() {
					// 現在のテキストを取得
					nameText := nameView.GetText(true)
					messageText := messageView.GetText(true)
					timeText := timeView.GetText(true)

					newText := runewidth.Wrap(fmt.Sprintf("あああああああああああ %d", i), 20)
					paddingLines := strings.Repeat("\n", len(strings.Split(newText, "\n"))-1)

					// 新しいテキストを設定
					nameView.SetText(fmt.Sprintf("Name %d\n%s%s", i, paddingLines, nameText))
					messageView.SetText(fmt.Sprintf("%s\n%s", newText, messageText))
					timeView.SetText(fmt.Sprintf("%s\n%s%s", time.Now().Format("15:04:05"), paddingLines, timeText))
				})
			}
		}
	}(ctx)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
