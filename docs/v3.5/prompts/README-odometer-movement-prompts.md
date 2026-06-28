# v3.5.6 Odometer Movement Prompt Bundle

Use these prompts after the current v3.5.6 PR has been merged to `main` and your local repo is freshly pulled.

Suggested flow:

1. Start from fresh `main`.
2. Run `01-document-odometer-movement-goal.md` first.
3. Review and merge that docs-only PR if clean.
4. Start again from fresh `main`.
5. Run `02-implement-odometer-movement-model.md`.

The intent is to avoid Codex inventing behaviour. The docs prompt creates the measuring tape. The implementation prompt cuts to that measurement.

Do not ask Codex to implement v3.5.7, snap/settle, backlash, or carry-drag as part of this work.
