package app

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	StartDate   int64  `json:"start_date" db:"start_date"`
	EndDate     int64  `json:"end_date" db:"end_date"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id" db:"owner_id"`
	RemindIn    int64  `json:"remind_in" db:"remind_in"`
}
