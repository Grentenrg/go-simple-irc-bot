package main

import (
	"fmt"
	"log/slog"
	"net"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gobs/pretty"
)

type handler func(Msg) error

type Channel struct {
	name            string
	modes           string
	Messages        []string
	Topic           string
	TopicChangeTime time.Time
	TopicChangeBy   string
	Users           []*UserIdentity

	shouldResetNames bool
}

func NewChannel(name string) *Channel {
	return &Channel{name: name, Messages: make([]string, 0)}
}

func (c *Channel) AddMessage(msg string) {
	c.Messages = append(c.Messages, msg)
}

type Client struct {
	conn     net.Conn
	logger   *slog.Logger
	handlers map[string]handler

	motd []string

	me UserIdentity

	userMods string
	channels map[string]*Channel
}

func (c *Client) setupHandlers() {
	c.handlers = map[string]handler{
		"RPL_WELCOME":      c.HandleRPL_WELCOME,
		"PING":             c.HandlePing,
		"RPL_MOTD":         c.HandleRPL_MOTD,
		"RPL_ENDOFMOTD":    c.HandleRPL_ENDOFMOTD,
		"RPL_UMODEIS":      c.HandleRPL_UMODEIS,
		"RPL_MOTDSTART":    c.HandleRPL_MOTDSTART,
		"JOIN":             c.HandleJOIN,
		"RPL_NAMREPLY":     c.HandleRPL_NAMREPLY,
		"RPL_ENDOFNAMES":   c.HandleRPL_ENDOFNAMES,
		"PRIVMSG":          c.HandlePRIVMSG,
		"RPL_TOPIC":        c.HandleRPL_TOPIC,
		"PART":             c.HandlePART,
		"QUIT":             c.HandleQUIT,
		"NICK":             c.HandleNICK,
		"MODE":             c.HandleMODE,
		"KICK":             c.HandleKICK,
		"TOPIC":            c.HandleRPL_TOPIC,
		"RPL_NOTOPIC":      c.HandleRPL_NOTOPIC,
		"RPL_TOPICWHOTIME": c.HandleRPL_TOPICWHOTIME,
	}
}

func NewClient(conn net.Conn, logger *slog.Logger) *Client {
	c := &Client{
		conn:   conn,
		logger: logger, channels: make(map[string]*Channel),
	}
	c.setupHandlers()

	c.motd = make([]string, 0)

	return c
}

