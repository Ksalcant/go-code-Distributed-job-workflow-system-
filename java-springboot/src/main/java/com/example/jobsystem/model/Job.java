package com.example.jobsystem.model;

import com.fasterxml.jackson.annotation.JsonProperty;

// Public fields keep this learning model close to the original Go struct.
public class Job {
    @JsonProperty("ID") public long id;
    @JsonProperty("Status") public String status;
    @JsonProperty("JobType") public String jobType;
    @JsonProperty("Input") public Object input;

    public Job(long id, String status, String jobType, Object input) {
        this.id = id;
        this.status = status;
        this.jobType = jobType;
        this.input = input;
    }

    // Go stores/returns Job values. Copy metadata to avoid Java reference aliasing.
    public Job copy() {
        return new Job(id, status, jobType, input);
    }
}
