package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// file extension name and their ffmpeg decoder/encoder name
type DeEnCoder struct {
	Decoder string
	Encoder string
	Args    []string
}

var mediaTypes = map[string]DeEnCoder{
	".jpg": {Decoder: "mjpeg", Encoder: "mjpeg"},
	".png": {Decoder: "png_pipe", Encoder: "image2"},
	".mov": {Decoder: "mov", Encoder: "mov", Args: []string{
		"-af", "loudnorm=I=-16:LRA=11:TP=-1.5",
	}},
	".mp4": {Decoder: "mp4", Encoder: "mp4", Args: []string{
		"-movflags", "frag_keyframe+empty_moov",
		"-af", "loudnorm=I=-16:LRA=11:TP=-1.5",
	}},
}

var ffmpegLock sync.Mutex

func getMedia(s string) string {
	lower := strings.ToLower(s)
	for p := range mediaTypes {
		if strings.HasSuffix(lower, p) {
			return p
		}
	}
	return ""
}

func fromMedia(s string) bool {
	return strings.HasPrefix(s, "ppt/media/") || // pptx format
		strings.HasPrefix(s, "Pictures") || // opendocument format
		strings.HasPrefix(s, "word/media") // docx format
}

func Transcode(r *zip.Reader, w *zip.Writer) []error {
	errs := make([]error, 0)
	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			log.Println("Can not decode", f.Name, err)
			errs = append(errs, err)
			continue
		}

		writer, err := w.Create(f.Name)
		if err != nil {
			log.Println("Can not create file", f.Name, err)
			errs = append(errs, err)
			continue
		}

		// not media, compress and continue
		pattern := getMedia(f.Name)
		if !fromMedia(f.Name) || pattern == "" {
			io.Copy(writer, rc)
			continue
		}

		var args = []string{
			"-f", mediaTypes[pattern].Decoder,
			"-i", "-",
			"-movflags", "frag_keyframe+empty_moov"}
		args = append(args, mediaTypes[pattern].Args...)
		args = append(args, "-f", mediaTypes[pattern].Encoder, "-")

		cmd := exec.Command("ffmpeg", args...)
		cmd.Stdout = writer
		// cmd.Stderr = os.Stderr

		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Println("Error obtaining stdin", err)
			errs = append(errs, err)
			continue
		}

		err = cmd.Start()
		if err != nil {
			log.Println("Start ffmpeg failed", err)
			errs = append(errs, err)
			continue
		}

		_, err = io.Copy(stdin, rc)
		if err != nil {
			log.Println("Can not copy pipe", err)
			errs = append(errs, err)
			continue
		}
		stdin.Close()

		err = cmd.Wait()
		if err != nil {
			log.Println("Ffmpeg report error", err)
			errs = append(errs, err)
			continue
		}
	}

	return nil
}

func logBegin(c *gin.Context) {
	log.Println("Start compress task")
}

func main() {
	route := gin.Default()
	route.GET("/*any", func(c *gin.Context) {
		c.File("index.html")
	})
	route.POST("/*any", logBegin, func(c *gin.Context) {
		header, err := c.FormFile("file")
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.Header("Content-Disposition", "inline; filename*=UTF-8''"+header.Filename)

		file, err := header.Open()
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		r, err := zip.NewReader(file, header.Size)
		w := zip.NewWriter(c.Writer)

		errs := Transcode(r, w)
		if len(errs) != 0 {
			c.AbortWithError(500, errors.New(fmt.Sprint(errs)))
			return
		}

		err = w.Close()
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

	})
	route.Run(":8888")
}
