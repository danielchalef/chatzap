package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	audioFile := flag.String("audio", "", "Path to audio file")
	outputFile := flag.String("output", "", "Path to output file")
	prompt := flag.String("prompt", "", "Prompt text")
	// outputFormat := flag.String("format", "text", "Output format: json, text, srt, verbose_json, or vtt")
	openai_key := flag.String("key", "", "OpenAI API key")

	flag.Parse()

	if *audioFile == "" || *outputFile == "" {
		fmt.Println("Please provide audio file and output file")
		os.Exit(1)
	}

	if *openai_key == "" {
		if os.Getenv("OPENAI_API_KEY") == "" {
			fmt.Println("Please provide OpenAI API key either as a flag or set the OPENAI_API_KEY environment variable")
			os.Exit(1)
		} else {
			*openai_key = os.Getenv("OPENAI_API_KEY")
		}
	}

	c := openai.NewClient(*openai_key)
	ctx := context.Background()

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: *audioFile,
		Prompt:   *prompt,
	}
	resp, err := c.CreateTranscription(ctx, req)
	check(err)
	fmt.Println(resp.Text)

	f, err := os.Create(*outputFile)
	check(err)

	defer f.Close()
	_, err = f.WriteString(resp.Text)
	check(err)
	f.Sync()

	fmt.Println("Done!")

}
