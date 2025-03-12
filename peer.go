package main

import (
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
}

func NewPeer(conn net.Conn, msgCh chan Message) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
	}
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case CommandSet:
					if len(v.Array()) < 3 {
						continue
					}
					p.msgCh <- Message{
						cmd: &SetCommand{
							key:   v.Array()[1].String(),
							value: v.Array()[2].String(),
						},
						peer: p,
					}
				case CommandGet:
					if len(v.Array()) < 2 {
						continue
					}
					p.msgCh <- Message{
						cmd: &GetCommand{
							key: v.Array()[1].String(),
						},
						peer: p,
					}
				}
			}
		}
	}
	return nil
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}
