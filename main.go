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

	entries, err := os.ReadDir(*inDir)
	if err != nil {
		log.Fatalf("Cannot read directory %s: %v", *inDir, err)
	}

	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".iso") {
			continue
		}
		inFilepath := filepath.Join(*inDir, e.Name())
		base := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
		outFilepath := filepath.Join(*outDir, base+".mp4")

		// encode
		if err := encode(inFilepath, outFilepath); err != nil {
			log.Printf("Failed to encode %s: %v", inFilepath, err)
			continue
		}
		if err := move(inFilepath, "iso"); err != nil {
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
			log.Printf("Failed to upload %s: %v", outFilepath, err)
			continue
		}
		if err := move(outFilepath, "mp4"); err != nil {
			log.Printf("Failed to archive %s: %v", outFilepath, err)
			continue
		}
	}

	fmt.Println("Done!!")
}

func move(srcFile, destDir string) error {
	info, err := os.Stat(destDir)
	if err != nil {
		return err
	}
	if info.IsDir() {
		filename := filepath.Base(srcFile)
		return os.Rename(srcFile, filepath.Join(destDir, filename))
	}
	return fmt.Errorf("destDir %s is not a directory", destDir)
}

// encode invokes ffmpeg
func encode(in, out string) error {
	// TODO: use context
	log.Printf("Encoding %s ...", in)
	cmd := exec.Command("ffmpeg", "-i", in, out)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
