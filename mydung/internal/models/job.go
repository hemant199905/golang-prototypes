package models

type Job struct {
	ID      string
	Name    string
	Payload map[string]string
	Status  string
	FailureReason string
}
