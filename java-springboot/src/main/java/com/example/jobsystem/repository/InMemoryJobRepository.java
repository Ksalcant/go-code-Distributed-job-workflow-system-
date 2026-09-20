package com.example.jobsystem.repository;

import com.example.jobsystem.model.Job;
import org.springframework.stereotype.Repository;
import java.util.HashMap;
import java.util.Map;

@Repository
public class InMemoryJobRepository implements JobRepository {
    // Same checkpoint as Go: synchronization remains a future milestone.
    private final Map<Long, Job> jobs = new HashMap<>();
    private long nextId = 1;

    @Override
    public Job save(Job job) {
        if (job == null) throw new IllegalArgumentException("Job must not be nil");
        job.id = nextId++;
        jobs.put(job.id, job.copy());
        return job.copy();
    }

    @Override
    public Job get(long jobId) {
        if (jobId <= 0) throw new IllegalArgumentException("Job ID must be greater than 0");
        Job job = jobs.get(jobId);
        if (job == null) throw new IllegalArgumentException("JobID not found");
        return job.copy();
    }

    @Override
    public void update(Job job) {
        if (job == null) throw new IllegalArgumentException("Job must not be nil");
        if (job.id <= 0) throw new IllegalArgumentException("JobID must not be negative");
        if (!jobs.containsKey(job.id)) throw new IllegalArgumentException("Job id does not exists!");
        jobs.put(job.id, job.copy());
    }

    @Override
    public void delete(long jobId) {
        if (jobId <= 0) throw new IllegalArgumentException("job ID must be greater than 0");
        if (!jobs.containsKey(jobId)) throw new IllegalArgumentException("job not found");
        jobs.remove(jobId);
    }
}
