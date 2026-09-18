package handler

import (
	"encoding/json"
	"job-system/model"
	"job-system/service"
	"net/http"
	"strconv"
)

type CreateJobRequest struct {
	JobType string          `json:"job_type"`
	Input   json.RawMessage `json:"input"`
}

type JobHandler struct {
	service *service.JobService
}

func NewJobHandler(srv *service.JobService) *JobHandler {
	return &JobHandler{
		service: srv,
	}
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req CreateJobRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	var cleanedInput any
	switch req.JobType {

	case "send_email":
		var emailInput model.EmailInput
		err := json.Unmarshal(req.Input, &emailInput)
		if err != nil {
			http.Error(w, "Bad job input", http.StatusBadRequest)
			return
		}
		cleanedInput = emailInput

	case "classify_meeting":
		var meetingInput model.MeetingInput
		err := json.Unmarshal(req.Input, &meetingInput)
		if err != nil {
			http.Error(w, "bad job input", http.StatusBadRequest)
			return
		}
		cleanedInput = meetingInput
	default:
		http.Error(w, "job type not supported", http.StatusBadRequest)
		return
	}

	job, err := h.service.CreateJob(req.JobType, cleanedInput)
	if err != nil {
		http.Error(w, "unable to create job", http.StatusBadRequest)
		return
	}

	responseBody, err := json.Marshal(job)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(responseBody)

	if err != nil {
		return
	}

}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	jobIdString := r.PathValue("id")
	jobID, err := strconv.ParseInt(jobIdString, 10, 64)
	if err != nil || jobID <= 0 {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}
	job, err := h.service.GetJob(jobID)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}
	rawJob, err := json.Marshal(job)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(rawJob)
	return
}
