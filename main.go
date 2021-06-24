package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	inDir   = flag.String("input", "input", "Directory of input iso files")
	outDir  = flag.String("output", "output", "Directory of output mp4 files")
	privacy = flag.String("privacy", "unlisted", "Video privacy status")
)

func main() {
	flag.Parse()

	if *inDir == "" {
		log.Fatalf("Input directory is required")
	}

	entries, err := os.ReadDir(*inDir)
	if err != nil {
		log.Fatalf("Cannot read directory %s: %v", *inDir, err)
	}

	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".iso") {
			continue
		}
		inFilepath := filepath.Join(*inDir, e.Name())
		base := e.Name()[:strings.LastIndex(e.Name(), ".iso")]
		outFilepath := filepath.Join(*outDir, base) + ".mp4"

		// encode
		if err := encode(inFilepath, outFilepath); err != nil {
			log.Printf("Failed to encode %s: %v", e.Name(), err)
			continue
		}
		if err := os.Rename(inFilepath, filepath.Join("iso", e.Name())); err != nil {
			log.Printf("Failed to archive %s: %v", inFilepath, err)
			continue
		}

		// upload
		err = uploadVideo(UploadParam{
			Title:    base,
			Privacy:  *privacy,
			Filename: outFilepath,
		})
		if err != nil {
			log.Printf("Failed to upload %s: %v", base, err)
			continue
		}
		if err := os.Rename(outFilepath, filepath.Join("mp4", base+".mp4")); err != nil {
			log.Printf("Failed to archive %s: %v", outFilepath, err)
			continue
		}
	}

	fmt.Println("Done!!")
}

// encode invokes ffmpeg
func encode(in, out string) error {
	cmd := exec.Command("ffmpeg", "-i", in, out)
	return cmd.Run()
}
