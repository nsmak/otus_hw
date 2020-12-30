package app

type Event struct {
	ID          string
	Title       string
	StartDate   int64
	EndDate     int64
	Description string
	OwnerID     string
	RemindIn    int64
}
