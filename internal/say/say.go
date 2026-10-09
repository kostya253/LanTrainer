package say

import (
	"os"
	"os/exec"
)

func Say(text string) error {
	voice := os.Getenv("SAY_VOICE")
	voiceSpeed := os.Getenv("SAY_SPEED")

	say := exec.Command("say", "-v", voice, "-r", voiceSpeed, "-o", "voice.wav", "--data-format=LEF32@22050", text)
	if err := say.Run(); err != nil {
		return err
	}

	ffmpeg := exec.Command(
		"ffmpeg",
		"-y",
		"-i", "voice.wav",
		"-c:a", "libopus",
		"-b:a", "32k",
		"-vbr", "on",
		"-compression_level", "10",
		"voice.ogg",
	)
	ffmpeg.Stderr = nil // capture_output=True discards stderr in Python
	if err := ffmpeg.Run(); err != nil {
		return err
	}

	return nil
}
