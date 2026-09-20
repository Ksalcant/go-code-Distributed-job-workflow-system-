# Job System — Java / Spring Boot track

A Java 21 / Spring Boot 3.5.16 translation of the Go implementation at commit `71a9a88` (`Add delete method`), inspected September 19, 2026. The original Go implementation stays at the repository root. This directory is an independent Maven project in the same Git repository.

## Exact stopping point

- `POST /jobs`: accept `send_email` or `classify_meeting`, validate typed input, assign an increasing ID starting at 1, and save with status `queued`.
- `GET /jobs/{id}`: retrieve an existing job.
- `DELETE /jobs/{id}`: delete an existing job.
- Repository `save`, `get`, `update`, and `delete`; no HTTP update or list endpoint.
- In-memory storage only. Restarting clears jobs and resets IDs. Jobs are never executed.

Workers, queues, scheduling, retries, workflows, persistence, classification providers, and observability remain future work. The map and counter are deliberately unsynchronized to preserve the Go learning checkpoint; use sequential requests. Spring's server can handle concurrent requests, so synchronization still needs to be learned and added before relying on concurrent access.

## Run and test

Install a JDK 21 and put its `java` on PATH (or set `JAVA_HOME` to that JDK). The Maven wrapper downloads Maven 3.9.11 on first use; no global Maven installation is required. Dependencies also need internet access on the first build.

From the repository root:

```sh
cd java-springboot
./mvnw test
./mvnw spring-boot:run
```

On Windows use `mvnw.cmd`. Java listens on **8081**, while Go remains on **8080**, so you can compare both without port conflicts. Stop with Ctrl+C. Optionally override Java's port with `SERVER_PORT=8091 ./mvnw spring-boot:run`.

To build a runnable jar:

```sh
./mvnw package
java -jar target/job-system-0.0.1-SNAPSHOT.jar
```

## Try the same contract

```sh
curl -i http://localhost:8081/jobs -H 'Content-Type: application/json' \
  -d '{"job_type":"send_email","input":{"Email":"learner@example.com","Message":"Hello"}}'
curl -i http://localhost:8081/jobs/1
curl -i -X DELETE http://localhost:8081/jobs/1
curl -i http://localhost:8081/jobs -H 'Content-Type: application/json' \
  -d '{"job_type":"classify_meeting","input":{"MeetingID":"meeting-001","Summary":"Project roadmap"}}'
```

On a fresh instance the first response is `201`:

```json
{"ID":1,"Status":"queued","JobType":"send_email","Input":{"Email":"learner@example.com","Message":"Hello"}}
```

GET succeeds with `200` and the same body; DELETE succeeds with `204` and no body. The envelope uses `job_type` and `input`; response and payload names match Go's exported fields. Case-insensitive field matching and ignored unknown fields match Go. Required fields must be nonempty strings; whitespace is accepted and email address syntax is not checked. Numeric/boolean payload fields are rejected instead of silently converted to strings.

| Failure | Status / body |
| --- | --- |
| Malformed request JSON | 400 / `invalid request body` |
| Unknown job type | 400 / `job type not supported` |
| Wrong/missing input shape | 400 / `Bad job input` (email) or `bad job input` (meeting) |
| Required field missing/empty or creation service failure | 400 / `unable to create job` |
| GET malformed, nonpositive, or overflowing ID | 400 / `Invalid job ID` |
| GET missing job/service failure | 404 / `Job not found` |
| DELETE malformed/overflowing ID | 400 / `invalid job ID` |
| DELETE nonpositive/missing ID/service failure | 400 / `service error` |

Controller error bodies are plain text with a trailing newline. Spring supplies framework-level responses for unmatched routes and unsupported methods; their bodies/headers need not match Go's router. GET `/jobs` and PUT `/jobs/1` return 405. This is application behavior parity, not a replacement for the Go HTTP server/parser in every protocol edge case.

## Read the layers side by side

| Go | Java under `src/main/java/com/example/jobsystem/` |
| --- | --- |
| `main.go` wiring | `JobSystemApplication` and Spring constructor injection |
| `model/model_job.go` | `model/Job`, `EmailInput`, `MeetingInput` |
| Repository interface and map implementation | `repository/JobRepository`, `InMemoryJobRepository` |
| `service/jobService.go` | `service/JobService` |
| `controller/handler/job_handler.go` | `controller/JobController` |
| `json.RawMessage` | Jackson `JsonNode` before payload conversion |
| `(value, error)` and `%w` | Return value or exception with preserved cause |
| `map[int64]Job` | `HashMap<Long, Job>` with explicit shallow metadata copies |

Spring scans `@Repository`, `@Service`, and `@RestController` and supplies their constructor dependencies. `@JsonProperty` preserves the existing API casing. Immutable payload records make the two implemented payloads easy to compare and share. The mutable Job model assigns IDs just like Go's pointer argument; repository copies prevent Java references from accidentally changing stored metadata without `update`.

Tests under `src/test/java` cover repository behavior, service validation/error causes, and requests through Spring's HTTP mappings using MockMvc. These extra parity checks validate the translation; they do not advance the system's feature milestone or count as independently completed learning exercises.

## Shared learning log

Continue using **`../learningLog.md`**, with separate Go history and Java track sections. The log was already intentionally ignored by Git and remains local; there is no second competing log. The Java section explains language/framework equivalents and records this assisted port without claiming independent mastery. If the log is absent in a fresh clone, that is the existing local-only policy.

Verified September 19, 2026: `./mvnw package` passed all 10 Java tests, existing Go tests passed, and a temporary live comparison matched 34 request status/body cases plus two 405 status checks. Both temporary servers were stopped.

Framework requirements: [Spring Boot 3.5 documentation](https://docs.spring.io/spring-boot/3.5/system-requirements.html).
