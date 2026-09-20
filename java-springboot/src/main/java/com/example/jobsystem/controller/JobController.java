package com.example.jobsystem.controller;

import com.example.jobsystem.model.*;
import com.example.jobsystem.service.JobService;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/jobs")
public class JobController {
    private final JobService service;
    private final ObjectMapper mapper;

    public JobController(JobService service, ObjectMapper mapper) {
        this.service = service;
        this.mapper = mapper;
    }

    public record CreateJobRequest(@JsonProperty("job_type") String jobType, JsonNode input) {}

    @PostMapping
    public ResponseEntity<?> create(@RequestBody(required = false) String body) {
        CreateJobRequest request;
        try {
            request = mapper.readValue(body == null ? "" : body, CreateJobRequest.class);
        } catch (Exception error) {
            return error(400, "invalid request body");
        }
        String jobType = request == null ? "" : request.jobType();
        if (!"send_email".equals(jobType) && !"classify_meeting".equals(jobType))
            return error(400, "job type not supported");
        Object input;
        try {
            if (request.input() == null) throw new IllegalArgumentException("missing input");
            // JSON null becomes a zero-value payload in Go, then fails validation.
            if (request.input().isNull()) {
                input = jobType.equals("send_email") ? new EmailInput(null, null) : new MeetingInput(null, null);
            } else if (jobType.equals("send_email")) {
                input = mapper.treeToValue(request.input(), EmailInput.class);
            } else {
                input = mapper.treeToValue(request.input(), MeetingInput.class);
            }
        } catch (Exception error) {
            return error(400, jobType.equals("send_email") ? "Bad job input" : "bad job input");
        }
        Job job;
        try {
            job = service.createJob(jobType, input);
        } catch (RuntimeException error) {
            return error(400, "unable to create job");
        }
        return json(201, job, "Internal error");
    }

    @GetMapping("/{id}")
    public ResponseEntity<?> get(@PathVariable String id) {
        long jobId;
        try {
            jobId = parseId(id);
            if (jobId <= 0) throw new NumberFormatException();
        } catch (NumberFormatException error) {
            return error(400, "Invalid job ID");
        }
        Job job;
        try {
            job = service.getJob(jobId);
        } catch (RuntimeException error) {
            return error(404, "Job not found");
        }
        return json(200, job, "Internal Server Error");
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<?> delete(@PathVariable String id) {
        long jobId;
        try {
            jobId = parseId(id);
        } catch (NumberFormatException error) {
            return error(400, "invalid job ID");
        }
        try {
            service.delete(jobId);
        } catch (RuntimeException error) {
            // Preserve Go's current DELETE behavior: missing/invalid IDs return 400.
            return error(400, "service error");
        }
        return ResponseEntity.noContent().header("Content-Type", "application/json").build();
    }

    private static long parseId(String id) {
        if (!id.matches("[+-]?[0-9]+")) throw new NumberFormatException();
        return Long.parseLong(id);
    }

    private ResponseEntity<?> json(int status, Job job, String failure) {
        try {
            return ResponseEntity.status(status).contentType(MediaType.APPLICATION_JSON)
                    .body(mapper.writeValueAsString(job));
        } catch (Exception error) {
            return error(500, failure);
        }
    }

    private static ResponseEntity<String> error(int status, String message) {
        return ResponseEntity.status(status).header("Content-Type", "text/plain; charset=utf-8")
                .header("X-Content-Type-Options", "nosniff").body(message + "\n");
    }
}
