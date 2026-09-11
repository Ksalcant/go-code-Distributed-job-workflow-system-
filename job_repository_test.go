package main

import "testing"

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
			job, err := repo.Save(tt.job)
			if err != nil && err.Error() != tt.wantErr {
				t.Errorf("Save(nil job) = %s; want %s", err, tt.wantErr)
			}
			if job.ID != tt.wantID {
				t.Errorf("Job ID = %d; wantID %d", job.ID, tt.wantID)
			}

		})
	}

}
