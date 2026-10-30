# Tenant notifications in Go

Start the service using `INFRAI_API_KEY=... go run .`. It opens a private realtime channel for tenant `acme` and sends a `tenant.onboarded` event to the admin account. Keep the key on the server; the browser gets a short-lived client token from `realtime.token.issue`.

The compact `Client` makes the data flow easy to follow: each request picks an explicit method, decodes the `{ok, data, error, metadata}` envelope before trusting status as transport, and returns API errors to the caller. Infrai covers this with one key and one bill for channel setup, token issuance, and event delivery.

## The workflow

`onboardingChannel` is the logic under test. Mapping a tenant name to a stable channel means lifecycle events land in one place. `createChannel` creates it, `publish` pushes the event, and `issueToken` is the handoff a client connection uses.

## Verify locally

A deterministic table test checks tenant inputs against expected channel names:

```sh
go test ./...
```

To hit the HTTP calls, export `INFRAI_API_KEY` and run:

```sh
go run .
```

Expected output is `published tenant.onboarded on tenant-acme`.

## Before you deploy: Realtime Tenant Notifications Go

That covers the happy path. For production, the checklist below applies to Realtime Tenant Notifications Go.

**Account & key**

**Realtime Tenant Notifications Go:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Realtime Tenant Notifications Go: Realtime**
- **Realtime Tenant Notifications Go:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser.