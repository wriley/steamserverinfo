// steamserverinfo.go
// https://github.com/wriley/steamserverinfo
//
// A server query program written in golang that uses the Steam Server Query API
// https://developer.valvesoftware.com/wiki/Server_queries
//
// Ported from original C version found here https://github.com/wriley/arma2serverinfo

package main

import (
	"fmt"
	"net"
	"os"
    "time"
    "encoding/binary"
    "bytes"
)

func CheckError(err error) bool {
	if err != nil {
		fmt.Println("Error: ", err)
        return 1 == 0
	} else {
        return 1 == 1
    }
}

func SendPacket(conn net.Conn, arr []byte, timeout time.Duration) (int, []byte) {
    fmt.Fprintln(os.Stderr, "Writing...")
    ret, err := conn.Write(arr)
    if CheckError(err) {
        fmt.Fprintf(os.Stderr, "Wrote %d bytes\n", ret)
        buffer := make([]byte, 1024)
        fmt.Fprintln(os.Stderr, "Reading...")
        conn.SetReadDeadline(time.Now().Add(timeout))
        n, err := conn.Read(buffer)
        if CheckError(err) {
            return n, buffer
        } else {
            return 0, nil
        }
    } else {
        return 0, nil
    }

}

func GetString(arr []byte, index int) (string, int) {
    data := ""
    for i := index; i < len(arr); i++ {
        index = i
        if arr[i] == 0x00 {
            break
        } else {
            data = data + string(arr[i])
        }
    }
    index++

    return data, index
}

