package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromId   int64
	TargetId int64
	Type     int `gorm:"size:8"`
	Media    int `gorm:"size:1"`
	Content  string
	pic      string
	Url      string `gorm:"size:128"`
	Desc     string `gorm:"size:128"`
	Amount   int    `gorm:"size:10"`
}

func (table *Message) MessageTableName() string {
	return "Message_basics"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

var clientMap map[int64]*Node = make(map[int64]*Node, 0)

var rwlocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {

	query := request.URL.Query()
	id := query.Get("userId")
	userId, _ := strconv.ParseInt(id, 10, 64)
	// token := query.Get("token")
	// target := query.Get("target")
	// context := query.Get("context")
	// msgtype := query.Get("type")
	isvalide := true
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return isvalide },
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	rwlocker.Lock()
	clientMap[userId] = node
	rwlocker.Unlock()
	go sendProc(node)
	go recvProc(node)
	SendMsg(userId, []byte("欢迎进入聊天室"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws]<<<<<<", string(data))
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(msg []byte) {
	udpsendChan <- msg
}

func init() {
	go udpsendproc()
	go udprecvproc()
}

func udpsendproc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 88, 255),
		Port: 3000,
	})

	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	for {
		select {
		case msg := <-udpsendChan:
			_, err = conn.Write(msg)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func udprecvproc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 3000,
	})

	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	for {
		var buf [512]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println(err)
			return
		}

		dispatch(buf[:n])
	}
}

func dispatch(buf []byte) {
	msg := Message{}
	err := json.Unmarshal(buf, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("msg type:", msg.Type)
	switch msg.Type {
	case 1:
		fmt.Println("dispatch:", buf)
		SendMsg(msg.TargetId, buf)
	case 2:
		// SendGroupMsg()
	case 3:
		// SendAllMsg()

	case 4:

	}
}

func SendMsg(userId int64, msg []byte) {
	rwlocker.RLock()
	node, ok := clientMap[userId]
	rwlocker.RUnlock()
	if !ok {
		fmt.Println("User not found")
		return
	}
	node.DataQueue <- msg

}
