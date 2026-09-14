package main

import "errors"

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

type JobRepository interface {
	Save(job *Job) (Job, error)
	Update(job *Job) error
	Delete(jobID int64) error
	Get(jobID int64) (Job, error)
}
type jobRepository struct {
	jobs   map[int64]Job
	nextID int64
}

func NewJobRepository() *jobRepository {
	return &jobRepository{
		jobs:   make(map[int64]Job),
		nextID: 1,
	}
}
func (repo *jobRepository) Save(job *Job) (Job, error) {
	if job == nil {
		return Job{}, errors.New("Job must not be nil")
	}
	job.ID = repo.nextID
	repo.jobs[job.ID] = *job
	repo.nextID++
	return *job, nil
}
func (repo *jobRepository) Get(jobID int64) (Job, error) {
	if jobID <= 0 {
		return Job{}, errors.New("Job ID must be greater than 0")
	}
	if job, found := repo.jobs[jobID]; found {
		return job, nil
	}
	return Job{}, errors.New("JobID not found")
}
func (repo *jobRepository) Update(job *Job) error {
	if job == nil {
		return errors.New("Job must not be nil")
	}

	if job.ID <= 0 {
		return errors.New("JobID must not be negative")
	}

	_, found := repo.jobs[job.ID]
	if !found {
		return errors.New("Job id does not exists!")
	}
	repo.jobs[job.ID] = *job
	return nil
}

func (repo *jobRepository) Delete(jobID int64) error {
	if jobID <= 0 {
		return errors.New("job ID must be greater than 0")
	}
	if _, found := repo.jobs[jobID]; !found {
		return errors.New("job not found")
	}
	delete(repo.jobs, jobID)
	return nil
}
