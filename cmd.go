package yt

import (
	"context"
	"os"
	"text/template"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/json"
	"github.com/rwxrob/bonzai/persisters/inprops"
	"github.com/rwxrob/bonzai/term"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var yt *youtube.Service
var chanid string
var apikey string
var nextpage string

// -------------------------------- Cmd -------------------------------

var Cmd = &bonzai.Cmd{
	Name:  `ytwee`,
	Short: `relay chat messages from YouTube to WeeChat`,
	Comp:  comp.CmdsAliases,
	Def:   help.Cmd,

	Cmds: []*bonzai.Cmd{
		help.Cmd, startCmd, videoCmd, messagesCmd, detailsCmd,
		chatidCmd, nextpageCmd, chanidCmd,
	},

	Persist: inprops.NewUserConfig(`ytwee`, `properties`),
	Vars: bonzai.Vars{
		{K: `yt-channel-id`, E: `YTCHANNELID`, P: true, G: &chanid},
		{K: `yt-api-key`, E: `YTAPIKEY`, P: true, R: true, G: &apikey},
		{K: `yt-chat-next-page`, P: true, G: &nextpage},
	},

	Init: func(x *bonzai.Cmd, _ ...string) error {
		ctx := context.Background()
		service, err := youtube.NewService(ctx, option.WithAPIKey(apikey))
		if err != nil {
			return err
		}
		yt = service
		return nil
	},
}

// ---------------------------- messagesCmd ---------------------------

var messagesCmd = &bonzai.Cmd{
	Name:  `messages`,
	Short: `print up to 200 messages`,
	Vars:  bonzai.Vars{{I: `yt-chat-next-page`}},
	Do: func(x *bonzai.Cmd, _ ...string) error {
		vidid := FetchVideoId(yt, chanid)
		chatid := FetchChatId(yt, vidid)
		messages, nextpage, err := FetchMessages(yt, chatid, nextpage)
		if err != nil {
			return err
		}
		x.Set(`yt-chat-next-page`, nextpage)
		tmpl, _ := template.New("message").Parse(`{{.Author}} {{.Text}}` + "\n")
		for _, message := range messages {
			tmpl.Execute(os.Stdout, message)
		}
		return nil
	},
}

// ----------------------------- startCmd -----------------------------

var startCmd = &bonzai.Cmd{
	Name:  `start`,
	Short: `start relaying messages`,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		println(`would start relaying`)
		return nil
	},
}

// ----------------------------- videoCmd -----------------------------

var videoCmd = &bonzai.Cmd{
	Name:  `video`,
	Short: `unique id of current live stream video`,
	Long: `
		This is the **most** _expensive_ operation available so ***use with
		caution*** and generally only once a stream to cache it someplace
		and reuse it.
	`,
	Do: func(*bonzai.Cmd, ...string) error {
		vidid := FetchVideoId(yt, chanid)
		if vidid != "" {
			term.Print(vidid)
		}
		return nil
	},
}

var chanidCmd = &bonzai.Cmd{
	Name:    `chanid`,
	Short:   `set or get the channel ID`,
	Vars:    bonzai.Vars{{I: `yt-channel-id`}},
	MaxArgs: 1,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) > 0 {
			x.Set(`yt-channel-id`, args[0])
			return nil
		}
		term.Print(x.Get(`yt-channel-id`))
		return nil
	},
}

// ---------------------------- detailsCmd ----------------------------

var detailsCmd = &bonzai.Cmd{
	Name:  `details`,
	Short: `live stream details`,
	Do: func(*bonzai.Cmd, ...string) error {
		vidid := FetchVideoId(yt, chanid)
		details, err := FetchStreamDetails(yt, vidid)
		if err != nil {
			return err
		}
		if vidid != "" {
			term.Print(json.This{details})
		}
		return nil
	},
}

// ----------------------------- chatidCmd ----------------------------

var chatidCmd = &bonzai.Cmd{
	Name:  `chatid`,
	Short: `live stream chat unique identifier`,
	Do: func(*bonzai.Cmd, ...string) error {
		vidid := FetchVideoId(yt, chanid)
		chatid := FetchChatId(yt, vidid)
		if chatid != "" {
			term.Print(chatid)
		}
		return nil
	},
}

// ---------------------------- nextpageCmd ---------------------------

var nextpageCmd = &bonzai.Cmd{
	Name:  `nextpage`,
	Short: `print the next page token that has been cached`,
	Do: func(*bonzai.Cmd, ...string) error {
		term.Print(nextpage)
		return nil
	},
}

var isliveCmd = &bonzai.Cmd{
	Name:  `islive`,
	Short: `check if channal is live`,
	Do: func(*bonzai.Cmd, ...string) error {
		term.Print(nextpage)
		return nil
	},
}
