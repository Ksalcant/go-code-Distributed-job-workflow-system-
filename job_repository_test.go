package main

import (
	"reflect"
	"testing"
)

func TestSaveTableDriven(t *testing.T) {
	tests := []struct {
		name    string
		job     *Job
		wantID  int64
		wantErr string
	}{
		{
			"Save email job", &Job{
				Status:  "queued",
				JobType: "Send Email",
				Input: EmailInput{
					Email:   "kaisel@uber.com",
					Message: "Hola",
				},
			}, 1, "",
		},
		{
			"Save Meeting job", &Job{
				Status: "queued", JobType: "classify meeting", Input: MeetingInput{
					MeetingID: "123445678",
					Summary:   "I am a summary from meeting input",
				},
			}, 1, "",
		},
		{
			"Save null job", nil, 0, "Job must not be nil",
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
				t.Fatalf("expected stored job ID %d,  but got error %v", savedJob.ID, storedJobErr)
			}

			if !reflect.DeepEqual(storedJob, savedJob) {
				t.Errorf("storedJobID= %d; savedJobID= %d", storedJob.ID, savedJob.ID)
			}
			if savedJob.ID != tt.wantID {
				t.Errorf("Job ID = %d; want %d", savedJob.ID, tt.wantID)
			}
		})
	}

}
