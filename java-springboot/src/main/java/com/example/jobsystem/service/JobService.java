package com.example.jobsystem.service;

import com.example.jobsystem.model.*;
import com.example.jobsystem.repository.JobRepository;
import org.springframework.stereotype.Service;

@Service
public class JobService {
    private final JobRepository repository;

    public JobService(JobRepository repository) {
        this.repository = repository;
    }

    public Job createJob(String jobType, Object input) {
        if (input == null) throw new IllegalArgumentException("Job input must not be nil");
        switch (jobType == null ? "" : jobType) {
            case "send_email" -> {
                if (!(input instanceof EmailInput email))
                    throw new IllegalArgumentException("job input must be EmailInput");
                requireField(email.email(), "email must not be empty");
                requireField(email.message(), "message must not be empty");
            }
            case "classify_meeting" -> {
                if (!(input instanceof MeetingInput meeting))
                    throw new IllegalArgumentException("job input must be meetingInput");
                requireField(meeting.meetingId(), "meetingID must not be empty");
                requireField(meeting.summary(), "summary must not be empty");
            }
            default -> throw new IllegalArgumentException("unsupported job type");
        }
        try {
            return repository.save(new Job(0, "queued", jobType, input));
        } catch (RuntimeException error) {
            throw new IllegalStateException("create job: " + error.getMessage(), error);
        }
    }

    public Job getJob(long jobId) {
        if (jobId <= 0) throw new IllegalArgumentException("Job ID must be greater than 0");
        try {
            return repository.get(jobId);
        } catch (RuntimeException error) {
            throw new IllegalStateException("get job: " + error.getMessage(), error);
        }
    }

    public void delete(long jobId) {
        if (jobId <= 0) throw new IllegalArgumentException("Job ID does not exists");
        try {
            repository.delete(jobId);
        } catch (RuntimeException error) {
            throw new IllegalStateException("delete job: " + error.getMessage(), error);
        }
    }

    private static void requireField(String value, String message) {
        // Match Go: whitespace and non-address email strings remain accepted.
        if (value == null || value.isEmpty()) throw new IllegalArgumentException(message);
    }
}
