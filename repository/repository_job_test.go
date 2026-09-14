package repository

import (
	"reflect"
	"testing"

	model "job-system/model"
)

func TestSaveTableDriven(t *testing.T) {
	tests := []struct {
		name    string
		job     *model.Job
		wantID  int64
		wantErr string
	}{
		{
			"Save email job",
			&model.Job{
				Status:  "queued",
				JobType: "Send Email",
				Input: model.EmailInput{
					Email:   "kaisel@uber.com",
					Message: "Hola",
				},
			},
			1,
			"",
		},
		{
			"Save Meeting job",
			&model.Job{
				Status:  "queued",
				JobType: "classify meeting",
				Input: model.MeetingInput{
					MeetingID: "123445678",
					Summary:   "I am a summary from meeting input",
				},
			},
			1,
			"",
		},
		{
			"Save null job",
			nil,
			0,
			"Job must not be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewJobRepository()
			savedJob, err := repo.Save(tt.job)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got none", tt.wantErr)
				}

				if err.Error() != tt.wantErr {
					t.Errorf("error = %q; want %q", err.Error(), tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			storedJob, storedJobErr := repo.Get(savedJob.ID)
			if storedJobErr != nil {
				t.Fatalf("expected stored job ID %d, but got error %v", savedJob.ID, storedJobErr)
			}

			if !reflect.DeepEqual(storedJob, savedJob) {
				t.Errorf("storedJobID = %d; savedJobID = %d", storedJob.ID, savedJob.ID)
			}

			if savedJob.ID != tt.wantID {
				t.Errorf("Job ID = %d; want %d", savedJob.ID, tt.wantID)
			}
		})
	}
}

func TestGetTableDriven(t *testing.T) {
	repo := NewJobRepository()

	storedJob, err := repo.Save(&model.Job{
		Status:  "queued",
		JobType: "Send Email",
		Input: model.EmailInput{
			Email:   "kaisel@uber.com",
			Message: "Hola",
		},
	})
	if err != nil {
		t.Fatalf("setup save failed: %v", err)
	}

	tests := []struct {
		name    string
		jobID   int64
		want    model.Job
		wantErr string
	}{
		{"Job less than 0", -1, model.Job{}, "Job ID must be greater than 0"},
		{"Job with ID = 0", 0, model.Job{}, "Job ID must be greater than 0"},
		{"Job not found", 999, model.Job{}, "JobID not found"},
		{"Return valid job", storedJob.ID, storedJob, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repo.Get(test.jobID)

			if test.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got none", test.wantErr)
				}

				if err.Error() != test.wantErr {
					t.Fatalf("error = %q; want %q", err.Error(), test.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("got = %#v; want = %#v", got, test.want)
			}
		})
	}
}
