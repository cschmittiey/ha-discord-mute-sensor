package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	TargetUserID = "TARGET_USER_ID"
)

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// We only care about voice state updates, so we can ignore other events.
	dg.Identify.Intents = discordgo.IntentsGuildVoiceStates

	dg.AddHandler(voiceStateUpdate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func voiceStateUpdate(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	// Check if the user ID matches the specific user
	if m.UserID == os.Getenv("DISCORD_USER_ID") {
		// Check if the user has just joined a voice channel
		if m.ChannelID != "" {
			if m.SelfMute {
				fmt.Println("muted")
				notifyHA(true)
			} else {
				fmt.Println("no longer muted")
				notifyHA(false)
			}
		} else {
			fmt.Println("no longer muted because no longer in channel")
			notifyHA(false)
		}
	}
}

func notifyHA(sensor_state bool) {
	// slap that HA HTTP endpoint with a big ol' POST
	// true should mean sensor on
	// false should mean sensor off

	url := os.Getenv("HA_BASE_URL") + "/api/states/binary_sensor." + os.Getenv("HA_SENSOR_NAME")

	fmt.Println("url: " + url)

	jsonStr := []byte(``) // i'm not proud of this. surely there is a better way, but i will not be doing it at this time

	if sensor_state {
		jsonStr = []byte(`{"state": "on", "attributes": {"friendly_name": "Discord Muted Sensor"}}`)
	} else {
		jsonStr = []byte(`{"state": "off", "attributes": {"friendly_name": "Discord Muted Sensor"}}`)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("HA_AUTH_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// handy debug statements
	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
}
