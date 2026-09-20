package com.example.jobsystem.model;

import com.fasterxml.jackson.annotation.JsonProperty;

public record MeetingInput(@JsonProperty("MeetingID") String meetingId,
        @JsonProperty("Summary") String summary) {}
