package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
)

//go:embed kickstart.mp3
var song []byte

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "error: %v\n", fmt.Errorf("no commands specified"))
		os.Exit(-1)
	}

	decoder, err := minimp3.NewDecoder(bytes.NewReader(song))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(-1)
	}
	defer decoder.Close()
	<-decoder.Started()

	var context *oto.Context
	if context, err = oto.NewContext(decoder.SampleRate, decoder.Channels, 2, 1024); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(-1)
	}

	var player = context.NewPlayer()
	defer player.Close()

	go io.Copy(player, decoder)

	shellExecutable := os.Args[1]
	shellParameters := []string{}

	if len(os.Args) > 2 {
		shellParameters = os.Args[2:]
	}

	cmd := exec.Command(shellExecutable, shellParameters...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	_ = cmd.Run()
}
