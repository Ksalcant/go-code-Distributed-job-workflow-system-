# Learning Log — Distributed Job / Workflow System

Last updated: September 12, 2026

## Goal

Build one Go backend system that accepts, tracks, and eventually executes jobs in the background. Grow it into a workflow system with priority scheduling, retries, dependencies, worker recovery, rate limiting, persistence, and a meeting-classification workload.

My learning goal is to strengthen backend engineering and Go while understanding and implementing the data structures and algorithms behind the system. I should be able to explain my design choices and write the core logic myself.

## Current checkpoint

**Week 3 — Go service foundation, in progress.** This is a milestone label from the original plan, not a claim that three calendar weeks of work are complete.

- Implemented the job model and a single-threaded, in-memory repository with `Save`, `Get`, `Update`, and `Delete`.
- Completed the first table-driven `Save` test, including retrieving the saved job to check storage.
- Started `Get` tests for existing, invalid, and missing IDs. These still need corrections.
- Discussed the service and handler responsibilities; neither layer nor the HTTP endpoints exists yet.
- The current implementation is in `main.go`, with tests in `job_repository_test.go`. Splitting into model/repository/service/handler packages was discussed but has not been done.

## Original 12-week roadmap and progress

| Week | Goal from the original plan | Current status |
| --- | --- | --- |
| 3 | Go foundation: `POST /jobs`, `GET /jobs/{id}`, `DELETE /jobs/{id}`; Handler → Service → Repository → in-memory storage; structs, interfaces, context, HTTP, maps, errors, and tests. | In progress: model and repository implemented; basic testing underway; service, HTTP, and context work remain. |
| 4 | Bounded worker pool with goroutines/channels, cancellation, graceful shutdown, mutexes, WaitGroups, and race-condition testing. | Planned; queue and worker concepts discussed, implementation not started. |
| 5 | Priority and `scheduled_at` fields; heap-based priority queue; insertion, removal, ordering, edge cases, and operation complexity. | Planned. |
| 6 | Maximum attempts, backoff, jitter, maximum delay, and idempotency keys; transient/permanent failures and exhausted retries. | Planned. |
| 7 | Workflow graph with job dependencies, DFS/cycle detection, and BFS/topological ordering; branching, merging, disconnected nodes, and cycles. | Planned; workflow/DAG purpose introduced conceptually. |
| 8 | Worker ownership, leases, heartbeats, expiration, and reassignment; reason about duplicate execution and idempotency. | Planned. |
| 9 | Sliding-window rate limiter; bursts, sustained traffic, timestamp boundaries, and concurrent callers. | Planned. |
| 10 | MySQL or PostgreSQL persistence, schema/migrations, indexes, transactions where appropriate, and pagination; restart and concurrency behavior. | Planned. |
| 11 | `classify_meeting` worker behind a provider interface; persist prediction, confidence, reasoning, model/version, and prompt version; handle timeouts and malformed output. | Planned. `MeetingInput` exists only as a data model; no classifier executes yet. Start with a mock provider as discussed. |
| 12 | Structured logs, metrics, simple dashboard, README, architecture/failure-mode documentation, and a 10× traffic thought experiment. | Planned; the current README contains only the repository title. |

The core plan is intended to run locally. Cloud deployment and a paid LLM provider are optional extensions, not prerequisites for progress.

## What I have built and learned

### 1. Understanding jobs, queues, workers, and workflows

I began by explaining the system in my own words before writing code:

- A **job** is a unit of work the system accepts, tracks, and eventually executes.
- A **queue** separates accepting work from executing it. With 100 jobs and 3 busy workers, the remaining jobs wait for capacity.
- A **worker** takes a job, executes it, and updates its status.
- `queued → running → completed` describes a job lifecycle; failure and retry states come later.
- A **workflow** relates jobs through dependencies. A queue alone does not enforce those dependencies; the graph/scheduler will handle them later.

These are concepts discussed so far. The current repository stores jobs but does not execute them or enforce status transitions.

### 2. Designing the job model

I moved from an email-specific nested input structure to a generic `Job` with separate payload types:

```go
type Job struct {
    ID      int64
    Status  string
    JobType string
    Input   any
}
```

`EmailInput` holds an email and message. `MeetingInput` holds a meeting ID and summary.

Lessons from that design:

- Metadata (`ID`, `Status`, `JobType`) describes the job; input contains the data needed to do the work.
- A job has one payload whose shape varies by job type. This is different from a one-to-many relationship.
- Capitalized fields are exported and accessible from other packages.
- `struct{}` has no fields; it is not a container that can acquire arbitrary fields later.
- I chose `any` for V1 because it was easier to understand than `json.RawMessage`. The tradeoff is that callers cannot assume a particular concrete input type without checking it.
- Representing input and dispatching execution are separate design problems. Worker dispatch and type assertions/switches remain future implementation work.

### 3. Choosing storage and separating responsibilities

I chose a map keyed by job ID because jobs need to be retrieved by ID:

```go
jobs map[int64]Job
```

This connects directly to DSA: hash-map lookup is average O(1), compared with O(n) for a linear search through a slice. The repository uses O(n) storage for n jobs.

I also worked through the intended architecture:

```text
HTTP → Handler → Service → Repository → Storage
```

- **Model:** the data passed between layers.
- **Handler/controller:** request parsing and HTTP responses.
- **Service:** business decisions, such as a new job starting as queued.
- **Repository:** storing, retrieving, replacing, and deleting records.

The repository does not decide what business fields should change. The service will supply the updated job. Keeping storage behind an interface is preparation for replacing the map with a database later.

### 4. Implementing the repository contract

I implemented this interface and matching methods:

