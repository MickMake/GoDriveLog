# GoDriveLog Pi4 Fyne Kiosk Setup

Design reference: [`docs/Designs/Runtime/setup-kiosk-mode-on-pi4.md`](../../Designs/Runtime/setup-kiosk-mode-on-pi4.md)

## Purpose
Tracks the Raspberry Pi kiosk setup note for a Fyne-based display stack.

## Implementation Status
Status: **Not implemented**.

The current repo does not implement the documented Fyne kiosk stack.

## Packages and Files
- [`cmd/GoDriveLog/main_ebiten.go`](../../../cmd/GoDriveLog/main_ebiten.go)

## Types
- None in current code.

## Functions and Methods
- `main` builds the current Ebiten-based app/CLI surface, not the documented Fyne display flow.

## Runtime Flow
There is no Pi4 kiosk bootstrap path, packaging flow, or service setup in the repo that matches this note.

## Configuration
The note describes deployment and OS setup rather than in-repo runtime config; none of that is wired into current code.

## Behaviour
GoDriveLog can be built and run in its current desktop/runtime modes, but not as the specific documented Pi4 Fyne kiosk stack.

## Rendering
Current dashboard rendering is Ebiten-based, which is a different stack from the Fyne setup described here.

## Tests
- [`cmd/GoDriveLog/main_ebiten_test.go`](../../../cmd/GoDriveLog/main_ebiten_test.go)

## Limitations
This document is closer to a deployment recipe than a code slice, and the named stack does not match current implementation.

## Deviations from Design
The codebase has moved around a different rendering/runtime path than the note describes.

## Remaining Work
Either rewrite the deployment doc to match the actual stack or implement a dedicated kiosk/deployment path intentionally.

## Verification Notes
Verified by reading the linked code and test files on 2026-07-12. This was a documentation audit only; no Go implementation changes were made as part of this pass.
