package models

import (
	"TM_chat/utils"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	UserId     int64
	TargetId   int64
	Type       int `gorm:"size:8"`
	Media      int `gorm:"size:1"`
	Content    string
	CreateTime int64
	ReadTime   int64
	pic        string
	Url        string `gorm:"size:128"`
	Desc       string `gorm:"size:128"`
	Amount     int    `gorm:"size:10"`
}

func init() {
	go udpsendproc()
	go udprecvproc()
	fmt.Println("init goroutine ")
}
func (table *Message) MessageTableName() string {
	return "Message_basics"
}

type Node struct {
	Conn          *websocket.Conn //连接
	Addr          string          //客户端地址
	FirstTime     uint64          //首次连接时间
	HeartbeatTime uint64          //心跳时间
	LoginTime     uint64          //登录时间
	DataQueue     chan []byte     //消息
	GroupSets     set.Interface   //好友 / 群
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
	currentTime := uint64(time.Now().Unix())
	node := &Node{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(),
		HeartbeatTime: currentTime,
		LoginTime:     currentTime,
		DataQueue:     make(chan []byte, 50),
		GroupSets:     set.New(set.ThreadSafe),
	}
	rwlocker.Lock()
	clientMap[userId] = node
	rwlocker.Unlock()
	go sendProc(node)
	go recvProc(node)
	SetUserOnlineInfo("online_"+id, []byte(node.Addr), time.Duration(viper.GetInt("timeout.RedisOnlineTime"))*time.Hour)
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
		msg := &Message{}
		err = json.Unmarshal(data, msg)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if msg.Type == 3 {
			currentTime := uint64(time.Now().Unix())
			node.Heartbeat(currentTime)
		} else {
			dispatch(data)
			broadMsg(data)
			fmt.Println("[ws]<<<<<<", string(data))
		}

	}
}

func (node *Node) Heartbeat(time uint64) {
	node.HeartbeatTime = time
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(msg []byte) {
	udpsendChan <- msg
}

func udpsendproc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IP(viper.GetString("port.url")),
		Port: viper.GetInt("port.udp"),
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
		IP:   net.IP(viper.GetString("port.url")),
		Port: viper.GetInt("port.udp"),
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
	msg := &Message{}
	err := json.Unmarshal(buf, msg)
	msg.CreateTime = time.Now().Unix()
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
		SendGroupMsg(msg.TargetId, buf)

	}
}

func SendGroupMsg(targetId int64, msg []byte) {
	fmt.Println("群发消息开始")
	userIds := SearchUserByGroupId(uint(targetId))
	for _, userId := range userIds {
		if targetId != int64(userId) {
			SendMsg(int64(userId), msg)
		}
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
	jsonMsg := Message{}
	json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	targetIdstr := strconv.Itoa(int(userId))
	userIdstr := strconv.Itoa(int(jsonMsg.UserId))
	r, err := utils.REDIS.Get(ctx, "online_"+userIdstr).Result()
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	if r != "" {
		fmt.Println("sendMsg >>> userId:", userId, "  Msg:", string(msg))
		node.DataQueue <- msg
	}
	var key string
	if userId > jsonMsg.UserId {
		key = "msg_" + userIdstr + "_" + targetIdstr
	} else {
		key = "msg_" + targetIdstr + "_" + userIdstr
	}
	res, err := utils.REDIS.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println("ZRevRange fail, err :", err)
		return
	}
	score := float64(cap(res)) + 1
	ress, err := utils.REDIS.ZAdd(ctx, key, &redis.Z{score, msg}).Result()
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	fmt.Println("res :", ress)

}
func (msg Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(msg)
}
func JoinGroup(userId uint, comId string) int {
	contact := Contact{}
	contact.OwnerId = userId
	//contact.TargetId = comId
	contact.Type = 2
	v, _ := strconv.Atoi(comId)
	s := uint(v)
	contact.TargetId = s
	community := Community{}

	utils.DB.Where("id=? or name=?", comId, comId).Find(&community)
	if community.Name == "" {
		return -1
	}

	type Contact struct {
		gorm.Model
		OwnerId  uint
		TargetId uint
		Type     int
		Desc     string
	}
	utils.DB.Where("owner_id=? and target_id=? and type =2 ", userId, s).Find(&contact)
	if !contact.CreatedAt.IsZero() {
		return -1
	} else {

		utils.DB.Create(&contact)
		return 0
	}
}

func RedisMsg(userIdA int64, userIdB int64, start int64, end int64, isRev bool) []string {
	rwlocker.RLock()
	node, ok := clientMap[userIdA]
	rwlocker.RUnlock()
	ctx := context.Background()
	userIdstr := strconv.Itoa(int(userIdA))
	targetIdstr := strconv.Itoa(int(userIdB))
	var key string
	if userIdA > userIdB {
		key = "msg_" + targetIdstr + "_" + userIdstr
	} else {
		key = "msg_" + userIdstr + "_" + targetIdstr
	}
	var rels []string
	var err error
	if !isRev {
		rels, err = utils.REDIS.ZRange(ctx, key, start, end).Result()
	} else {
		rels, err = utils.REDIS.ZRevRange(ctx, key, start, end).Result()
	}
	if err != nil {
		fmt.Println("err :", err)
		return nil
	}
	if ok {
		for _, val := range rels {
			node.DataQueue <- []byte(val)
			fmt.Println("cleanConnection userIdA:", userIdA, "  Msg:", val)
		}
	} else {
		fmt.Println("User not found")
		return nil
	}
	return rels
}

func CleanConnection(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("cleanConnection err", r)
		}
	}()

	currentTime := uint64(time.Now().Unix())
	for i := range clientMap {
		node := clientMap[i]
		if node.IsHeartbeatTimeOut(currentTime) {
			fmt.Println("心跳超时..... 关闭连接：", node)
			node.Conn.Close()
		}
	}
	return result
}
func (node *Node) IsHeartbeatTimeOut(currentTime uint64) (timeout bool) {
	if node.HeartbeatTime+viper.GetUint64("timeout.HeartbeatMaxTime") <= currentTime {
		fmt.Println("心跳超时。。。自动下线", node)
		timeout = true
	}
	return
}
