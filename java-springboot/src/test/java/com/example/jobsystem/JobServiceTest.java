package com.example.jobsystem;

import com.example.jobsystem.model.*;
import com.example.jobsystem.repository.*;
import com.example.jobsystem.service.JobService;
import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

class JobServiceTest {
    @Test void createsQueuedJobsWithTypedInputsThenRetrievesAndDeletes() {
        var service = new JobService(new InMemoryJobRepository());
        var email = new EmailInput("learner@example.com", "hello");
        var job = service.createJob("send_email", email);
        assertEquals(1, job.id);
        assertEquals("queued", job.status);
        assertEquals("send_email", job.jobType);
        assertEquals(email, service.getJob(1).input);
        var meeting = service.createJob("classify_meeting", new MeetingInput("m", "summary"));
        assertEquals(2, meeting.id);
        assertEquals("queued", meeting.status);
        service.delete(1);
        assertThrows(IllegalStateException.class, () -> service.getJob(1));
        assertThrows(IllegalStateException.class, () -> service.delete(1));
    }

    @Test void rejectsWrongTypesAndEveryMissingField() {
        var service = new JobService(new InMemoryJobRepository());
        Object[][] cases = {
            {"other", new EmailInput("a", "b")}, {"send_email", null},
            {"send_email", new MeetingInput("m", "s")},
            {"classify_meeting", new EmailInput("a", "b")},
            {"send_email", new EmailInput("", "b")}, {"send_email", new EmailInput("a", "")},
            {"classify_meeting", new MeetingInput("", "s")}, {"classify_meeting", new MeetingInput("m", "")}
        };
        for (var c : cases) assertThrows(IllegalArgumentException.class, () -> service.createJob((String)c[0], c[1]));
        assertThrows(IllegalArgumentException.class, () -> service.getJob(0));
        assertThrows(IllegalArgumentException.class, () -> service.delete(-1));
        assertEquals(1, service.createJob("send_email", new EmailInput(" ", " ")).id);
    }

    @Test void preservesRepositoryFailureAsCause() {
        var repo = mock(JobRepository.class);
        var failure = new IllegalStateException("storage failure");
        when(repo.save(any())).thenThrow(failure);
        when(repo.get(1)).thenThrow(failure);
        doThrow(failure).when(repo).delete(1);
        var service = new JobService(repo);
        assertSame(failure, assertThrows(IllegalStateException.class,
                () -> service.createJob("send_email", new EmailInput("a", "b"))).getCause());
        assertSame(failure, assertThrows(IllegalStateException.class, () -> service.getJob(1)).getCause());
        assertSame(failure, assertThrows(IllegalStateException.class, () -> service.delete(1)).getCause());
    }
}
