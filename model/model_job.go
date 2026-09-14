package model

type Job struct {
	ID      int64
	Status  string
	JobType string
	Input   any
}

type EmailInput struct {
	Email   string
	Message string
}

type MeetingInput struct {
	MeetingID string
	Summary   string
}
