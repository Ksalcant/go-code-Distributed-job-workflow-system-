package service

/*
	Wich operation is valid in our system?

Is this a supported job type?
What status should a new job have?
Is the input valid for this operation?
Is this state transition allowed?
*/
import (
	"errors"
	"fmt"
	"job-system/model"
	"job-system/repository"
)

type JobService struct {
	repo repository.JobRepository
}

func NewJobService(repo repository.JobRepository) *JobService {

	return &JobService{
		repo: repo,
	}
}
func (srv *JobService) CreateJob(jobType string, jobInput any) (model.Job, error) {

	if jobInput == nil {
		return model.Job{}, errors.New("Job input must not be nil")
	}

	switch jobType {
	case "send_email":
		emailInput, ok := jobInput.(model.EmailInput)
		if !ok {
			return model.Job{}, errors.New("job input must be EmailInput")
		}
		if emailInput.Email == "" {
			return model.Job{}, errors.New("email must not be empty")
		}
		if emailInput.Message == "" {
			return model.Job{}, errors.New("message must not be empty")
		}
	case "classify_meeting":
		meetingInput, ok := jobInput.(model.MeetingInput)
		if !ok {
			return model.Job{}, errors.New("job input must be meetingInput")
		}
		if meetingInput.MeetingID == "" {
			return model.Job{}, errors.New("meetingID must not be empty")
		}
		if meetingInput.Summary == "" {
			return model.Job{}, errors.New("summary must not be empty")
		}

	default:
		return model.Job{}, errors.New("unsupported job type")
	}

	job := &model.Job{
		Status:  "queued",
		JobType: jobType,
		Input:   jobInput,
	}
	savedJob, err := srv.repo.Save(job)
	if err != nil {
		return model.Job{}, fmt.Errorf("create job: %w", err) // propagate
	}
	return savedJob, nil

}
