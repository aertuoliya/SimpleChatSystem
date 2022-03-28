package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}

//单独处理server的回执
func (client *Client) DealResponse() {
	//client.conn有数据则拷贝到stdout标准输出上
	io.Copy(os.Stdout, client.conn)

}

//菜单
func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公共聊天")
	fmt.Println("2.私聊")
	fmt.Println("3.更改用户名")
	fmt.Println("0.退出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法的命令\n")
		return false
	}
}

//公共聊天
func (client *Client) PublicChat() {
	//用户输入信息
	var chatMsg string
	fmt.Printf("请输入聊天内容(或输入exit退出)：")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		//发给服务器
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println("请输入聊天内容(或输入exit退出)：")
		fmt.Scanln(&chatMsg)
	}

}

//查询在线的用户
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

//私聊
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println("请输入您要私聊的在线用户（或输入exit退出）：")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("请输入聊天内容（或输入exit退出）：")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println("请输入聊天内容（或输入exit退出）：")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUsers()
		fmt.Println("请输入您要私聊的在线用户（或输入exit退出）：")
		fmt.Scanln(&remoteName)
	}
}

//更改用户名
func (client *Client) UpdateName() bool {
	fmt.Printf("请输入用户名：")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

//菜单选项的运行
func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}
		switch client.flag {
		case 1:
			//fmt.Println("公共聊天中..")
			client.PublicChat()
		case 2:
			//fmt.Println("私聊中..")
			client.PrivateChat()
		case 3:
			//fmt.Println("更改用户名中..")
			client.UpdateName()

		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口（默认是8888）")
}
func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	//client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("链接服务器失败")
		return
	}
	//单独一个goroutine处理server的回执
	go client.DealResponse()
	fmt.Println("链接服务器成功")
	//启动客户端业务
	client.Run()
	//select {}
}
