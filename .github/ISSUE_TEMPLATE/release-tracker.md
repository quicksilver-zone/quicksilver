---
name: Release tracker
about: Create an issue to track release progress

---

<!-- < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < < ☺ 
v                            ✰  Thanks for opening an issue! ✰    
v    Before smashing the submit button please review the template.
v    Word of caution: poorly thought-out proposals may be rejected 
v                     without deliberation 
☺ > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > > >  -->

## Milestones

<!-- Links to alpha, beta, RC or final milestones -->

## QA

### Testing

<!-- List of tests that need be performed with previous
versions of quicksilver to guarantee that no regression is introduced -->

- [ ] [Simulation tests](https://github.com/ingenuity-build/quicksilver/tree/main/test/simulation) pass for the release branch.
- [ ] [Upgrade tests](https://github.com/ingenuity-build/quicksilver/tree/main/app) pass for the release branch.
- [ ] [E2E tests](https://github.com/ingenuity-build/quicksilver/tree/main/test/e2e) pass for the release branch.
- [ ] Testnet deployment testing.

## Migration

<!-- Link to migration document -->

## Checklist

<!-- Remove any items that are not applicable. -->

- [ ] Bump [go package version](https://github.com/ingenuity-build/quicksilver/blob/main/go.mod#L3). (Major release only)
- [ ] Change all imports starting with `github.com/ingenuity-build/quicksilver/v{x}` to `github.com/ingenuity-build/quicksilver/v{x+1}`. (Major release only)
- [ ] Branch off main to create release branch in the form  of `release/vx.y.z` and add branch protection rules.
- [ ] Add branch protection rules to new release branch.
- [ ] Update [`CHANGELOG.md`](https://github.com/ingenuity-build/quicksilver/blob/main/CHANGELOG.md)
- [ ] Add any necessary [retract](https://go.dev/ref/mod#go-mod-file-retract) statements to `go.mod`.
- [ ] Create new binary, tag and release.
- [ ] Build and push corresponding docker image.

## Post-release checklist

- [ ] Update [`CHANGELOG.md`](https://github.com/ingenuity-build/quicksilver/blob/main/CHANGELOG.md)
- [ ] Update [version matrix](https://github.com/ingenuity-build/quicksilver/blob/main/RELEASES.md#version-matrix) in `RELEASES.md`:
    - Add the new release.
    - Remove any tags that might not be recommended anymore.
- [ ] Update docs site with versioned docs:
____

#### For Admin Use

- [ ] Not duplicate issue
- [ ] Appropriate labels applied
- [ ] Appropriate contributors tagged/assigned