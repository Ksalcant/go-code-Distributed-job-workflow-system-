package com.example.jobsystem;

import com.example.jobsystem.model.*;
import com.example.jobsystem.repository.InMemoryJobRepository;
import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;

class JobRepositoryTest {
    @Test void savesBothPayloadsAndAssignsIncreasingIds() {
        var repo = new InMemoryJobRepository();
        var email = new Job(0, "queued", "send_email", new EmailInput("a", "b"));
        var saved = repo.save(email);
        assertEquals(1, saved.id);
        assertEquals(1, email.id);
        assertEquals(saved.input, repo.get(1).input);
        assertEquals(2, repo.save(new Job(0, "queued", "classify_meeting", new MeetingInput("m", "s"))).id);
    }

    @Test void rejectsNullInvalidAndMissingJobs() {
        var repo = new InMemoryJobRepository();
        assertThrows(IllegalArgumentException.class, () -> repo.save(null));
        assertThrows(IllegalArgumentException.class, () -> repo.update(null));
        for (long id : new long[]{-1, 0, 999}) {
            assertThrows(IllegalArgumentException.class, () -> repo.get(id));
            assertThrows(IllegalArgumentException.class, () -> repo.delete(id));
            assertThrows(IllegalArgumentException.class, () -> repo.update(new Job(id, "queued", "send_email", null)));
        }
    }

    @Test void copiesMetadataReplacesRecordsAndDoesNotReuseDeletedIds() {
        var repo = new InMemoryJobRepository();
        var original = new Job(0, "queued", "send_email", new EmailInput("a", "b"));
        var saved = repo.save(original);
        original.status = "changed";
        saved.status = "also changed";
        assertEquals("queued", repo.get(1).status);
        var fetched = repo.get(1);
        fetched.status = "replacement";
        assertEquals("queued", repo.get(1).status);
        repo.update(fetched);
        fetched.status = "changed again";
        assertEquals("replacement", repo.get(1).status);
        repo.delete(1);
        assertThrows(IllegalArgumentException.class, () -> repo.get(1));
        assertEquals(2, repo.save(original).id);
    }
}
