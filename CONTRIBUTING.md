# Contribution Guide

Welcome! We gladly accept contributions from the community. If you wish to
contribute code and you have not signed our [Contributor License Agreement
(CLA)](https://cla.vmware.com), our bot will update the issue when you open a
Pull Request. For any questions about the CLA process, please refer to our
[FAQ](https://cla.vmware.com/faq).

Once again, we hope you have read our [code of conduct](CODE_OF_CONDUCT.md)
before starting.

## Community

Wavefront by VMware Adapter for Istio contributors can be found on the
[VMware {code} Slack](https://vmwarecode.slack.com). To join the community,
please visit the VMware {code} [website](https://code.vmware.com/join).

## Logging Bugs

Anyone can log a bug using the GitHub 'New Issue' button. Please use a short
title and give as much information as you can about what the problem is,
relevant software versions, and how to reproduce it. If you know the fix or a
workaround include that too.

## Code Contribution Flow

We use GitHub Pull Requests to incorporate code changes from external
contributors. Typical contribution flow steps are:

* Sign the [CLA](https://cla.vmware.com)
* Fork the wavefront-adapter-for-istio repository
* Clone the forked repository locally and configure the upstream repository
* Open an Issue describing what you propose to do (unless the change is so
  trivial that an issue is not needed)
* Wait for discussion and possible direction hints in the issue thread
* Once you know which steps to take in your intended contribution, make changes
  in a topic branch and commit (don't forget to add or modify tests too)
* Fetch changes from upstream, rebase with master and resolve any merge
  conflicts so that your topic branch is up-to-date
* Build and test the project locally (see [BUILD.md](BUILD.md) for more details)
* Push all commits to the topic branch in your forked repository
* Submit a Pull Request to merge topic branch commits to upstream master

If this process sounds unfamiliar have a look at the excellent overview of
[collaboration via Pull Requests on GitHub](https://help.github.com/categories/collaborating-with-issues-and-pull-requests/)
for more information.

## Commit Messages

Commit messages should be self-sufficient to describe the problem the PR
addresses i.e. ideally without any references to private communication channels
like Slack, email, etc.

Below is a commit message template for reference.

```
Short (50 chars or less) summary of changes

More detailed explanatory text, if necessary.  Wrap it to about 72
characters or so.  In some contexts, the first line is treated as the
subject of an email and the rest of the text as the body.  The blank
line separating the summary from the body is critical (unless you omit
the body entirely); tools like rebase can get confused if you run the
two together.

Further paragraphs come after blank lines.

  - Bullet points are okay, too

  - Typically a hyphen or asterisk is used for the bullet, preceded by a
    single space, with blank lines in between, but conventions vary here

Signed-off-by: Some Developer <some.developer@example.com>
```
