package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

//用户上线
func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	//广播用户上线消息
	u.server.BroadCast(u, "已经上线")
}

//用户下线
func (u *User) Offline() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	//广播用户上线消息
	u.server.BroadCast(u, "已经下线")
}

//当前user对应的客户端发送消息
func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

//用户处理消息
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		//查询当前在线用户
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		//判断name是否存在
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMsg("该用户名已存在\n")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()
			u.Name = newName
			u.SendMsg("用户名成功更新为：" + u.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式

		//获取用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			u.SendMsg("私聊格式错误，请使用\"to|aa|nihao\"格式\n")
		}
		//得到user对象
		remoteUser, ok := u.server.OnlineMap[remoteName]
		if !ok {
			u.SendMsg("该用户不存在")
			return
		}
		//发送消息
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMsg("无消息内容，请重发\n")
		}
		remoteUser.SendMsg(u.Name + "私聊您：" + content)
	} else {
		u.server.BroadCast(u, msg)
	}

}

//监听当前user channel 的方法，一旦有消息直接发送给客户端
func (u *User) ListenMessage() {
	for msg := range u.C {
		//msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
