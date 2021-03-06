package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int
	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

//server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
		//
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

//监听message消息channel的go程，一旦有消息发送给所有在线用户
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message

		//msg发送个全部在线user
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

//广播消息的方法
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	//
	//fmt.Println("链接建立成功")
	//用户上线,将用户加入表
	user := NewUser(conn, s)

	user.Online()
	isLive := make(chan bool)
	// s.mapLock.Lock()
	// s.OnlineMap[user.Name] = user
	// s.mapLock.Unlock()
	// //广播用户上线消息
	// s.BroadCast(user, "已经上线")
	//接受客户端发的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			//提取用户的消息去除换行
			msg := string(buf[:n-1])
			//广播消息
			user.DoMessage(msg)
			isLive <- true
		}
	}()
	for {
		select {
		//用户活跃
		case <-isLive:
		//已超时强制登出
		case <-time.After(time.Second * 300):
			user.SendMsg("由于长时间不活跃您已经被踢")
			close(user.C)
			conn.Close()
			return
		}
	}

}

//启动服务器的接口
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()
	//启动监听message的
	go s.ListenMessager()
	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//do handler
		go s.Handler(conn)
	}

}
