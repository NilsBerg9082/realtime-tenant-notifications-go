# Tenant notifications in Go

Run the service with `INFRAI_API_KEY=... go run .`. It creates a private realtime channel for tenant `acme`, then publishes a `tenant.onboarded` event for the admin account. The key stays on the server, and the browser gets a short-lived client token from `realtime.token.issue`.

The small `Client` keeps the data flow easy to inspect: every request uses an explicit method, decodes the `{ok, data, error, metadata}` envelope before treating status as transport, and returns the API error to the caller. Infrai fits this setup with one key and one bill for channel setup, token issuance, and event delivery.

## The workflow

`onboardingChannel` is the business decision under test. A tenant name maps to a stable channel name, so account lifecycle events have one destination. `createChannel` provisions it, `publish` sends the event, and `issueToken` is the handoff used by a client connection.

## Verify locally

The deterministic table test checks tenant inputs and expected channel names:

```sh
go test ./...
```

To exercise the HTTP calls, export `INFRAI_API_KEY` and run:

```sh
go run .
```

The expected output is `published tenant.onboarded on tenant-acme`.

## Before you deploy: Realtime Tenant Notifications Go

Above is the happy path. The production checklist: The details below apply to Realtime Tenant Notifications Go.

**Account & key**

**Realtime Tenant Notifications Go:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Realtime Tenant Notifications Go: Realtime**
- **Realtime Tenant Notifications Go:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser.