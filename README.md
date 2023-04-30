# Jessica CLI

Jessica is an AI assistant for software developers to help them with their code by explaining, refactoring, and answering questions. It is primarily used as a Command Line Interface (CLI) tool.

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

   List all dialogs:
   ```bash
   jess dialog -l
   ```

   Continue dialog with context ID:
   ```bash
   jess dialog -c <context_id>
   ```

   Show dialog with context ID:
   ```bash
   jess dialog -s <context_id>
   ```

   Delete dialog with context ID:
   ```bash
   jess dialog -d <context_id>
   ```

2. **Code processing**

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
