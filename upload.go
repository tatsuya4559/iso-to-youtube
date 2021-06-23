package main

import (
	"fmt"
	"os"

	"google.golang.org/api/youtube/v3"
)

type UploadParam struct {
	Title    string
	Privacy  string
	Filename string
}

func uploadVideo(p UploadParam) error {
	client := getClient(youtube.YoutubeUploadScope)
	service, err := youtube.New(client)
	if err != nil {
		return fmt.Errorf("Error creating YouTube client: %w", err)
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{Title: p.Title},
		Status:  &youtube.VideoStatus{PrivacyStatus: p.Privacy},
	}

	call := service.Videos.Insert([]string{"snippet", "status"}, upload)

	file, err := os.Open(p.Filename)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("Error opening %s: %w", p.Filename, err)
	}

	fmt.Println("Uploading...")
	response, err := call.Media(file).Do()
	if err != nil {
		return fmt.Errorf("Error while uploading: %w", err)
	}
	fmt.Printf("Upload successful! Video ID: %v\n", response.Id)

	return nil
}
