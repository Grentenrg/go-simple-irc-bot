package main

import (
	"errors"
	"log"
	"log/slog"
	"net"
	"os"
)

var (
	ErrNotCRLFTerminated = errors.New("not CRLF terminated")
	ErrUnknwonCommand    = errors.New("unknown command")
)

func main() {
	addr := "localhost:6667"
	username := "[bot]Gral-irc"
	realName := "gral.irc bot"

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	logger.Info("connecting to server", "addr", addr)

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := NewClient(conn, logger)

	exit := make(chan struct{})

	if err := client.Pass(""); err != nil {
		logger.Error("error sending password", "error", err)
		os.Exit(1)
	}

	if err := client.Nick(username); err != nil {
		logger.Error("error sending nickname", "error", err)
		os.Exit(1)
	}

	if err := client.User(username, realName); err != nil {
		logger.Error("error sending user", "error", err)
		os.Exit(1)
	}

	go func() {
		var previous []byte
		packet := ""
		for {
			data := make([]byte, 1024)

			n, err := client.Read(data)
			if err != nil {
				logger.Error("error reading from server", "error", err)
				break
			}

			// fmt.Println("<- ", string(data[:n]))
			packets, rest, err := parsePacket(data[:n])
			if err != nil {
				logger.Error("error parsing packet", "error", err)
				break
			}

			if len(previous) > 0 {
				packet = string(previous) + packets[0]
				packets[0] = packet
			}

			previous = []byte(rest)

			for _, p := range packets {
				client.logger.Debug(p)
				// msg, err := parseMessage(p)
				// if err != nil {
				// 	logger.Error("error parsing message", "error", err)
				// 	continue
				// }

				m, err := ParseMessage(p)
				if err != nil {
					logger.Error("error parsing message", "error", err)
					continue
				}

				if err = client.Handle(*m); err != nil {
					logger.Error("error handling message", "error", err, "message", m)
					continue
				}
			}

		}
		close(exit)
	}()

	<-exit
}
