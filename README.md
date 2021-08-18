# Activate SSH key (RSA private key)

[![Step changelog](https://shields.io/github/v/release/bitrise-io/steps-activate-ssh-key?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-io/steps-activate-ssh-key/releases)

Setup the SSH Key to use with the current workflow

<details>
<summary>Description</summary>

This Step makes sure Bitrise has access to your repository and thus able to clone your code to our virtual machines. The Step saves the provided private key of your SSH keypair to a file and then loads it into the user's ssh-agent with `ssh-add`.

### Configuring the Step

By default, you do not have to change anything about the Step's configuration. All you need to do is make sure that you registered your key pair on Bitrise and the public key at your Git provider. You can generate and register an SSH keypair in two ways.

- Automatically during the [app creation process](https://devcenter.bitrise.io/getting-started/adding-a-new-app/#setting-up-ssh-keys).
- Manually during the app creation process or at any other time. You [generate your own SSH keys](https://devcenter.bitrise.io/faq/how-to-generate-ssh-keypair/) and register them on Bitrise and at your Git provider. The SSH key should not have a passphrase! 

Optionally, you can save the private key on the virtual machine. If a key already exists on the path you specified in the **(Optional) path to save the private key** input, it will be overwritten.

### Troubleshooting

If the Step fails, check the public key registered to your Git repository and compare it to the public key registered on Bitrise. The most frequent issue is that someone deleted or revoked the key on your Git provider's website.

You can also set the **Enable verbose logging** input to `true`. This provides additional information in the log.

### Useful links

- [Setting up SSH keys](https://devcenter.bitrise.io/getting-started/adding-a-new-app/#setting-up-ssh-keys)
- [How can I generate an SSH key pair?](https://devcenter.bitrise.io/faq/how-to-generate-ssh-keypair/)

### Related Steps

- [Git Clone Repository](https://www.bitrise.io/integrations/steps/git-clone)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `ssh_rsa_private_key` |  | sensitive | `$SSH_RSA_PRIVATE_KEY` |
| `ssh_key_save_path` |  |  | `$HOME/.ssh/bitrise_step_activate_ssh_key` |
| `is_remove_other_identities` | (Optional) Remove other or previously loaded keys and restart ssh-agent?  Options:  * "true" * "false" |  | `true` |
| `verbose` | Enable verbose log option for better debug | required | `false` |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `SSH_AUTH_SOCK` | If the `is_should_start_new_agent` option is enabled, and no accessible ssh-agent is found, the step will start a new ssh-agent.  This output environment variable will contain the path of the SSH Auth Socket, which can be used to access the started ssh-agent. |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-io/steps-activate-ssh-key/pulls) and [issues](https://github.com/bitrise-io/steps-activate-ssh-key/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

**Note:** this step's end-to-end tests (defined in `e2e/bitrise.yml`) are working with secrets which are intentionally not stored in this repo. External contributors won't be able to run those tests. Don't worry, if you open a PR with your contribution, we will help with running tests and make sure that they pass.

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
