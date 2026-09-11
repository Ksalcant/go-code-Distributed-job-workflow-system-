package main

import "testing"

func Save_Test_Driven_Table(t *testing.T) {
	repo := NewJobRepository()
	tests := []struct {
		name    string
		status  string
		jobType string
		input   any
		want    any
	}{
		{
			"Save email job", "queued", "Send Email", EmailInput{
				Email:   "kaisel@uber.com",
				Message: "Hola"}, 1,
		},
		{
			"Save Meeting job", "queued", "classify meeting", MeetingInput{
				MeetingID: "123445678",
				Summary:   "I am a summary from meeting input",
			}, 2,
		},
		{
			"Save null job", "", "nil job", nil, "Job must not be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job, err := repo.Save(&Job{
				Status:  tt.status,
				JobType: tt.jobType,
				Input:   tt.input,
			})
			if err != nil && err != tt.want {
				t.Errorf("Save(nil job) = %s; want %s", err, tt.want)
			}
			if job.ID != tt.want {
				t.Errorf("Job ID = %d; want %d", job.ID, tt.want)
			}

		})
	}

}
