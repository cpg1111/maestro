# How to Contribute

Maestro and Maestrod are [Apache 2.0 licensed](LICENSE) and accept contributions via
GitHub pull requests.  This document outlines some of the conventions on
development workflow, commit message formatting, contact points and other
resources to make it easier to get your contribution accepted.

# Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. See the [DCO](DCO) file for details.

# Email

Currently you can reach the maintainer of Maestro and Maestrod at:

- Email: christiang11754@gmail.com

## Getting Started

- Fork the repository on GitHub
- Read the [README](README.md) for build instructions
- Play with the project, submit bugs, submit patches, submit feature requests!

## Feature Requests

When submitting a feature request please create a pull request with the following:

- An explanation of why this project needs this feature
- A proposal of how to implement it

## Contribution Flow

This is a rough outline of what a contributor's workflow looks like:

- Create a topic branch from where you want to base your work (usually master).
- Make commits of logical units.
- Make sure your commit messages are in the proper format (see below).
- Push your changes to a topic branch in your fork of the repository.
- Make sure the tests pass, and add any new tests as appropriate.
- Submit a pull request to the original repository.

Thanks for your contributions!

### Format of the Commit Message

We follow a rough convention for commit messages that is designed to answer two
questions: what changed and why. The subject line should feature the what and
the body of the commit should describe the why.

```
<subsystem>: <what changed>
<BLANK LINE>
<why this change was made>
<BLANK LINE>
<footer>
```

The first line is the subject and should be no longer than 70 characters, the
second line is always blank, and other lines should be wrapped at 80 characters.
This allows the message to be easier to read on GitHub as well as in various
git tools.

## Bug Reports


For reporting bugs, please create an issue with the following info:

- What version of maestro you are using

- Your maestro config sans any sensitive info

- any possible stacktraces either in build or runtime
