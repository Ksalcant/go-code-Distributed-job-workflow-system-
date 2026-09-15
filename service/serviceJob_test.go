package service

import (
	"job-system/repository"
	"testing"

	"job-system/model"
)

func TestCreateJob_Email(t *testing.T) {
	repo := repository.NewJobRepository()
	service := NewJobService(repo)

	input := model.EmailInput{
		Email:   "kaisel@uber.com",
		Message: "Hola",
	}

	job, err := service.CreateJob("send_email", input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if job.ID != 1 {
		t.Errorf("job ID = %d; want 1", job.ID)
	}

	if job.Status != "queued" {
		t.Errorf("job status = %q; want %q", job.Status, "queued")
	}

	if job.JobType != "send_email" {
		t.Errorf("job type = %q; want %q", job.JobType, "send_email")
	}

	emailInput, ok := job.Input.(model.EmailInput)

	if !ok {
		t.Fatalf("job input type =%T; want model.EmailInput", job.Input)
	}
	if emailInput.Email != input.Email {
		t.Errorf("email %q; want %q", emailInput.Email, input.Email)
	}
	if emailInput.Message != input.Message {
		t.Errorf("message = %q; want %q", emailInput.Message, input.Message)
	}

}

func TestCreateJob_UnsupportedJobType(t *testing.T) {
	repo := repository.NewJobRepository()
	service := NewJobService(repo)

	input := model.EmailInput{
		Email:   "kaisel@uber.com",
		Message: "Hola",
	}

	_, err := service.CreateJob("something_else", input)
	if err == nil {
		t.Fatalf("expected unsupported error")
	}

}

func TestCreateJob_WrongInputType(t *testing.T) {
	repo := repository.NewJobRepository()
	service := NewJobService(repo)

	input := model.MeetingInput{
		MeetingID: "12323343545463645",
		Summary:   "Test1",
	}

	_, err := service.CreateJob("send_email", input)
	if err == nil {
		t.Fatalf("expected wrong inputtype error")
	}

}
func TestCreateJob_MissingTypeFields(t *testing.T) {
	repo := repository.NewJobRepository()
	service := NewJobService(repo)

	input := model.EmailInput{
		Email:   "",
		Message: "Hola",
	}

	_, err := service.CreateJob("send_email", input)
	if err == nil {
		t.Fatalf("expected email must not be empty error")
	}

}
