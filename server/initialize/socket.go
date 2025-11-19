package initialize

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// 用户信息结构
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// WebSocket 消息结构
type Message struct {
	Type    string      `json:"type"`    // message, join, leave, etc.
	From    User        `json:"from"`    // 发送者信息
	Content interface{} `json:"content"` // 消息内容
	Time    int64       `json:"time"`    // 时间戳
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	user   *User
	userID string
}

type Hub struct {
	clients    map[*Client]bool
	users      map[string]*Client // 用户ID到客户端的映射
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		users:      make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			if client.userID != "" {
				h.users[client.userID] = client
			}
			log.Printf("用户 %s 连接，当前连接数: %d", client.userID, len(h.clients))

			// 广播用户上线消息
			joinMsg := Message{
				Type: "user_join",
				From: *client.user,
				Content: map[string]interface{}{
					"userCount":   len(h.clients),
					"onlineUsers": h.GetOnlineUsers(),
				},
				Time: time.Now().Unix(),
			}
			h.BroadcastMessage(joinMsg)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if client.userID != "" {
					delete(h.users, client.userID)
				}
				close(client.send)
				log.Printf("用户 %s 断开，当前连接数: %d", client.userID, len(h.clients))

				// 广播用户下线消息
				leaveMsg := Message{
					Type: "user_leave",
					From: *client.user,
					Content: map[string]interface{}{
						"userCount":   len(h.clients),
						"onlineUsers": h.GetOnlineUsers(),
					},
					Time: time.Now().Unix(),
				}
				h.BroadcastMessage(leaveMsg)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					if client.userID != "" {
						delete(h.users, client.userID)
					}
				}
			}
		}
	}
}

// 获取在线用户列表
func (h *Hub) GetOnlineUsers() []User {
	users := make([]User, 0, len(h.users))
	for _, client := range h.users {
		if client.user != nil {
			users = append(users, *client.user)
		}
	}
	return users
}

// 向特定用户发送消息
func (h *Hub) SendToUser(userID string, message Message) {
	if client, exists := h.users[userID]; exists {
		msgBytes, _ := json.Marshal(message)
		client.send <- msgBytes
	}
}

// 广播消息
func (h *Hub) BroadcastMessage(message Message) {
	msgBytes, _ := json.Marshal(message)
	for client := range h.clients {
		select {
		case client.send <- msgBytes:
		default:
			close(client.send)
			delete(h.clients, client)
			if client.userID != "" {
				delete(h.users, client.userID)
			}
		}
	}
}

// Token 验证函数（示例）
func validateToken(token string) (*User, error) {
	// 这里应该是你的实际 Token 验证逻辑
	// 示例：简单的基于 token 前缀的验证
	if strings.HasPrefix(token, "user_") {
		return &User{
			ID:       token,
			Username: "用户_" + token[5:],
			Role:     "user",
		}, nil
	} else if strings.HasPrefix(token, "admin_") {
		return &User{
			ID:       token,
			Username: "管理员_" + token[6:],
			Role:     "admin",
		}, nil
	}
	return nil, http.ErrAbortHandler
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// 从查询参数获取 token
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token required", http.StatusUnauthorized)
		return
	}

	// 验证 token
	user, err := validateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket 升级失败: %v", err)
		return
	}

	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		user:   user,
		userID: user.ID,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("用户 %s 读取错误: %v", c.userID, err)
			}
			break
		}

		// 解析消息
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("消息解析错误: %v", err)
			continue
		}

		// 设置发送者信息
		msg.From = *c.user
		msg.Time = time.Now().Unix()

		log.Printf("收到来自 %s 的消息: %s", c.user.Username, msg.Content)

		// 处理不同类型的消息
		switch msg.Type {
		case "chat":
			// 广播聊天消息
			msgBytes, _ := json.Marshal(msg)
			c.hub.broadcast <- msgBytes

		case "private_message":
			// 私聊消息
			if targetUserID, ok := msg.Content.(map[string]interface{})["target"].(string); ok {
				privateMsg := Message{
					Type: "private_message",
					From: *c.user,
					Content: map[string]interface{}{
						"message": msg.Content.(map[string]interface{})["message"],
						"from":    c.user.Username,
					},
					Time: time.Now().Unix(),
				}
				c.hub.SendToUser(targetUserID, privateMsg)

				// 同时给自己也发一份
				c.hub.SendToUser(c.userID, privateMsg)
			}

		default:
			// 广播其他类型消息
			msgBytes, _ := json.Marshal(msg)
			c.hub.broadcast <- msgBytes
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// API: 获取在线用户列表
func onlineUsersHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := hub.GetOnlineUsers()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"onlineUsers": users,
			"userCount":   len(users),
		})
	}
}

// API: 向特定用户发送消息（管理员功能）
func sendMessageHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			UserID  string      `json:"user_id"`
			Message interface{} `json:"message"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		msg := Message{
			Type:    "system_message",
			From:    User{ID: "system", Username: "系统", Role: "system"},
			Content: request.Message,
			Time:    time.Now().Unix(),
		}

		hub.SendToUser(request.UserID, msg)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "message_sent",
		})
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok", "service": "websocket-server"}`))
}

func InitWebsocket() {
	hub := NewHub()
	go hub.Run()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/api/online-users", onlineUsersHandler(hub))
	http.HandleFunc("/api/send-message", sendMessageHandler(hub))

	port := ":8081"
	log.Printf("WebSocket 服务器启动在 http://localhost%s", port)
	log.Printf("WebSocket 端点: ws://localhost%s/ws?token=user_123", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("服务器启动失败: ", err)
	}
}
