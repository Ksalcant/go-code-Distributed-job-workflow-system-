package com.example.jobsystem.repository;

import com.example.jobsystem.model.Job;

public interface JobRepository {
    Job save(Job job);
    Job get(long jobId);
    void update(Job job);
    void delete(long jobId);
}
