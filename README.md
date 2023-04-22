Jessica (Jess) is a program that helps developers to develop code and acts as a co-developer(s). It is designed to assist in improving code, initiating dialogs, and more. Jess is built on GPT-4, thus having GPT-4 is a requirement for its proper functioning.

## Features

* Explain/improve/refactor the code
* Usual dialogs with Jess (GPT4)

## Usage
* Goto https://platform.openai.com/account/api-keys and create a new secret key.
* create a file `.open-ai.key` under $HOME dir and copy the value of secret key generated in previous step.

### Generic Dialog

To interact with Jess, use the following commands and flags:

`jess dialog`:

* Manage dialogs
* Flags:
  * `--list` or `-l`: list all dialogs
  * `--continue ID` or `-c ID`: continue a dialog with the specified ID
  * `--show ID` or `-s ID`: show a dialog with the specified ID
  * `--delete ID` or `-d ID`: delete a dialog with the specified ID

### Coding Requests

`jessica file`:

* Read and process files
* Flags:
  * `--input PATH` or `-i PATH`: path to the input file to show to Jessica (required)
  * `--prompt TEXT` or `-p TEXT`: prompt input to be passed along with the file
  * `--refactor` or `-r`: suggest refactoring of the file by applying best practices

You need to provide either prompt or refactor flag. If refactor flag set, prompt will be ignored.

## Install

```bash
curl -sSL https://raw.githubusercontent.com/assistant-ai/jess/master/install.sh | bash
```
