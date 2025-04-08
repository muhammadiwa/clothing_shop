// package models

// import "time"

// type TicketStatus string

// const (
// 	TicketStatusOpen     TicketStatus = "open"
// 	TicketStatusPending  TicketStatus = "pending"
// 	TicketStatusResolved TicketStatus = "resolved"
// 	TicketStatusClosed   TicketStatus = "closed"
// )

// type SupportTicket struct {
// 	ID          uint         `json:"id" gorm:"primaryKey"`
// 	UserID      uint         `json:"user_id"`
// 	User        User         `json:"user" gorm:"foreignKey:UserID"`
// 	OrderID     *uint        `json:"order_id"`
// 	Order       *Order       `json:"order" gorm:"foreignKey:OrderID"`
// 	Subject     string       `json:"subject" gorm:"not null"`
// 	Description string       `json:"description" gorm:"not null"`
// 	Status      TicketStatus `json:"status" gorm:"type:enum('open','pending','resolved','closed');default:'open'"`
// 	Priority    string       `json:"priority" gorm:"type:enum('low','medium','high');default:'medium'"`
// 	CreatedAt   time.Time    `json:"created_at" gorm:"autoCreateTime"`
// 	UpdatedAt   time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
// 	ClosedAt    *time.Time   `json:"closed_at"`

// 	Replies []TicketReply `json:"replies" gorm:"foreignKey:TicketID"`
// }

// type TicketReply struct {
// 	ID          uint      `json:"id" gorm:"primaryKey"`
// 	TicketID    uint      `json:"ticket_id"`
// 	UserID      uint      `json:"user_id"`
// 	User        User      `json:"user" gorm:"foreignKey:UserID"`
// 	Content     string    `json:"content" gorm:"not null"`
// 	IsFromAdmin bool      `json:"is_from_admin" gorm:"default:false"`
// 	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
// 	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
// }
