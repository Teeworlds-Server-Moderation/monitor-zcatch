package parse

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/Teeworlds-Server-Moderation/common/dto"
	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/common/mqtt"
)

var (

	// 0: full 1: ID 2: IP 3: port 4: version 5: name 6: clan 7: country
	playerEnteredRegex = regexp.MustCompile(`id=([\d]+) addr=([a-fA-F0-9\.\:\[\]]+):([\d]+) version=(\d+) name='(.{0,20})' clan='(.{0,16})' country=([-\d]+)$`)
)

// PlayerJoined parsing and creation of the corresponding event JSON struct,
// as well as marshalling that struct into a message payload.
func PlayerJoined(source, timestamp, logLine string) (mqtt.Message, error) {
	match := playerEnteredRegex.FindStringSubmatch(logLine)
	if len(match) != 8 {
		return emptyMsg, fmt.Errorf("Invalid PlayerJoined line format: %s", logLine)
	}
	port, _ := strconv.Atoi(match[3])
	id, _ := strconv.Atoi(match[1])
	country, _ := strconv.Atoi(match[7])
	version, _ := strconv.Atoi(match[4])

	playerJoinEvent := events.NewPlayerJoinedEvent()
	playerJoinEvent.Timestamp = timestamp
	playerJoinEvent.EventSource = source

	player := dto.Player{
		Name:    match[5],
		Clan:    match[6],
		IP:      match[2],
		Port:    port,
		ID:      id,
		Country: country,
		Version: version,
	}

	playerJoinEvent.Player = player
	ServerState.PlayerJoin(id, player)

	msg := mqtt.Message{
		Topic:   events.TypePlayerJoined,
		Payload: playerJoinEvent.Marshal(),
	}
	return msg, nil
}
