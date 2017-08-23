# So you wanna update the WRT metadata?

## The Check-list

Great! We really appreciate it and there are only a few small things
you should check before throwing the PR to RelEng:

- [ ] Are you submitting this PR against the correct release branch?
- [ ] Have you manually verified that your changes deploy/work?
- [ ] Would you expect WATs to be affected by this change? If so, have you run WATs against a deployment with this change?
- [ ] Does this change rely on a particular release version? If so, which one? ______
- [ ] Are you available for a cross-team pair to help troubleshoot this change?
- [ ] Are there corresponding changes to the ERT? If so, have you submitted a corresponding PR to [p-runtime](https://github.com/pivotal-cf/p-runtime)?
- [ ] Have you made separate PRs against the branches for each legacy version of the WRT that is affected?

### Useful Information for Consuming this Change

Please fill out this section with any information that we will find helpful
when consuming this change. This might include a link to a story in your
tracker backlog with more details on the change and acceptance criteria for our PM.

Feel free to join us the #pcf-releng channel if you run into any problems.
