package com.example.jobsystem.model;

import com.fasterxml.jackson.annotation.JsonProperty;

public record EmailInput(@JsonProperty("Email") String email,
        @JsonProperty("Message") String message) {}