func (c *Client) Write(data []byte) (int, error) {
	return c.conn.Write(data)
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// send data to server
func (c *Client) Send(data []byte) (int, error) {
	c.logger.Debug("->", "data", string(data))
	data = append(data, '\r', '\n')
	return c.Write(data)
}

func (c *Client) Read(data []byte) (int, error) {
	return c.conn.Read(data)
}

// Handle RPL_MOTD
func (c *Client) HandleRPL_MOTD(msg Msg) error {
	c.motd = append(c.motd, strings.Join(msg.Args[0:], " "))
	return nil
}

// Handle RPL_ENDOFMOTD
func (c *Client) HandleRPL_ENDOFMOTD(msg Msg) error {
	// join #gral.irc
	if _, err := c.Send([]byte("JOIN #gral.irc")); err != nil {
		return fmt.Errorf("error sending join: %w", err)
	}
	return nil
}

// Handle RPL_UMODEIS
func (c *Client) HandleRPL_UMODEIS(msg Msg) error {
	c.userMods = msg.Args[1]

	c.me = UserIdentity{Nick: msg.Args[0]}

	return nil
}

func (c *Client) HandleRPL_MOTDSTART(msg Msg) error {
	c.motd = make([]string, 0)
	return nil
}

func (c *Client) HandleRPL_WELCOME(msg Msg) error {
	c.logger.Info("WELCOME")
	return nil
}

// send PRIVMSG
func (c *Client) SendPRIVMSG(target, message string) error {
	if _, err := c.Send([]byte("PRIVMSG " + target + " :" + message)); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

func (c *Client) HandlePRIVMSG(msg Msg) error {
	target := msg.Target

	if target[0] != '#' {
		// Private message
		return nil
	} else {
		if channel, ok := c.channels[target]; ok {
			channel.AddMessage(msg.Raw)

			message := strings.Join(msg.Args[1:], " ")

			c.logger.Info("message", "channel", target, "message", message, "from", msg.Nick)

			words := strings.Split(message, " ")
			if words[0] == "!topic" {
				if err := c.SendPRIVMSG(target, "Topic: "+channel.Topic); err != nil {
					return fmt.Errorf("error sending topic: %w", err)
				}
			} else if words[0] == "!users" {
				users := make([]string, 0)
				for _, user := range channel.Users {
					users = append(users, user.Nick)
				}
				if err := c.SendPRIVMSG(target, "Users: "+strings.Join(users, ", ")); err != nil {
					return fmt.Errorf("error sending names: %w", err)
				}
			}
		} else {
			c.logger.Error("channel not found", "channel", target)
		}
		return nil
	}
}

func (c *Client) Handle(msg Msg) error {
	handler, ok := c.handlers[msg.CommandName()]
	if !ok {
		return fmt.Errorf("unknown command: %s: %w", msg.CommandName(), ErrUnknwonCommand)
	}

	pretty.PrettyPrint(msg)
	if err := handler(msg); err != nil {
		return fmt.Errorf(
			"error handling message command:%s|%s: %w",
			msg.Command,
			msg.CommandName(),
			err,
		)
	}

	return nil
}

func (c *Client) HandlePing(msg Msg) error {
	if _, err := c.Send([]byte("PONG " + msg.Args[0])); err != nil {
		return fmt.Errorf("error sending pong: %w", err)
	}

	return nil
}

// topics

func (c *Client) HandleRPL_TOPIC(msg Msg) error {
	channel := msg.Args[1]
	if _, ok := c.channels[channel]; !ok {
		c.channels[channel] = NewChannel(channel)
		c.logger.Error("channel not found", "channel", channel)
	}
	c.channels[channel].Topic = msg.Trailing
	c.channels[channel].TopicChangeTime = time.Now()

	return nil
}

func (c *Client) HandleRPL_NOTOPIC(msg Msg) error {
	channel := msg.Target
	if _, ok := c.channels[channel]; !ok {
		c.logger.Error("channel not found", "channel", channel)
	}
	c.channels[channel].Topic = ""
	c.channels[channel].TopicChangeTime = time.Now()

	return nil
}

// handle 333 RPL_TOPICWHOTIME
func (c *Client) HandleRPL_TOPICWHOTIME(msg Msg) error {
	channel := msg.Args[1]
	if _, ok := c.channels[channel]; !ok {
		c.logger.Error("channel not found", "channel", channel)
	}

	unixTimestamp := msg.Args[3]
	intUnixTimestamp, err := strconv.ParseInt(unixTimestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing unix timestamp: %w", err)
	}
	t := time.Unix(intUnixTimestamp, 0)
	c.channels[channel].TopicChangeTime = t.In(time.FixedZone("UTC", 0))

	c.channels[channel].TopicChangeBy = msg.Args[2]

	return nil
}

// JOIN
func (c *Client) Join(channel string) error {
	if _, err := c.Send([]byte("JOIN " + channel)); err != nil {
		return fmt.Errorf("error sending join: %w", err)
	}
	return nil
}

// Handle JOIN
func (c *Client) HandleJOIN(msg Msg) error {
	user := UserIdentity{Nick: msg.Nick, User: msg.User, Host: msg.Host}

	channel := msg.Target
	if _, ok := c.channels[channel]; !ok {
		if user.Nick != c.me.Nick {
			c.logger.Error("channel not found", "channel", channel)
		}
		c.channels[channel] = NewChannel(channel)
	}

	c.channels[channel].Users = append(c.channels[channel].Users, &user)

	if user.Nick == c.me.Nick {
		c.logger.Info("joined channel", "channel", channel)
	}

	c.channels[channel].shouldResetNames = true

	return nil
}

// handle RPL_NAMEREPLY
func (c *Client) HandleRPL_NAMREPLY(msg Msg) error {
	channel := msg.Args[2]
	if _, ok := c.channels[channel]; !ok {
		c.logger.Error("channel not found", "channel", channel)
		c.channels[channel] = NewChannel(channel)
	}

	if c.channels[channel].shouldResetNames {
		c.channels[channel].Users = make([]*UserIdentity, 0)
		c.channels[channel].shouldResetNames = false
	}

	users := strings.Fields(msg.Trailing)

	for _, user := range users {
		c.channels[channel].Users = append(
			c.channels[channel].Users,
			&UserIdentity{Nick: strings.TrimLeft(user, "@+")},
		)
	}
	return nil
}

// handle RPL_ENDOFNAMES
func (c *Client) HandleRPL_ENDOFNAMES(msg Msg) error {
	channel := msg.Args[1]
	if _, ok := c.channels[channel]; !ok {
		c.logger.Error("channel not found", "channel", channel)
		c.channels[channel] = NewChannel(channel)
	}

	c.channels[channel].shouldResetNames = true
	return nil
}

// handle PART
func (c *Client) HandlePART(msg Msg) error {
	channel := msg.Target
	if _, ok := c.channels[channel]; !ok {
		c.logger.Error("channel not found", "channel", channel)
	}

	user := UserIdentity{Nick: msg.Nick, User: msg.User, Host: msg.Host}

	for i, u := range c.channels[channel].Users {
		if u.Nick == user.Nick {
			c.channels[channel].Users = slices.Delete(c.channels[channel].Users, i, 1)
			break
		}
	}

	if user.Nick == c.me.Nick {
		c.logger.Info("left channel", "channel", channel)
		delete(c.channels, channel)
	}

	return nil
}

// handle QUIT
func (c *Client) HandleQUIT(msg Msg) error {
	user := UserIdentity{Nick: msg.Nick, User: msg.User, Host: msg.Host}

	for _, channel := range c.channels {
		for i, u := range channel.Users {
			if u.Nick == user.Nick {
				channel.Users = slices.Delete(channel.Users, i, 1)
				break
			}
		}
	}

	return nil
}

// handle NICK
func (c *Client) HandleNICK(msg Msg) error {
	user, err := ParseUserIdentity(msg.Nick)
	if err != nil {
		return fmt.Errorf("error parsing user identity: %w", err)
	}

	newNick := msg.Args[0]
	for _, channel := range c.channels {
		for i, u := range channel.Users {
			if u.Nick == user.Nick {
				channel.Users[i].Nick = newNick
			}
		}
	}

	if user.Nick == c.me.Nick {
		c.me.Nick = newNick
	}

	return nil
}

// handle MODE
func (c *Client) HandleMODE(msg Msg) error {
	if msg.Target[0] != '#' {
		// User mode
		return nil
	} else {
		channel := msg.Target
		if ch, ok := c.channels[channel]; !ok {
			c.logger.Error("channel not found", "channel", channel)
		} else {
			ch.modes = msg.Args[0]
		}

	}
	return nil
}

// handle KICK
func (c *Client) HandleKICK(msg Msg) error {
	channel := msg.Target
	if _, ok := c.channels[channel]; !ok {
		c.logger.Error("channel not found", "channel", channel)
	}

	// user := UserIdentity{Nick: msg.Nick, User: msg.User, Host: msg.Host}
	targettedUser := UserIdentity{Nick: msg.Args[1]}

	for i, u := range c.channels[channel].Users {
		if u.Nick == targettedUser.Nick {
			c.channels[channel].Users = slices.Delete(c.channels[channel].Users, i, 1)
			break
		}
	}

	pretty.PrettyPrint(c.channels)

	return nil
}

// send NICK
func (c *Client) SendNICK(nick string) error {
	if _, err := c.Send([]byte("NICK " + nick)); err != nil {
		return fmt.Errorf("error sending nick: %w", err)
	}
	return nil
}

// send PART
func (c *Client) SendPART(channel string) error {
	if _, err := c.Send([]byte("PART " + channel)); err != nil {
		return fmt.Errorf("error sending part: %w", err)
	}
	return nil
}

// send QUIT
func (c *Client) SendQUIT() error {
	if _, err := c.Send([]byte("QUIT")); err != nil {
		return fmt.Errorf("error sending quit: %w", err)
	}
	return nil
}

// send KICK
func (c *Client) SendKICK(channel, nick, reason string) error {
	if _, err := c.Send([]byte("KICK " + channel + " " + nick + " :" + reason)); err != nil {
		return fmt.Errorf("error sending kick: %w", err)
	}
	return nil
}

// send MODE
func (c *Client) SendMODE(channel, mode string) error {
	if _, err := c.Send([]byte("MODE " + channel + " " + mode)); err != nil {
		return fmt.Errorf("error sending mode: %w", err)
	}
	return nil
}

// send PASS
func (c *Client) Pass(password string) error {
	if _, err := c.Send([]byte("PASS " + password)); err != nil {
		return fmt.Errorf("error sending pass: %w", err)
	}
	return nil
}

// send USER
func (c *Client) User(username, realname string) error {
	if _, err := c.Send([]byte("USER " + username + " ignored ignored :" + realname)); err != nil {
		return fmt.Errorf("error sending user: %w", err)
	}
	return nil
}

// send NICK
func (c *Client) Nick(nick string) error {
	if _, err := c.Send([]byte("NICK " + nick)); err != nil {
		return fmt.Errorf("error sending nick: %w", err)
	}
	return nil
}
