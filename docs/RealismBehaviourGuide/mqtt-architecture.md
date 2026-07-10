# MQTT Architecture Notes

**Status:** Design discussion  
**Implementation:** Not started  
**Applies to:** Future GoDriveLog architecture

This document captures the current design direction for introducing MQTT into GoDriveLog. It is architectural guidance, not an implementation contract.

It currently lives in `docs/RealismBehaviourGuide/` because there is no broader architecture-notes directory yet. MQTT is not a gauge realism option.

## Executive summary

MQTT is a good fit for decoupling telemetry producers from dashboards and other live consumers.

The recommended first slice is:

```text
real OBD daemon / fake OBD daemon -> MQTT -> dashboard
```

MQTT should not automatically become the source of truth for stored trip or sensor history. It is a message bus, not a database. Durable local persistence still needs an explicit design such as SQLite, JSONL, or an append-only spool.

Do not combine MQTT transport, local storage, remote sync, fake data, replay, dashboard migration, and cloud deployment into one feature. That way lies distributed-system goblinry.

## Goals

- Decouple OBD acquisition from dashboard rendering.
- Let real and fake telemetry sources use the same interface and schema.
- Support multiple dashboards or other live consumers.
- Allow dashboard development without OBD hardware connected.
- Prepare for replay, logging, and remote bridging without requiring them in the first slice.
- Keep transport concerns separate from dashboard rendering and persistence.

## Non-goals

The first MQTT slice should not:

- replace durable local storage;
- provide complete trip-history persistence;
- implement remote cloud sync;
- expose an internet-facing broker;
- introduce binary payload formats without evidence that JSON is insufficient;
- migrate every GoDriveLog data path at once.

## Architectural principle

GoDriveLog already separates gauge behaviour from presentation. The same principle should apply here:

> MQTT is transport, not truth, storage, or rendering behaviour.

Dashboard code should consume a telemetry abstraction and should not need to know whether values came from real OBD hardware, a fake producer, replay data, or another valid source.

## Proposed daemon model

### `godrivelog-obd`

Real OBD producer.

Responsibilities:

- connect to an OBD adapter, ELM327, or equivalent source;
- poll configured PIDs;
- normalize units;
- publish telemetry to MQTT;
- publish source health and status;
- optionally write a local spool later if capture reliability requires it.

### `godrivelog-fake-obd`

Fake telemetry producer.

Responsibilities:

- publish the same schema as the real OBD daemon;
- support deterministic test values;
- support waveform, random, and demo values;
- later support replay scenarios.

### `godrivelog-dashboard`

Dashboard consumer.

Responsibilities:

- subscribe to telemetry topics;
- retain the latest known value for each signal;
- render dashboard gauges;
- show stale, missing, disconnected, and simulated states clearly;
- avoid depending on whether the source is real or fake.

### `godrivelog-logger`

Optional persistence consumer.

Responsibilities:

- subscribe to telemetry or import local spool files;
- write SQLite, JSONL, or another durable format;
- manage trips and sessions;
- handle reconnects and duplicate samples;
- support export and audit.

### `godrivelog-bridge`

Possible later remote bridge.

Responsibilities:

- subscribe to the local broker;
- publish to a remote broker or server;
- handle TLS, authentication, retry, and backoff;
- keep remote-sync concerns out of the OBD producer.

## Initial data flow

Recommended first implementation:

```text
fake OBD publisher -> MQTT -> dashboard
real OBD publisher -> MQTT -> dashboard
```

This proves:

- topic naming;
- payload schema;
- publisher and subscriber abstractions;
- stale-data behaviour;
- fake-source compatibility;
- dashboard decoupling.

Persistence and remote sync should follow only after this path is stable.

## Suggested topic structure

Keep the hierarchy predictable:

```text
godrivelog/{vehicle_id}/telemetry/{signal}
godrivelog/{vehicle_id}/status/source
godrivelog/{vehicle_id}/status/health
godrivelog/{vehicle_id}/control/{command}
```

Examples:

```text
godrivelog/caddy/telemetry/rpm
godrivelog/caddy/telemetry/speed
godrivelog/caddy/telemetry/coolant_temp
godrivelog/caddy/status/obd
godrivelog/caddy/status/health
```

Avoid clever topic hierarchies until there is a demonstrated need for them.

## Suggested telemetry envelope

Do not publish raw numbers alone. Include enough metadata to judge freshness, source, units, and quality.

```json
{
  "ts": "2026-07-10T18:15:30.123+10:00",
  "vehicle_id": "caddy",
  "source": "obd",
  "signal": "rpm",
  "value": 1840,
  "unit": "rpm",
  "quality": "ok"
}
```

A possible Go model:

```go
type TelemetrySample struct {
    Timestamp time.Time `json:"ts"`
    VehicleID string    `json:"vehicle_id"`
    Source    string    `json:"source"`
    Signal    string    `json:"signal"`
    Value     float64   `json:"value"`
    Unit      string    `json:"unit"`
    Quality   string    `json:"quality"`
}
```

A batched sample may later be preferable for logging because split signal topics can arrive slightly out of sync. The first slice should stay simple unless current code or testing shows a clear need for batching.

## Suggested interfaces

```go
type TelemetryPublisher interface {
    Publish(ctx context.Context, sample TelemetrySample) error
}
```

```go
type TelemetrySubscriber interface {
    Subscribe(ctx context.Context, handler func(TelemetrySample) error) error
}
```

