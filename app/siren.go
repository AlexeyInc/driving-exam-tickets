package app

import (
	"fmt"
	"os/exec"
	"time"
)

func PlaySiren() {
	playCount := 3
	for range playCount {
		// Play the sound using afplay (macOS)
		cmd := exec.Command("afplay", "../media/police.wav")
		err := cmd.Start()
		if err != nil {
			fmt.Println("Error playing sound:", err)
			return
		}
		cmd.Wait()

		time.Sleep(time.Millisecond * 5)
	}
	fmt.Println("Good luck :)")
}
