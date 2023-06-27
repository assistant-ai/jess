# Jessica CLI

Jessica is an AI assistant for software developers to help them with their code by explaining, refactoring, and answering questions. It is primarily used as a Command Line Interface (CLI) tool.

* [Discord](https://discord.gg/7cUrMbcm)
* [Video With Example](https://youtu.be/j1ChHnMqmP4)

## Installation

To install, first, clone the repository and navigate to the project directory:

```bash
curl -sSL https://raw.githubusercontent.com/assistant-ai/jess/master/install.sh | bash
```

## Requirements

- Golang >= 1.15
- OpenAI API Key (stored in `~/.open-ai.key` file)

## Features

Jessica offers the following features:

1. **Manage Dialog**: Start, continue, list, show, and delete conversations with the AI assistant.
2. **Code Processing**: Perform various tasks such as:
   - Explain: Describe the code in plain English.
   - Refactor: Refactor the code following best practices.
   - Answer Questions: Answer questions about the code with possible code examples.

## Usage

1. **Dialog management**

If you want to just chat with Jess you should use dialog command. Dialog command allows you to either start a new dialog or to continue existing one. To start a new dialog that is persistent just come up with the unique id and start it like this:

   ```bash
   jess dialog -c <context_id>
   ```

This either will start a new dialog with the context id, or will continue dialog with this context (if it already exist). You can start dialog without the context_id:

   ```bash
   jess dialog
   ```

In this case dialog will NOT be persistent and will dissapear right after the end. You can check all dialogs that you had in the past by using -l key:

   ```bash
   jess dialog -l
   ```
If you do not want to continue dialog, but want to see all the messages from it, use -s key with the context id like this:

   ```bash
   jess dialog -s <context_id>
   ```

Finally, if you no longer want to have the dialogs sored, you can delete it:

   ```bash
   jess dialog -d <context_id>
   ```

2. **Context Work**

Context is something that is attached to everything you do with Jess. If you start dialog with key -c for the first time it will also store the empty context. Context is a persistent message that Jess should know about the dialog. For example, if I am going to start a dialog about coding a Go project I might set in context that I am coding in Golang, I use MacOS, project is a REST API and I deploy it on GCP.

Show context message:

   ```bash
   jess context -s <context_id>
   ```

Show list of contexts:

   ```bash
   jess context -l
   ```

Delete context. IMPORTANT: deleteing context will delete a dialog attached to it with all messages in it:

   ```bash
   jess context -d <context_id>
   ```

3. **Code processing**

   Explain code files in English:
   ```bash
   jess code explain -i <input_file1> -i <input_file2> -o <output_file>
   ```

   Refactor code file:
   ```bash
   jess code refactor -i <input_file> -o <output_file>
   ```

   Ask questions about code files:
   ```bash
   jess code question -p "Your question" -i <input_file1> -i <input_file2> -o <output_file>
   ```

Note: Replace `<input_file>` with the actual file paths and `<context_id>` with an actual context ID.

## Contributing and Support

Feel free to open issues, submit pull requests, or contact us if you need help.
