package repository

import (
	"errors"

	"job-system/model"
)

type JobRepository interface {
	Save(job *model.Job) (model.Job, error)
	Update(job *model.Job) error
	Delete(jobID int64) error
	Get(jobID int64) (model.Job, error)
}

type jobRepository struct {
	jobs   map[int64]model.Job
	nextID int64
}

func NewJobRepository() *jobRepository {
	return &jobRepository{
		jobs:   make(map[int64]model.Job),
		nextID: 1,
	}
}

func (repo *jobRepository) Save(job *model.Job) (model.Job, error) {
	if job == nil {
		return model.Job{}, errors.New("Job must not be nil")
	}

	job.ID = repo.nextID
	repo.jobs[job.ID] = *job
	repo.nextID++
	return *job, nil
}

func (repo *jobRepository) Get(jobID int64) (model.Job, error) {
	if jobID <= 0 {
		return model.Job{}, errors.New("Job ID must be greater than 0")
	}

	if job, found := repo.jobs[jobID]; found {
		return job, nil
	}

	return model.Job{}, errors.New("JobID not found")
}

func (repo *jobRepository) Update(job *model.Job) error {
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
