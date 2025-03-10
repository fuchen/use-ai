# use-ai

`use-ai` is a command-line tool that uses Large Language Models (LLMs) to help users generate shell commands suitable for their operating system and shell environment.

## Features

- Automatically detects the user's operating system and shell environment
- Generates appropriate shell commands based on user queries
- Supports continuous dialogue, maintaining context
- Customizable OpenAI API endpoint, model, and system prompt

## Installation

### Building from Source

Ensure you have Go installed (Go 1.16+ recommended), then run:

```bash
git clone https://github.com/fuchen/use-ai
cd use-ai
go build
```

## Configuration

When first run, the program will create a `.use-ai.json` configuration file in your home directory. You need to edit this file to add your OpenAI API key:

```json
{
  "openai": {
    "endpoint": "https://api.openai.com/v1",
    "model": "gpt-4-turbo",
    "api_key": "YOUR_API_KEY_HERE",
    "system_prompt": "You are a professional command line assistant. Based on the user's query, only return the appropriate shell command without explanation. Make sure the command is suitable for the user's operating system and shell environment."
  }
}
```

## Usage

```bash
./use-ai
```

After starting the program, you can directly enter a description of the task you want to perform, and the AI will try to generate the appropriate shell command.

For example:
```bash
> How do I find all PDF files in the current directory?
```

The AI will return an appropriate command based on your operating system and shell environment.

Enter an empty line (just press Enter) to end the conversation and exit the program.

## License

MIT

## Notes

- Using this tool requires an OpenAI API key
- Generated commands may need manual verification, especially for potentially dangerous operations