package phone

import (
	"math/rand"
	"sync"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)
import "gorm.io/gorm"

// Service represents the phone service
type Service struct {
	DB *gorm.DB
}

// WebSocketHandler represents a WebSocket handler
type WebSocketHandler struct {
    clients map[*websocket.Conn]bool
    broadcast chan Phone
    upgrader websocket.Upgrader
}
// NewWebSocketHandler creates a new WebSocketHandler instance
func NewWebSocketHandler() *WebSocketHandler {
    return &WebSocketHandler{
        clients:   make(map[*websocket.Conn]bool),
        broadcast: make(chan Phone),
        upgrader: websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 1024,
        },
    }
}

// StartWebSocket starts the WebSocket handler
func (h *WebSocketHandler) StartWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := h.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()

    h.clients[conn] = true

    for {
        var msg Phone
        err := conn.ReadJSON(&msg)
        if err != nil {
            log.Println(err)
            delete(h.clients, conn)
            break
        }

        h.broadcast <- msg
    }
}

// BroadcastWebSocket sends updates to all connected clients
func (h *WebSocketHandler) BroadcastWebSocket() {
    for {
        msg := <-h.broadcast
        for client := range h.clients {
            err := client.WriteJSON(msg)
            if err != nil {
                log.Println(err)
                client.Close()
                delete(h.clients, client)
            }
        }
    }
}
// Phone represents a phone number
type Phone struct {
	ID    uint `gorm:"primaryKey;autoIncrement" json:"id"`
	Number string `json:"number"`
}

// PhoneList represents a list of phone numbers
type PhoneList struct {
	mu     sync.Mutex
	Phones []Phone
}

// PaginatedPhonesResponse represents the response structure for paginated phone numbers
type PaginatedPhonesResponse struct {
	TotalData   int      `json:"total_data"`
	PageNow     int      `json:"page_now"`
	NextPage    string   `json:"next_page,omitempty"`
	TotalPages  int      `json:"total_pages"`
	PrevPage    string   `json:"prev_page,omitempty"`
	PhoneNumbers []Phone `json:"phone_numbers"`
}

// GetPaginatedPhones retrieves a paginated list of phone numbers
func (s *Service) GetPaginatedPhones(pageNum, pageSize int) (PaginatedPhonesResponse, error) {
	var result PaginatedPhonesResponse

	// Count total number of records
	var totalData int64
	if err := s.DB.Model(&Phone{}).Count(&totalData).Error; err != nil {
		return result, err
	}

	// Calculate total pages
	totalPages := int((totalData + int64(pageSize) - 1) / int64(pageSize))

	// Calculate offset
	offset := (pageNum - 1) * pageSize

	// Retrieve paginated phone numbers
	var phones []Phone
	if err := s.DB.Offset(offset).Limit(pageSize).Find(&phones).Error; err != nil && err != gorm.ErrRecordNotFound {
		return result, err
	}

	// Set response metadata
	result.TotalData = int(totalData)
	result.PageNow = pageNum
	result.TotalPages = totalPages
	result.PhoneNumbers = phones

	// Set next page link if applicable
	if pageNum < totalPages {
		result.NextPage = fmt.Sprintf("/phones?page=%d&pageSize=%d", pageNum+1, pageSize)
	}

	// Set previous page link if applicable
	if pageNum > 1 {
		result.PrevPage = fmt.Sprintf("/phones?page=%d&pageSize=%d", pageNum-1, pageSize)
	}

	return result, nil
}
// GetPhones retrieves the list of phone numbers from the database
func (s *Service) GetPhones() ([]Phone, error) {
	var phones []Phone
	if err := s.DB.Find(&phones).Error; err != nil {
		return nil, err
	}
	return phones, nil
}

// AddPhone adds a new phone number to the database
func (s *Service) AddPhone(newPhone *Phone) error {
	return s.DB.Create(newPhone).Error
}

// UpdatePhone updates the information of a phone number in the database
func (s *Service) UpdatePhone(id string, updatedPhone *Phone) error {
	return s.DB.Where("id = ?", id).Updates(updatedPhone).Error
}

// DeletePhone deletes a phone number from the database
func (s *Service) DeletePhone(id string) error {
	return s.DB.Where("id = ?", id).Delete(&Phone{}).Error
}
// generateID generates a simple unique ID (for demonstration purposes)
func generateID() string {
	// You may use a more sophisticated ID generation mechanism in a real-world scenario
	return "id-" + randomString(6)
}

// randomString generates a random string of the specified length
func randomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[randomInt(len(chars))]
	}
	return string(result)
}

// randomInt generates a random integer in the range [0, max)
func randomInt(max int) int {
	return rand.Intn(max)
}
