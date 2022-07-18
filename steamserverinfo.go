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
	"path/filepath"
	"encoding/hex"
)

var debug bool = false
var colorize bool = false

func CheckNoError(err error) bool {
	if err != nil {
		fmt.Println("Error: ", err)
        return  false
	} else {
        return true
    }
}

func Colorize(s string) string {
	return "\033[1;31m" + s + "\033[0m"
}

func MyHexDump(arr []byte, s int) string {
	var b = make([]byte, s)
	for i := 0; i < s; i++ {
		b[i] = arr[i]
	}
	if colorize {
		return Colorize(hex.Dump(b))
	} else {
		return hex.Dump(b)
	}
}

func SendPacket(conn net.Conn, arr []byte, timeout time.Duration) (int, []byte) {
    if(debug) {	fmt.Fprintln(os.Stderr, "Writing...") }
    ret, err := conn.Write(arr)
	if(debug) { fmt.Fprintf(os.Stderr, MyHexDump(arr, ret)) }
    if CheckNoError(err) {
        if(debug) {	fmt.Fprintf(os.Stderr, "Wrote %d bytes\n", ret) }
        buffer := make([]byte, 1500)
        if(debug) {	fmt.Fprintln(os.Stderr, "Reading...") }
        conn.SetReadDeadline(time.Now().Add(timeout))
        n, err := conn.Read(buffer)
        if CheckNoError(err) {
			if(debug) {	fmt.Fprintf(os.Stderr, "Read %d bytes\n", n) }
			if(debug) { fmt.Fprintf(os.Stderr, MyHexDump(buffer, n)) }
            return n, buffer
        } else {
            return 0, nil
        }
    } else {
        return 0, nil
    }

}

func stripCtlAndExtFromBytes(str string) string {
	b := make([]byte, len(str))
	var bl int
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
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
	data = stripCtlAndExtFromBytes(data)
    return data, index
}

func GetUInt16(arr []byte, index int) (uint16, int) {
    num1 := arr[index]
    index++
    num2 := arr[index]
    index++
    num :=  uint16(num1) | uint16(num2)<<8
    return num, index
}

func GetUInt32(arr []byte, index int) (uint32, int) {
    num1 := arr[index]
    index++
    num2 := arr[index]
    index++
    num3 := arr[index]
    index++
    num4 := arr[index]
    index++
    num := uint32(num4)<<24 | uint32(num3)<<16 | uint32(num2)<<8 | uint32(num1)
    return num, index
}

func CheckHeader(hdr byte, chk byte) bool {
	if hdr != chk {
	    fmt.Fprintf(os.Stderr, "Header was 0x%x instead of 0x%x\n", hdr, chk)
		return false
	} else {
		return true
	}
}

