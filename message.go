package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Message struct {
	Name    string
	Command string
	Params  []string
}

func (m Message) String() string {
	return fmt.Sprintf("Name: %s, Command: %s, Params: %v", m.Name, m.Command, m.Params)
}

func (m Message) CommandName() string {
	if name, ok := Commands[m.Command]; ok {
		return name
	}

	return m.Command
}

func parseMessage(msg string) (Message, error) {
	if len(msg) == 0 {
		return Message{}, errors.New("empty message")
	}

	if msg[0] == ':' {
		spl := strings.SplitN(msg, " ", 3)
		if len(spl) != 3 {
			return Message{}, errors.New("invalid message")
		}

		return Message{
			Name:    spl[0][1:],
			Command: spl[1],
			Params:  strings.Split(spl[2], " "),
		}, nil
	} else {
		spl := strings.SplitN(msg, " ", 2)
		if len(spl) != 2 {
			return Message{}, errors.New("invalid message")
		}

		return Message{
			Command: spl[0],
			Params:  strings.Split(spl[1], " "),
		}, nil
	}
}

var ErrEmptyPacket = errors.New("empty packet")

func parsePacketOld(in []byte) ([]string, error) {
	if len(in) == 0 {
		return nil, ErrEmptyPacket
	}

	packets := make([]string, 0)

	for i := 0; i < len(in); {
		if i+1 >= len(in) {
			break
		}

		if in[i+1] == '\n' && in[i] == '\r' {
			packets = append(packets, string(in[0:i]))
			in = in[i+2:]
			i = 0
			continue
		}

		i++
	}

	rest := string(in)
	if len(rest) > 0 {
		packets = append(packets, rest)
	}

	return packets, nil
}

func parsePacket(in []byte) ([]string, string, error) {
	rest := ""
	if len(in) == 0 {
		return nil, rest, ErrEmptyPacket
	}

	// Preallocate slice based on estimated packet count
	estimatedPackets := bytes.Count(in, []byte{'\r', '\n'})
	packets := make([]string, 0, estimatedPackets)

	// Split but keep non-empty segments only
	for len(in) > 0 {
		idx := bytes.Index(in, []byte{'\r', '\n'})
		if idx == -1 {
			// No more CRLF, append remaining if not empty
			if len(in) > 0 {
				rest = string(in)
			}
			break
		}

		// Append segment if not empty
		if idx > 0 {
			packets = append(packets, string(in[:idx]))
		}

		// Move past the CRLF
		in = in[idx+2:]
	}

	return packets, rest, nil
}

type UserIdentity struct {
	Nick string
	User string
	Host string
}

func ParseUserIdentity(in string) (UserIdentity, error) {
	if len(in) == 0 {
		return UserIdentity{}, errors.New("empty name")
	}

	spl := strings.SplitN(in, "!", 2)
	if len(spl) != 2 {
		return UserIdentity{}, errors.New("invalid name")
	}

	spl2 := strings.SplitN(spl[1], "@", 2)
	if len(spl2) != 2 {
		return UserIdentity{}, errors.New("invalid name")
	}

	return UserIdentity{
		Nick: spl[0],
		User: spl2[0],
		Host: spl2[1],
	}, nil
}

type Msg struct {
	Raw      string            // Raw message
	Prefix   string            // Message prefix (sender)
	Nick     string            // Nickname part of prefix
	User     string            // Username part of prefix
	Host     string            // Hostname part of prefix
	Command  string            // IRC command or numeric
	Target   string            // Message target (channel or user)
	Args     []string          // Command arguments
	Trailing string            // Trailing argument (after :)
	Tags     map[string]string // IRCv3 message tags
}

func ParseMessage(line string) (*Msg, error) {
	if len(line) == 0 {
		return nil, errors.New("empty message")
	}

	msg := &Msg{
		Raw:  line,
		Tags: make(map[string]string),
	}

	// Remove \r\n if present
	line = strings.TrimRight(line, "\r\n")

	// Parse IRCv3 tags if present
	if strings.HasPrefix(line, "@") {
		if idx := strings.Index(line, " "); idx != -1 {
			tags := line[1:idx]
			line = line[idx+1:]

			// Parse tags
			for _, tag := range strings.Split(tags, ";") {
				if kv := strings.SplitN(tag, "=", 2); len(kv) == 2 {
					msg.Tags[kv[0]] = kv[1]
				} else if len(kv) == 1 {
					msg.Tags[kv[0]] = ""
				}
			}
		}
	}

	// Parse prefix if present
	if strings.HasPrefix(line, ":") {
		if idx := strings.Index(line, " "); idx != -1 {
			msg.Prefix = line[1:idx]
			line = line[idx+1:]

			// Parse prefix parts (nick!user@host)
			if idx := strings.Index(msg.Prefix, "!"); idx != -1 {
				msg.Nick = msg.Prefix[:idx]
				if idx2 := strings.Index(msg.Prefix[idx+1:], "@"); idx2 != -1 {
					msg.User = msg.Prefix[idx+1 : idx+1+idx2]
					msg.Host = msg.Prefix[idx+1+idx2+1:]
				}
			} else if idx := strings.Index(msg.Prefix, "@"); idx != -1 {
				msg.Nick = msg.Prefix[:idx]
				msg.Host = msg.Prefix[idx+1:]
			} else {
				msg.Nick = msg.Prefix
			}
		}
	}

	// Parse command and arguments
	parts := strings.SplitN(line, " :", 2)
	args := strings.Fields(parts[0])

	if len(args) == 0 {
		return nil, errors.New("no command found")
	}

	msg.Command = strings.ToUpper(args[0])

	if len(args) > 1 {
		msg.Args = args[1:]
		// Set target for common commands
		switch msg.Command {
		case "PRIVMSG", "NOTICE", "JOIN", "PART", "MODE", "TOPIC", "INVITE", "KICK":
			if len(msg.Args) > 0 {
				msg.Target = msg.Args[0]
			}
		}
	}

	// Add trailing argument if present
	if len(parts) > 1 {
		msg.Trailing = parts[1]
		msg.Args = append(msg.Args, msg.Trailing)
	}

	return msg, nil
}

// Msg.CommandName
func (m *Msg) CommandName() string {
	if name, ok := Commands[m.Command]; ok {
		return name
	}

	return m.Command
}
