*This is a WIP project*

*This README.MD file is paritally written by Gemini*

# GZCLP Tracker

A RESTful backend service written in Go for tracking and progressing workouts following the **GZCLP** (GZCL Linear Progression) strength training routine. The application manages user workout routines, handles tier-based progression/failure logic, advances rotating workout sessions, and persists state in Google Cloud Firestore.

---

## Overview

The GZCLP routine is a 4-day linear progression program structured around 3 exercise tiers per workout session. This backend provides:

- **Current Workout Retrieval**: Fetches the user's active workout session (e.g. `A1`, `B1`, `A2`, or `B2`) with the exact set, repetition, and weight scheme derived from the user's current progression stage.
- **Result Processing & Progression**: Ingests workout results (success/failure per exercise tier), automatically updates weights and progression stages according to GZCLP rules, rotates the user's routine to the next workout session, and saves the state transactionally.

---

## Architecture & Layering

The codebase follows Clean/Hexagonal Architecture principles, separating transport, domain logic, and persistence into clear layers with interfaces:

```
                  ┌────────────────────────┐
                  │    HTTP / REST API     │  (internal/http/rest)
                  │ (chi router, handlers) │
                  └───────────┬────────────┘
                              │
                              ▼
                  ┌────────────────────────┐
                  │      User Service      │  (internal/user)
                  │   (orchestration)      │
                  └──────┬──────────┬──────┘
                         │          │
        ┌────────────────┘          └───────────────┐
        ▼                                           ▼
┌───────────────┐                          ┌─────────────────┐
│ Domain Models │ (Routines, Workout)      │ Firestore Repo  │ (internal/storage/firestore)
│ & Progression │ (internal/user)          │ & Entities      │
└───────────────┘                          └────────┬────────┘
                                                    │
                                                    ▼
                                           ┌─────────────────┐
                                           │ Google Cloud    │
                                           │ Firestore / DB  │
                                           └─────────────────┘
```

## Configuration & Environment Variables

| Variable | Description | Example |
| :--- | :--- | :--- |
| `SERVICE_ACCOUNT_PATH` | File path to Google Cloud / Firebase service account credentials JSON. | `/path/to/serviceAccountKey.json` |

---

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) (version 1.22+ recommended)
- A Google Cloud Project with Cloud Firestore enabled, or a local Firestore emulator.

### Running the Server

1. Set the service account credentials path:
   ```bash
   export SERVICE_ACCOUNT_PATH="/path/to/serviceAccount.json"
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

   The server will start on port `3000`:
   ```
   Server running on :3000
   ```

---

## Testing

Run the full test suite across all packages:

```bash
go test -v ./...
```