func main() {
	A2S_INFO := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x54, 0x53, 0x6F, 0x75, 0x72, 0x63, 0x65, 0x20, 0x45, 0x6E, 0x67, 0x69, 0x6E, 0x65, 0x20, 0x51, 0x75, 0x65, 0x72, 0x79, 0x00}
    CHALLENGE_RESPONSE := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x54, 0x53, 0x6F, 0x75, 0x72, 0x63, 0x65, 0x20, 0x45, 0x6E, 0x67, 0x69, 0x6E, 0x65, 0x20, 0x51, 0x75, 0x65, 0x72, 0x79, 0x00, 0x00, 0x00, 0x00, 0x00}
    A2S_RULES := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x56, 0xFF, 0xFF, 0xFF, 0xFF}
    A2S_PLAYER := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x55, 0xFF, 0xFF, 0xFF, 0xFF}
    
	argsWithProg := os.Args
	if len(argsWithProg) < 3 {
		fmt.Printf("Usage: %s <server> <port>\n", filepath.Base(argsWithProg[0]))
		os.Exit(1)
	}

	server := argsWithProg[1]
	port := argsWithProg[2]
	if len(argsWithProg) > 3 {
		debug = true
	}
	if len(argsWithProg) > 4 {
		colorize = true
	}
    seconds := 15
    timeout := time.Duration(seconds) * time.Second

    if(debug) {	fmt.Fprintln(os.Stderr, "Opening UDP connection...") }
	Conn, err := net.DialTimeout("udp", server+":"+port, timeout)
	if !CheckNoError(err) {
        os.Exit(2)
    }

	defer Conn.Close()

    // Get Info

    if(debug) {	fmt.Fprintln(os.Stderr, "Sending A2S_INFO...") }
    
    start := time.Now()
	n, BytesReceived := SendPacket(Conn, A2S_INFO, timeout)
    t := time.Now()
    elapsed1 := t.Sub(start)
    
    if BytesReceived == nil || n == 0 {
        fmt.Fprintln(os.Stderr, "Received no data!")
        os.Exit(2)
    }
    
    if CheckHeader(BytesReceived[4], 0x41) {
        CHALLENGE_RESPONSE[25] = BytesReceived[5]
        CHALLENGE_RESPONSE[26] = BytesReceived[6]
        CHALLENGE_RESPONSE[27] = BytesReceived[7]
        CHALLENGE_RESPONSE[28] = BytesReceived[8]
        n, BytesReceived = SendPacket(Conn, CHALLENGE_RESPONSE, timeout)
    }

    if !CheckHeader(BytesReceived[4], 0x49) { os.Exit(2) }

    if(debug) {	fmt.Fprintf(os.Stderr, "HEADER: 0x%x\n", BytesReceived[4]) }
    if(debug) {	fmt.Fprintf(os.Stderr, "PROTOCOL: 0x%x\n", BytesReceived[5]) }
    
    var sPtr int
	var info string
    sPtr = 5
    info, sPtr = GetString(BytesReceived, sPtr)
    fmt.Printf("NAME: %s\n", info)
    
    info, sPtr = GetString(BytesReceived, sPtr)
    fmt.Printf("MAP: %s\n", info)

    info, sPtr = GetString(BytesReceived, sPtr)
    fmt.Printf("FOLDER: %s\n", info)

    info, sPtr = GetString(BytesReceived, sPtr)
    fmt.Printf("GAME: %s\n", info)

	var id uint16
    id, sPtr = GetUInt16(BytesReceived, sPtr)
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
			var port uint16
            port, _ = GetUInt16(BytesReceived, sPtr)
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

    if(debug) {	fmt.Fprintln(os.Stderr, "Sending A2S_RULES...") }
    
    start = time.Now()
	n, BytesReceived = SendPacket(Conn, A2S_RULES, timeout)
    t = time.Now()
    elapsed2 := t.Sub(start)

    if BytesReceived == nil || n == 0 {
        fmt.Fprintln(os.Stderr, "Received no data!")
        os.Exit(2)
    }

    if !CheckHeader(BytesReceived[4], 0x41) { os.Exit(2) }

    // Challenge number
	var chnum uint32
    chnum, sPtr = GetUInt32(BytesReceived, sPtr)
    if(debug) {	fmt.Fprintf(os.Stderr,"Challenge number: %d\n", chnum) }

    A2S_RULES[5] = byte(chnum)
    A2S_RULES[6] = byte(chnum >> 8)
    A2S_RULES[7] = byte(chnum >> 16)
    A2S_RULES[8] = byte(chnum >> 24)

    if(debug) {	fmt.Fprintln(os.Stderr, "Sending A2S_RULES...") }
    
    start = time.Now()
	n, BytesReceived = SendPacket(Conn, A2S_RULES, timeout)
    t = time.Now()
    elapsed3 := t.Sub(start)

    if BytesReceived == nil || n == 0 {
        fmt.Fprintln(os.Stderr, "Received no data!")
        os.Exit(2)
    }
    
    elapsed := (elapsed1 + elapsed2 + elapsed3) / 3

    fmt.Printf("PING: %d\n", int(elapsed) / 1000000)
    
    if !CheckHeader(BytesReceived[4], 0x45) { os.Exit(2) }

    // reset sPtr
    sPtr = 5
	var rules uint16
    rules, sPtr = GetUInt16(BytesReceived, sPtr)

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

    if(debug) {	fmt.Fprintln(os.Stderr, "Sending A2S_PLAYER...") }
	n, BytesReceived = SendPacket(Conn, A2S_PLAYER, timeout)

    if BytesReceived == nil || n == 0 {
        fmt.Fprintln(os.Stderr, "Received no data!")
        os.Exit(2)
    }

    if !CheckHeader(BytesReceived[4], 0x41) { os.Exit(2) }

    // Challenge number
    chnum, sPtr = GetUInt32(BytesReceived, sPtr)
    if(debug) {	fmt.Fprintf(os.Stderr,"Challenge number: %d\n", chnum) }

    A2S_PLAYER[5] = byte(chnum)
    A2S_PLAYER[6] = byte(chnum >> 8)
    A2S_PLAYER[7] = byte(chnum >> 16)
    A2S_PLAYER[8] = byte(chnum >> 24)

    if(debug) {	fmt.Fprintln(os.Stderr, "Sending A2S_PLAYER...") }
	n, BytesReceived = SendPacket(Conn, A2S_PLAYER, timeout)

    if BytesReceived == nil || n == 0 {
        fmt.Fprintln(os.Stderr, "Received no data!")
        os.Exit(2)
    }

    if !CheckHeader(BytesReceived[4], 0x44) { os.Exit(2) }

    sPtr = 5
    players := BytesReceived[sPtr]
    sPtr++

    if players > 0 {
        fmt.Println("PLAYER LIST:");
    }

	var score uint32

    for i := 0; i < int(players); i++ {
        // Index (this seems to always be 0, so skipping it)
        sPtr++

        // Name
        info, sPtr = GetString(BytesReceived, sPtr)

        // Score
        score, sPtr = GetUInt32(BytesReceived, sPtr)

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
