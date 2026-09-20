package com.example.jobsystem;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.web.servlet.MockMvc;
import com.fasterxml.jackson.databind.ObjectMapper;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

@SpringBootTest
@AutoConfigureMockMvc
class JobHttpTest {
    @Autowired MockMvc http;
    @Autowired ObjectMapper mapper;

    @Test void roundTripsBothTypesAndDeletes() throws Exception {
        for (String payload : new String[]{
                "{\"job_type\":\"send_email\",\"input\":{\"Email\":\"a\",\"Message\":\"b\"}}",
                "{\"job_type\":\"classify_meeting\",\"input\":{\"MeetingID\":\"m\",\"Summary\":\"s\"}}"}) {
            var body = http.perform(post("/jobs").contentType("application/json").content(payload))
                    .andExpect(status().isCreated()).andExpect(jsonPath("$.Status").value("queued"))
                    .andExpect(jsonPath("$.id").doesNotExist()).andReturn().getResponse().getContentAsString();
            long id = mapper.readTree(body).get("ID").asLong();
            http.perform(get("/jobs/" + id)).andExpect(status().isOk()).andExpect(content().json(body));
            http.perform(delete("/jobs/" + id)).andExpect(status().isNoContent()).andExpect(content().string(""));
            http.perform(get("/jobs/" + id)).andExpect(status().isNotFound()).andExpect(content().string("Job not found\n"));
            http.perform(delete("/jobs/" + id)).andExpect(status().isBadRequest()).andExpect(content().string("service error\n"));
        }
    }

    @Test void rejectsMalformedUnsupportedAndIncompleteRequests() throws Exception {
        for (String body : new String[]{"{", "", "[]", "{}",
                "{\"job_type\":\"other\",\"input\":{}}",
                "{\"job_type\":\"send_email\"}",
                "{\"job_type\":\"send_email\",\"input\":null}",
                "{\"job_type\":\"send_email\",\"input\":[]}",
                "{\"job_type\":\"send_email\",\"input\":{\"Email\":3,\"Message\":\"x\"}}",
                "{\"job_type\":\"send_email\",\"input\":{\"Email\":\"a\"}}",
                "{\"job_type\":\"classify_meeting\",\"input\":{\"meeting_id\":\"m\",\"Summary\":\"s\"}}"}) {
            http.perform(post("/jobs").contentType("application/json").content(body)).andExpect(status().isBadRequest());
        }
    }

    @Test void acceptsCaseInsensitiveFieldsAndIgnoresUnknownFieldsLikeGo() throws Exception {
        http.perform(post("/jobs").contentType("application/json").content(
                "{\"JOB_TYPE\":\"send_email\",\"input\":{\"email\":\" \",\"message\":\"hello\",\"extra\":1}}"))
                .andExpect(status().isCreated()).andExpect(jsonPath("$.Input.Email").value(" "));
    }

    @Test void validatesIdsAndDoesNotExposeUnimplementedEndpoints() throws Exception {
        for (String id : new String[]{"abc", "0", "-1", "9223372036854775808"}) {
            http.perform(get("/jobs/" + id)).andExpect(status().isBadRequest());
            http.perform(delete("/jobs/" + id)).andExpect(status().isBadRequest());
        }
        http.perform(get("/jobs/999999")).andExpect(status().isNotFound());
        http.perform(get("/jobs")).andExpect(status().isMethodNotAllowed());
        http.perform(put("/jobs/1")).andExpect(status().isMethodNotAllowed());
    }
}