```go
type JobRepository interface {
    Save(job *Job) (Job, error)
    Update(job *Job) error
    Delete(jobID int64) error
    Get(jobID int64) (Job, error)
}
```

Important decisions and lessons:

- Use `int64` consistently for IDs, map keys, method arguments, and expected test IDs.
- Interface and implementation method signatures must match exactly.
- `NewJobRepository` initializes the map with `make` and starts `nextID` at 1.
- `Save` assigns the next ID, stores the job value, increments the counter, and returns the saved job plus an error.
- Returning the saved job gives callers explicit access to the assigned ID and stored fields.
- `Save` also mutates the supplied job pointer by assigning its ID. Test setup must account for that.
- Check a pointer for `nil` before accessing its fields. A nil `*Job` is different from `&Job{Input: nil}`.
- Use the map's comma-ok lookup to distinguish a missing entry from a zero-value job.
- `Update` replaces an existing record; `Delete` verifies existence before removal.
- Errors communicate invalid input and missing records; returning a success error value (`nil`) for a failed operation hides the failure.
- Job IDs should not be recycled after deletion. The current counter preserves this within one repository instance; it resets for a new instance and is not a durable or distributed ID generator.
- Storing `*job` in `map[int64]Job` copies the struct value. This is not a guaranteed deep copy of arbitrary data held by `Input`.

### 5. Learning to test behavior

In `TestSaveTableDriven`, I added email, meeting, and nil-job cases. Each case gets a fresh repository, so the first successful save in each case expects ID 1.

The successful cases now check that `Save` succeeds, assigns the expected ID, and that `Get(savedJob.ID)` retrieves an equivalent job. This proves more than checking the return value of `Save` alone.

Lessons practiced:

- Go discovers test functions named `Test...` in `_test.go` files.
- A test table should model its inputs and expectations explicitly, such as `job`, `wantID`, and `wantErr`, instead of putting unrelated expectations into one `any` field.
- `t.Run` names each case; a `return` inside its callback exits that case.
- `t.Errorf` records failure and continues; `t.Fatalf` stops the current test path when later checks cannot be meaningful.
- An expected-error case must fail if the error is missing. Checking only `err != nil && ...` does not prove an expected error occurred.
- Verify the expected error and return before checking successful-result fields.
- `reflect.DeepEqual` can compare whole job values; comparing a `Job` value with a `*Job` pointer does not compare matching shapes.
- Test setup must preserve the condition being tested. Saving a supposedly missing job first destroys the missing-ID scenario.

## Open work at this checkpoint

The local `TestGetTableDriven` draft still has issues identified in the mentoring conversation and visible in the current code:

1. It saves every positive-ID fixture, including the supposedly missing job. `Save` changes that fixture's ID, so the test stops asking for the original missing ID.
2. It shares one repository across cases instead of isolating their setup.
3. It can miss an expected error when `err == nil`.
4. It compares a returned `Job` with a `*Job`, and also compares job contents after expected errors.
5. It ignores errors from setup calls to `Save`.

The correction is to separate the requested `jobID`, optional seed job, expected job value, and expected error. Keep a few clear cases: existing job, zero ID, negative ID, and absent positive ID.

Other known limitations: `Update` and `Delete` have no dedicated tests yet; error wording is inconsistent; `Update` rejects zero as well as negative IDs despite its message mentioning only negative IDs. The map and `nextID` are unsynchronized, which will need attention before concurrent workers use them. All state is currently in memory.

Verification on September 12, 2026: ran `go test ./...` with a temporary Go build cache. `TestSaveTableDriven` passed; all five `TestGetTableDriven` cases failed at the result comparison. The suite is currently failing, consistent with the unfinished `Get` tests above. No implementation or test code was changed while creating this log.

## Learning approach going forward

The September 12 mentoring discussion adjusted the emphasis: spend most effort on core logic, implementation, and tradeoffs. Use a few meaningful tests for ordinary CRUD behavior, then move on. Spend deeper reasoning and edge-case testing on the heap, workflow graph, scheduler, worker queue, and rate limiter.

I will continue writing the code first and using the mentoring conversation for questions, feedback, and explanations. Pair the project with a small number of targeted interview problems matching the current data structure; building this system does not replace timed coding practice.

## Next steps

1. Correct the small set of `Get` cases and rerun the tests; avoid spending more time perfecting the table structure.
2. Build the service layer and decide its job-creation and validation rules.
3. Add `POST /jobs`, `GET /jobs/{id}`, and `DELETE /jobs/{id}` through handlers, learning JSON and HTTP concerns as they arise.
4. Move into the bounded worker queue/pool milestone and its concurrency concepts.

## Evidence and continuity

This checkpoint was reconstructed from the Chat conversations **Explain Build Plan** (original roadmap and foundation discussions) and **Branch · Explain Build Plan** (September 12 testing/progress discussion), plus the current repository files and Git history. Planned features and mentor suggestions are not counted as completed implementation.

Repository milestones:

- September 3, 2026 — `48e2731`: first commit.
- September 11, 2026 — `3af02db`: repository and test file added.
- September 11, 2026 — `bd45c72`: test cases improved.
- September 11, 2026 — `d70aa44`: retrieval checks added to verify successful storage.
- At this checkpoint, `main.go` and the new `Get` test have local, uncommitted changes included in this assessment.

Update this file after meaningful milestones. Use it to resume mentoring or start a focused conversation for the next phase.

### Future entry template

- **Date / milestone:**
- **What I implemented:**
- **What I can now explain:**
- **Bug or misconception I corrected:**
- **Design decision and tradeoff:**
- **DSA connection / complexity:**
- **Verification performed and result:**
- **Next step:**