These should be refined against the existing GoDriveLog packages before implementation. MQTT details should sit behind these abstractions rather than spreading through dashboard code.

## Dashboard stale-data handling

The dashboard must not display old data indefinitely as though it remains current.

Each signal should track:

- last update time;
- source connection state;
- freshness threshold;
- quality state;
- whether the signal is unsupported or missing.

Possible display states:

- live;
- stale;
- missing;
- disconnected;
- simulated.

An old speed value is not useful telemetry. It is a small lie wearing boots.

## MQTT QoS

QoS must be selected deliberately.

| QoS | Meaning | Likely fit |
|---|---|---|
| `0` | at most once | live dashboards where occasional drops are acceptable |
| `1` | at least once | logger input where duplicates are handled explicitly |
| `2` | exactly once at the MQTT protocol level | probably unnecessary for the first design |

QoS does not mean the data is durably stored. It only describes broker delivery semantics.

## Local persistence options

### Option A: Logger subscribes to MQTT

```text
OBD daemon -> MQTT -> dashboard
                  -> logger -> SQLite/JSONL
```

Advantages:

- clean separation;
- real and fake sources are handled identically;
- dashboard and logger are peer consumers;
- other consumers are easy to add.

Risks:

- broker downtime stops logging;
- late logger startup loses earlier messages unless they were buffered elsewhere;
- QoS 1 can duplicate samples;
- broker delivery is not database persistence.

### Option B: OBD daemon writes locally and publishes MQTT

```text
OBD daemon -> SQLite/JSONL/spool
           -> MQTT -> dashboard
```

Advantages:

- strongest local reliability;
- local capture survives broker failure;
- storage occurs close to the source.

Risks:

- producer becomes more complex;
- fake producers need matching behaviour if their data should also be logged;
- capture and dashboard paths differ.

### Option C: Append-only spool plus MQTT live stream

```text
OBD daemon -> append-only spool
           -> MQTT live telemetry

logger/importer -> spool -> SQLite/export
```

Advantages:

- broker failure does not lose captured data;
- supports replay and debugging;
- provides a clear audit trail;
- allows SQLite import later.

Risks:

- needs rotation and session rules;
- duplicate import must be prevented;
- requires more engineering than the first MQTT slice.

This is probably the best long-term design if GoDriveLog becomes both a live dashboard and a serious trip-history tool.

## Remote bridge

Remote MQTT is a separate feature, not a checkbox.

A safe design needs:

- TLS;
- authentication;
- per-device credentials;
- topic access controls;
- no anonymous write access;
- reconnect and backoff behaviour;
- privacy decisions for vehicle and location data;
- a clear distinction between local and remote brokers.

Recommended shape:

```text
local OBD/fake publisher -> local MQTT -> dashboard
                                      -> local logger
                                      -> remote bridge -> remote broker/server
```

The OBD daemon should not also become the secure remote-sync service.

## Recommended implementation order

1. Define the telemetry event model.
2. Define the MQTT topic contract.
3. Add `TelemetryPublisher` and `TelemetrySubscriber` abstractions.
4. Implement an MQTT publisher.
5. Implement the fake OBD publisher.
6. Implement the dashboard subscriber.
7. Add stale, missing, disconnected, and simulated display states.
8. Implement the real OBD publisher.
9. Decide local logging or spool architecture using evidence from the first slices.
10. Consider remote bridging only after the local system is stable.

Possible release slicing:

| Slice | Scope |
|---|---|
| MQTT-1 | telemetry model and topic contract |
| MQTT-2 | fake publisher and dashboard subscriber |
| MQTT-3 | real OBD publisher |
| MQTT-4 | stale data and status handling |
| MQTT-5 | logger or local spool decision |
| MQTT-6 | remote bridge, if still required |

## Main risks

### Over-architecture

The largest risk is turning a small useful application into a distributed system before it needs to be one. MQTT should solve the specific problem of decoupling producers from consumers.

### Local data loss

If GoDriveLog is used for trip records, tax reconstruction, or historical analysis, local durability matters more than elegant message flow.

### Broker dependency

Once MQTT is central to live telemetry, the dashboard depends on broker availability. Installation, service startup, reconnect behaviour, and diagnostics therefore matter.

### Duplicates and ordering

QoS 1 may duplicate messages. Separate signal topics may arrive out of sync. Consumers and loggers must not assume perfect delivery order.

### Remote security

An internet-facing broker must not be added casually. Poorly secured MQTT is an invitation for the internet to come in and rearrange the cutlery.

## Design principles

1. MQTT is transport, not truth.
2. Storage remains explicit and local-first.
3. Real and fake sources publish the same schema.
4. Dashboards do not care whether the source is real or simulated.
5. Dashboards show stale and disconnected states.
6. Remote sync remains a separate bridge feature.
7. JSON is adequate for the first implementation.
8. Topic structure stays boring and predictable.
9. Interfaces are introduced before MQTT is hardcoded throughout the application.
10. Work is delivered in narrow, testable slices.

## Current recommendation

Start with:

```text
fake OBD publisher -> MQTT -> dashboard
```

Then add:

```text
real OBD publisher -> MQTT -> dashboard
```

Do not make MQTT the mandatory durable-storage path yet. The likely long-term shape is:

```text
OBD/fake source -> local event model -> MQTT live stream
                                -> optional append-only local spool
MQTT -> dashboard
spool/logger -> SQLite/export
optional bridge -> remote MQTT/server
```

MQTT can be GoDriveLog's live nervous system. It should not be mistaken for the memory, the database, or the tax records.
