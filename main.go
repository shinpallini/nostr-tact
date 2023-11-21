package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nbd-wtf/go-nostr"
	"github.com/rivo/tview"
)

type metadataContent struct {
	Name        string `json:"name"`
	DispalyName string `json:"display_name"`
	Picture     string `json:"picture"`
}

func main() {
	// sk := nostr.GeneratePrivateKey()
	// pk, _ := nostr.GetPublicKey(sk)
	app := tview.NewApplication()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// リレーサーバーの設定
	relay, err := nostr.RelayConnect(ctx, "wss://relay-jp.nostr.wirednet.jp")
	if err != nil {
		panic(err)
	}
	filters := []nostr.Filter{{
		Kinds: []int{nostr.KindTextNote},
		Limit: 1,
	}}

	sub, err := relay.Subscribe(ctx, filters)
	if err != nil {
		panic(err)
	}

	// userNames := func() map[string]string {
	// 	m := make(map[string]string)
	// 	filters := []nostr.Filter{{
	// 		Kinds: []int{nostr.KindProfileMetadata},
	// 		Limit: 10000,
	// 	}}
	// 	sub, err := relay.Subscribe(ctx, filters)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	for {
	// 		select {
	// 		case ev := <-sub.Events:
	// 			var content metadataContent
	// 			err = json.Unmarshal([]byte(ev.Content), &content)
	// 			if err != nil {
	// 				// panic(err)
	// 				continue
	// 			}
	// 			m[ev.PubKey] = content.DispalyName
	// 		case <-time.After(2 * time.Second):
	// 			return m
	// 		case <-ctx.Done():
	// 			return m
	// 		}
	// 	}
	// }()
	// f, err := os.Create("length.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// fmt.Fprintf(f, "user count: %d", len(userNames))

	// 名前、メッセージ、時刻用のTextViewを作成
	nameView := tview.NewTextView().SetDynamicColors(true)
	messageView := tview.NewTextView().SetDynamicColors(true)
	timeView := tview.NewTextView().SetDynamicColors(true)
	imageView := tview.NewImage()
	resp, err := http.Get("https://pomf2.lain.la/f/989dxs36.jpg")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	photo, err := jpeg.Decode(resp.Body)
	if err != nil {
		panic(err)
	}
	imageView.SetImage(photo)

	// レイアウト用のFlexを作成し、3つのTextViewを追加
	flex := tview.NewFlex().
		AddItem(nameView, 0, 1, false).
		AddItem(messageView, 0, 4, false).
		AddItem(timeView, 0, 2, false).
		AddItem(imageView, 0, 1, false)

	// var mu sync.Mutex
	go func(ctx context.Context) {
		for ev := range sub.Events {
			select {
			case <-ctx.Done():
				return
			default:
				app.QueueUpdateDraw(func() {
					// 現在のテキストを取得
					nameText := truncate(nameView.GetText(true))
					messageText := truncate(messageView.GetText(true))
					timeText := truncate(timeView.GetText(true))

					newText := runewidth.Wrap(fmt.Sprintf(ev.Content), 70)
					paddingLines := strings.Repeat("\n", len(strings.Split(newText, "\n"))-1)

					// 新しいテキストを設定
					// mu.Lock()
					nameView.SetText(fmt.Sprintf("%s\n%s%s", getName(ctx, relay, ev.PubKey), paddingLines, nameText))
					// mu.Unlock()
					messageView.SetText(fmt.Sprintf("%s\n%s", newText, messageText))
					timeView.SetText(fmt.Sprintf("%s\n%s%s", ev.CreatedAt.Time(), paddingLines, timeText))
					// time.Sleep(1 * time.Second)
				})
			}
		}
	}(ctx)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func getName(ctx context.Context, relay *nostr.Relay, pubKey string) string {
	filters := []nostr.Filter{{
		Kinds:   []int{nostr.KindProfileMetadata},
		Authors: []string{pubKey},
		Limit:   1,
	}}
	sub, err := relay.Subscribe(ctx, filters)
	if err != nil {
		panic(err)
	}
	select {
	case ev := <-sub.Events:
		var content metadataContent
		err = json.Unmarshal([]byte(ev.Content), &content)
		if err != nil {
			panic(err)
		}
		return content.DispalyName
	case <-ctx.Done():
		return ""
	case <-time.After(2 * time.Second):
		return ""
	}
}

func truncate(s string) string {
	if len(s) > 1001 {
		return s[:1000]
	}
	return s
}
