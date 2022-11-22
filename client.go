package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		if len(cmd) > 0 && cmd[0] == '/' {
			switch cmd {
			case "/nick":
				c.commands <- command{
					id:     CMD_NICK,
					client: c,
					args:   args,
				}
			case "/join":
				c.commands <- command{
					id:     CMD_JOIN,
					client: c,
					args:   args,
				}
			case "/rooms":
				c.commands <- command{
					id:     CMD_ROOMS,
					client: c,
				}
			case "/help":
				c.commands <- command{
					id:     CMD_HELP,
					client: c,
				}
			case "/quit":
				c.commands <- command{
					id:     CMD_QUIT,
					client: c,
				}
			default:
				c.err(fmt.Errorf("unknown command: %s", cmd))
			}
		} else {
			if c.room != nil {
				args = append([]string{cmd}, args...)

				c.commands <- command{
					id:     CMD_MSG,
					client: c,
					args:   args,
				}
			}
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
