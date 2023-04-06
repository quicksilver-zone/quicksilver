---
name: New Release
about: Prepare for a new Quicksilver release.
---

## 1. Summary
Release #
<!-- What are is the release we are reviewing? Major or minor release? -->

## 2. Checklist

- [ ] Upgrade handler added (must be tested if not no-op).
- [ ] Changelog up to date.

## 3. Post-merge Work 

- [ ] Build and tag binary from: 
  - `develop` if testnet release
  - `main` if mainnet release
- [ ] Build, tag and push docker image from: 
  - `develop` if testnet release
  - `main` if mainnet release

## 4. Additional Details (optional)

<!-- Add any additional info here. -->

