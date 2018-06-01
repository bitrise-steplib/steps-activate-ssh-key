## Changelog (Current version: 4.0.0)

-----------------

### 4.0.0 (2018 Jun 01)

* [54ff734] Prepare for 4.0.0
* [dee7fed] Merge pull request #9 from bitrise-io/rewrite_in_go
* [a5cb039] PR clean - fix (+8 squashed commits) Squashed commits: [842bdc7] go package_name URL fix [eb832f3] PR clean - fix [49d4a11] PR clean - fix [d7fd90c] oops.. [26ebc52] PR clean - fix [29dc4eb] PR fix - clean [d7776ad] SSH_AUTH_SOCK output export [bf3db8f] eval command removed (+11 squashed commits) Squashed commits: [4b9e53c]  - [88a3320]  - [74152b7]  - [594dd48] eval ssh-agent [9c68bd1]  - [87864f2]  - [57f5b19]  - [e282d36]  - [8d021d6] step.sh removed [93c6886] bitrise.yml update [02c6fb9] - rewrite in go :white_check_mark: - ssh_rsa_private_key will be removed after the config parse :rotating_light:
* [8a8b88c] main.go added:  - Writhe SSH RSA Private Key to file :white_check_mark:  - Merging from step.sh to main.go still in progress :construction:
* [77960f0] Merge pull request #7 from bitrise-io/update-tags
* [ca461cd] Rename ssh-key env
* [9f69b81] Update tags
* [99ff8c4] Update the quite outdated README
* [f2f96c6] `- MY_STEPLIB_REPO_FORK_GIT_URL: $MY_STEPLIB_REPO_FORK_GIT_URL`

### 3.1.1 (2016 Aug 01)

* [ba11196] v3.1.1
* [6f6eafb] minimal update
* [060a351] removed project type tags
* [6dedc0f] SSH_AUTH_SOCK title typo

### 3.1.0 (2015 Nov 26)

* [ca5c5ae] logging
* [9fae52d] instead of adding the `SSH_AUTH_SOCK` env to `~/.bashrc` this step'll now expose it with `envman`

### 3.0.3 (2015 Oct 29)

* [5c42e6f] Merge pull request #5 from bazscsa/patch-2
* [74100c4] added apt-get dependencies: git and expect
* [6c4a574] Merge pull request #4 from bazscsa/patch-1
* [abd6650] Update step.yml

### 3.0.2 (2015 Sep 03)

* [d4d437d] Merge pull request #3 from gkiki90/new_yml_format
* [221fd98] yml fix

### 3.0.1 (2015 Sep 03)

* [c493749] Merge pull request #2 from gkiki90/new_yml_format
* [9a1c4b1] required field fix

### 3.0.0 (2015 Sep 03)

* [72bafe6] Merge pull request #1 from gkiki90/new_yml_format
* [5572549] removed source
* [4ae8214] fixes
* [716991f] bitrise
* [eccb034] yml fix, gitignore
* [d58c15e] new format
* [e343bda] utils 'fail_if_cmd_error' fix

### 2.0.0 (2014 Nov 13)

* [89ae2f5] even more debug log
* [9682d04] minor debug / info log
* [9804655] input name fix: ssh rss PRIVATE key instead of public (because this is a private key, not the public key); change in how the agent reload works - ignores the case when it can't kill the already existing and functional agent, in that case it just removes all the identities from it
* [e7a9582] minor debug log revision
* [24eeab8] option to remove previously added identities before activating the new one

-----------------

Updated: 2018 Jun 01