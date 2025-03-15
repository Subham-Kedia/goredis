package main

import (
	"log/slog"
	"net"
)

const defaultListenAddress = ":5001"

type Config struct {
	ListenAddress string
}

type Message struct {
	cmd  Command
	peer *Peer
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitCh    chan any
	msgCh     chan Message
	kv        *KV
}

func NewServer(config Config) *Server {
	if len(config.ListenAddress) == 0 {
		config.ListenAddress = defaultListenAddress
	}
	return &Server{
		Config:    config,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan any),
		msgCh:     make(chan Message),
		kv:        NewKV(),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddress)
	if err != nil {
		return err
	}
	s.ln = ln
	go s.loop()
	slog.Info("server listening", "address", s.ListenAddress)
	return s.AcceptLoop()
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("message handling error", "err", err)
			}
		case <-s.quitCh:
			return
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
			slog.Info("peer added to server", "remoteAddr", peer.conn.RemoteAddr())
		}
	}
}

func (s *Server) AcceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("new connection accept error", "err", err)
			continue
		}
		slog.Info("connection request accepted", "remoteAddr", conn.RemoteAddr())
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer
	if err := peer.readLoop(); err != nil {
		slog.Error("read error", "err", err, "remoteAddr", conn.RemoteAddr())
	}
}

func (s *Server) handleMessage(msg Message) error {
	switch c := msg.cmd.(type) {
	case *SetCommand:
		s.kv.Set(c.key, c.value)
		msg.peer.Send([]byte("+OK\r\n"))
	case *GetCommand:
		val, ok := s.kv.Get(c.key)
		if ok {
			msg.peer.Send(val)
		} else {
			msg.peer.Send([]byte("$-1\r\n"))
		}
	}
	return nil
}
