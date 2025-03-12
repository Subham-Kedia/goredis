package main

import (
	"fmt"
	"log/slog"
	"net"
)

const defaultListenAddress = ":5001"

type Config struct {
	ListenAddress string
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitCh    chan struct{}
	msgCh     chan []byte
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
		quitCh:    make(chan struct{}),
		msgCh:     make(chan []byte),
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
		case rawMsg := <-s.msgCh:
			if err := s.handleRawMessage(rawMsg); err != nil {
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

func (s *Server) handleRawMessage(rawMsg []byte) error {
	slog.Info("message recieved", "msg", string(rawMsg))
	cmd, err := parseMessage(string(rawMsg))
	if err != nil {
		return err
	}
	switch c := cmd.(type) {
	case *SetCommand:
		s.kv.Set(c.key, c.value)
	case *GetCommand:
		val, ok := s.kv.Get(c.key)
		if ok {
			fmt.Println(string(val))
		}
	}
	return nil
}
