package yt

import (
	"fmt"
	"log"

	"google.golang.org/api/youtube/v3"
)

type ChatMessage struct {
	Id     string
	Author string
	Text   string
	Time   string
}

func FetchMessages(yt *youtube.Service, chatid, pagetok string) ([]ChatMessage, string, error) {
	call := yt.LiveChatMessages.List(
		chatid, []string{"snippet", "authorDetails"}).
		MaxResults(200).
		PageToken(pagetok)

	response, err := call.Do()
	if err != nil {
		return nil, pagetok, err
	}

	messages := []ChatMessage{}
	for _, i := range response.Items {
		m := ChatMessage{}
		m.Id = i.Id
		m.Text = i.Snippet.DisplayMessage
		m.Time = i.Snippet.PublishedAt
		m.Author = i.AuthorDetails.DisplayName
		messages = append(messages, m)
	}

	pagetok = response.NextPageToken
	return messages, pagetok, nil
}

func FetchVideoId(yt *youtube.Service, chanid string) string {

	// Search for live streams on the channel
	call := yt.Search.List([]string{"id", "snippet"}).
		ChannelId(chanid). // Filter by channel ID
		EventType("live"). // Filter for live events
		Type("video").     // Restrict to videos
		MaxResults(1)      // Only retrieve one result

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making API call: %v", err)
	}

	// Check if a live stream was found
	if len(response.Items) == 0 {
		return ""
	}

	// Extract and print the live stream ID
	return response.Items[0].Id.VideoId
}

func FetchStreamDetails(yt *youtube.Service, vidid string) (
	*youtube.VideoLiveStreamingDetails, error) {
	call := yt.Videos.List([]string{"liveStreamingDetails"}).Id(vidid)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}
	if len(resp.Items) == 0 {
		return nil, fmt.Errorf(`no video found for %v`, vidid)
	}
	return resp.Items[0].LiveStreamingDetails, nil
}

func FetchChatId(yt *youtube.Service, vidid string) string {
	details, err := FetchStreamDetails(yt, vidid)
	if err != nil {
		log.Print(err)

		return ""
	}
	return details.ActiveLiveChatId
}

/*
func IsChannelLive(name string) bool {
	url := fmt.Sprintf("https://www.youtube.com/@%s/live", channelName)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}
*/
