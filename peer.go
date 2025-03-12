package main

import (
	"log/slog"
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
	// Read from the connection and send to the message channel
	// all peers have access to server message channel
	// so that they can send messages to the server
	// and the server can broadcast the message to all peers
	// we need to identify which peer sent the message
	// so that we can exclude the sender from the broadcast
	for {
		n, err := p.conn.Read(buff)
		if err != nil {
			slog.Error("read error", "err", err)
			return err
		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buff[:n])
		p.msgCh <- msgBuf
	}
}
