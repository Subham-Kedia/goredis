package main

import (
	"net"
)

type Peer struct {
	conn  net.Conn
	msgCh chan []byte
}

func NewPeer(conn net.Conn, msgCh chan []byte) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
	}
}

func (p *Peer) readLoop() error {
	buff := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buff)
		if err != nil {
			return err
		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buff[:n])
		p.msgCh <- msgBuf
	}
}
