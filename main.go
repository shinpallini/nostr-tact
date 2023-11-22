package main

import (
	"context"
	"encoding/json"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/rivo/tview"
)

const (
	cfgname = "config.json"
)

type Config struct {
	Relays       []string `json:"relays"`
	PublicKey    string   `json:"publickey"`
	HexPublicKey string
}

func (c *Config) UnmarshalJSON(b []byte) error {
	type alias Config
	var cfg alias
	err := json.Unmarshal(b, &cfg)
	if err != nil {
		return err
	}
	if strings.HasPrefix(cfg.PublicKey, "npub") {
		_, v, err := nip19.Decode(cfg.PublicKey)
		if err != nil {
			return err
		}
		cfg.HexPublicKey = v.(string)
	}
	*c = Config(cfg)
	return nil
}

type Profile struct {
	Website     string `json:"website"`
	Nip05       string `json:"nip05"`
	Picture     string `json:"picture"`
	Lud16       string `json:"lud16"`
	DisplayName string `json:"display_name"`
	About       string `json:"about"`
	Name        string `json:"name"`
}

func main() {
	// load config
	f, err := os.Open(cfgname)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		panic(err)
	}
	_ = cfg

	// build tui application
	app := tview.NewApplication()
	mainView := tview.NewFlex().SetDirection(tview.FlexRow)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// create in-memory thumbnail storage
	thumbnails := make(map[string]image.Image)
	containers := make([]*tview.Flex, 0)
	// stacking posts on timeline
	go func() {
		for i := 0; ; i++ {
			select {
			case <-time.NewTicker(1 * time.Second).C:
				imageView := tview.NewImage()
				var pic image.Image
				var ok bool
				pic, ok = thumbnails["aaa"]
				if !ok {
					resp, err := http.Get("https://pomf2.lain.la/f/989dxs36.jpg")
					if err != nil {
						panic(err)
					}
					defer resp.Body.Close()
					pic, err = jpeg.Decode(resp.Body)
					if err != nil {
						panic(err)
					}
					thumbnails["aaa"] = pic
				}
				imageView.SetImage(pic)
				content, err := os.UserConfigDir()
				if err != nil {
					panic(err)
				}
				newContainer := tview.NewFlex().
					SetDirection(tview.FlexColumn).
					AddItem(tview.NewTextView().SetText("added: "+strconv.Itoa(i)), 0, 1, false).
					AddItem(imageView, 0, 2, false).
					AddItem(tview.NewTextView().SetText("content: "+content), 0, 6, false).
					AddItem(tview.NewTextView().SetText(time.Now().Format(time.Stamp)), 0, 1, false)

				containers = append(containers, newContainer)
				app.QueueUpdateDraw(func() {
					mainView.Clear()
					stack(containers, mainView)
					// mainView.AddItem(newContainer, 4, 0, false)
					if len(containers) > 10 {
						containers = containers[1:]
					}
				})
			case <-ctx.Done():
				return
			}
		}
	}()

	if err := app.SetRoot(mainView, true).Run(); err != nil {
		panic(err)
	}
}

func stack(containers []*tview.Flex, view *tview.Flex) {
	if len(containers) == 0 {
		return
	}
	for i := len(containers) - 1; i > 0; i-- {
		view.AddItem(containers[i], 4, 0, false)
	}
}