func main() {
	A2S_INFO := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x54, 0x53, 0x6F, 0x75, 0x72, 0x63, 0x65, 0x20, 0x45, 0x6E, 0x67, 0x69, 0x6E, 0x65, 0x20, 0x51, 0x75, 0x65, 0x72, 0x79, 0x00}
    A2S_RULES := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56, 0xFF, 0xFF, 0xFF, 0xFF}
    A2S_PLAYER := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x55, 0xFF, 0xFF, 0xFF, 0xFF}

	argsWithProg := os.Args
	if len(argsWithProg) != 3 {
		fmt.Printf("Usage: %s <server> <port>\n", argsWithProg[0])
		os.Exit(1)
	}

	server := argsWithProg[1]
	port := argsWithProg[2]
    seconds := 5
    timeout := time.Duration(seconds) * time.Second

    fmt.Fprintln(os.Stderr, "Opening UDP connection...")
	Conn, err := net.DialTimeout("udp", server+":"+port, timeout)
	if !CheckError(err) {
        os.Exit(2)
    }

	defer Conn.Close()

    // Get Info

    fmt.Fprintln(os.Stderr, "Sending A2S_INFO...")
	n, BytesReceived := SendPacket(Conn, A2S_INFO, timeout)

    if BytesReceived == nil || n == 0 {
        fmt.Fprintln(os.Stderr, "Received no data!")
        os.Exit(2)
    }

    if BytesReceived[4] != 0x49 {
        fmt.Fprintf(os.Stderr, "Header was 0x%x instead of 0x49\n", BytesReceived[4])
        os.Exit(2)
    }

    fmt.Fprintf(os.Stderr, "HEADER: 0x%x\n", BytesReceived[4])
    fmt.Fprintf(os.Stderr, "PROTOCOL: 0x%x\n", BytesReceived[5])

    sPtr := 5
    info, sPtr := GetString(BytesReceived, sPtr)
    fmt.Printf("NAME: %s\n", info)

    info, sPtr = GetString(BytesReceived, sPtr)
    fmt.Printf("MAP: %s\n", info)

    info, sPtr = GetString(BytesReceived, sPtr)
    fmt.Printf("FOLDER: %s\n", info)

    info, sPtr = GetString(BytesReceived, sPtr)
    fmt.Printf("GAME: %s\n", info)

    id1 := BytesReceived[sPtr]
    sPtr++
    id2 := BytesReceived[sPtr]
    sPtr++
    id :=  uint16(id1) | uint16(id2)<<8
    fmt.Printf("ID: %d\n", id)

    fmt.Printf("PLAYERS: %d\n", BytesReceived[sPtr])
    sPtr++

    fmt.Printf("MAXPLAYERS: %d\n", BytesReceived[sPtr])
    sPtr++

    fmt.Printf("BOTS: %d\n", BytesReceived[sPtr])
    sPtr++

    fmt.Printf("SERVERTYPE: %c\n", BytesReceived[sPtr])
    sPtr++

    fmt.Printf("ENVIRONMENT: %c\n", BytesReceived[sPtr])
    sPtr++

    fmt.Printf("VISIBILITY: %d\n", BytesReceived[sPtr])
    sPtr++

    fmt.Printf("VAC: %d\n", BytesReceived[sPtr])
    sPtr++

    info, sPtr = GetString(BytesReceived, sPtr)
    fmt.Printf("VERSION: %s\n", info)

    if n > sPtr {
        // EDF
        edf := BytesReceived[sPtr]
        sPtr++

        // PORT
        if edf & 0x80 != 0 {
            s1 := BytesReceived[sPtr]
            sPtr++
            s2 := BytesReceived[sPtr]
            sPtr++
            port := uint16(s1) | uint16(s2)<<8
            fmt.Printf("PORT: %d\n", port)
        }

        // STEAMID
        if edf & 0x10 != 0 {
            sPtr += 8
        }

        // Keywords
        if edf & 0x20 != 0 {
            info, sPtr = GetString(BytesReceived, sPtr)
            fmt.Printf("KEYWORDS: %s\n", info)
        }
    }

    // Get Rules
    sPtr = 5

    fmt.Fprintln(os.Stderr, "Sending A2S_RULES...")
	n, BytesReceived = SendPacket(Conn, A2S_RULES, timeout)

    if BytesReceived == nil || n == 0 {
        fmt.Fprintln(os.Stderr, "Received no data!")
        os.Exit(2)
    }

    if BytesReceived[4] != 0x41 {
        fmt.Fprintf(os.Stderr, "Header was 0x%x instead of 0x41\n", BytesReceived[4])
        os.Exit(2)
    }

    // Challenge number
    i1 := BytesReceived[sPtr]
    sPtr++
    i2 := BytesReceived[sPtr]
    sPtr++
    i3 := BytesReceived[sPtr]
    sPtr++
    i4 := BytesReceived[sPtr]
    sPtr++
    chnum := uint32(i4)<<24 | uint32(i3)<<16 | uint32(i2)<<8 | uint32(i1)
    fmt.Fprintf(os.Stderr,"Challenge number: %d\n", chnum)

    A2S_RULES[5] = byte(chnum)
    A2S_RULES[6] = byte(chnum >> 8)
    A2S_RULES[7] = byte(chnum >> 16)
    A2S_RULES[8] = byte(chnum >> 24)

    fmt.Fprintln(os.Stderr, "Sending A2S_RULES...")
	n, BytesReceived = SendPacket(Conn, A2S_RULES, timeout)

    if BytesReceived == nil || n == 0 {
        fmt.Fprintln(os.Stderr, "Received no data!")
        os.Exit(2)
    }

    if BytesReceived[4] != 0x45 {
        fmt.Fprintf(os.Stderr, "Header was 0x%x instead of 0x45\n", BytesReceived[4])
        os.Exit(2)
    }

    // reset sPtr
    sPtr = 5

    s1 := BytesReceived[sPtr]
    sPtr++
    s2 := BytesReceived[sPtr]
    sPtr++
    rules := uint16(s1) | uint16(s2)<<8

    if(rules > 0) {
        fmt.Println("RULE LIST:")
    }

    for i := uint16(0); i < rules; i++ {
        // Name
        info, sPtr = GetString(BytesReceived, sPtr)
        // Value
        val := ""
        val, sPtr = GetString(BytesReceived, sPtr)

        fmt.Printf("%s %s\n", info, val)
    }

    // Get Players
    sPtr = 5

    fmt.Fprintln(os.Stderr, "Sending A2S_PLAYER...")
	n, BytesReceived = SendPacket(Conn, A2S_PLAYER, timeout)

    if BytesReceived == nil || n == 0 {
        fmt.Fprintln(os.Stderr, "Received no data!")
        os.Exit(2)
    }

    if BytesReceived[4] != 0x41 {
        fmt.Fprintf(os.Stderr, "Header was 0x%x instead of 0x41\n", BytesReceived[4])
        os.Exit(2)
    }

    // Challenge number
    i1 = BytesReceived[sPtr]
    sPtr++
    i2 = BytesReceived[sPtr]
    sPtr++
    i3 = BytesReceived[sPtr]
    sPtr++
    i4 = BytesReceived[sPtr]
    sPtr++
    chnum = uint32(i4)<<24 | uint32(i3)<<16 | uint32(i2)<<8 | uint32(i1)
    fmt.Fprintf(os.Stderr,"Challenge number: %d\n", chnum)

    A2S_PLAYER[5] = byte(chnum)
    A2S_PLAYER[6] = byte(chnum >> 8)
    A2S_PLAYER[7] = byte(chnum >> 16)
    A2S_PLAYER[8] = byte(chnum >> 24)

    fmt.Fprintln(os.Stderr, "Sending A2S_PLAYER...")
	n, BytesReceived = SendPacket(Conn, A2S_PLAYER, timeout)

    if BytesReceived == nil || n == 0 {
        fmt.Fprintln(os.Stderr, "Received no data!")
        os.Exit(2)
    }

    if BytesReceived[4] != 0x44 {
        fmt.Fprintf(os.Stderr, "Header was 0x%x instead of 0x44\n", BytesReceived[4])
        os.Exit(2)
    }

    sPtr = 5
    players := BytesReceived[sPtr]
    sPtr++

    if players > 0 {
        fmt.Println("PLAYER LIST:");
    }

    for i := 0; i < int(players); i++ {
        // Index (this seems to always be 0, so skipping it)
        sPtr++

        // Name
        info, sPtr = GetString(BytesReceived, sPtr)

        // Score
        i1 = BytesReceived[sPtr]
        sPtr++
        i2 = BytesReceived[sPtr]
        sPtr++
        i3 = BytesReceived[sPtr]
        sPtr++
        i4 = BytesReceived[sPtr]
        sPtr++
        score := uint32(i4)<<24 | uint32(i3)<<16 | uint32(i2)<<8 | uint32(i1)

        // Duration
        b := []byte{0x00, 0x00, 0x00, 0x00}
        b[0] = BytesReceived[sPtr]
        sPtr++
        b[1] = BytesReceived[sPtr]
        sPtr++
        b[2] = BytesReceived[sPtr]
        sPtr++
        b[3] = BytesReceived[sPtr]
        sPtr++
        var duration float32
        buf := bytes.NewReader(b)
        err := binary.Read(buf, binary.LittleEndian, &duration)
        if err != nil {
            fmt.Fprintln(os.Stderr, "Float conversion failed:", err)
        }

        fmt.Printf("%s %d %.0f\n", info, score, duration)
    }
}
