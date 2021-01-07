package app

type MQEventNotification struct {
	EventID string `json:"event_id"`
	Title   string `json:"title"`
	Date    int64  `json:"date"`
	UserID  string `json:"user_id"`
}

func NewMQEventNotification(event Event) MQEventNotification {
	return MQEventNotification{
		EventID: event.ID,
		Title:   event.Title,
		Date:    event.StartDate,
		UserID:  event.OwnerID,
	}
}

type MQMessage struct {
	Notif MQEventNotification
	Err   error
}
