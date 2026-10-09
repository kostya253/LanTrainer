# LanTrainer

LanTrainer is a Telegram bot designed to help users learn new languages using the "Listen, Shadow, Recall" technique. It allows users to manage personal and global sentence dictionaries and quiz themselves on them.

## Features

- **Dictionary Management**: Add and manage sentences in your personal dictionary or explore a global one.
- **Interactive Learning**: Use the bot to translate and recall sentences (Listen, Shadow, Recall).
- **Customizable**: Configure translation prompts and voice settings via an `.env` file.

## Prerequisites

- **Hardware**: macOS (recommended: Mac mini with 24 GB RAM).
- **Software**:
  - [Ollama](https://ollama.com/) (with `translategemma:12b` model: `ollama pull translategemma:12b`).
  - [ffmpeg](https://ffmpeg.org/) (Install via Homebrew: `brew install ffmpeg`).
  - [Go](https://go.dev/) (Install via Homebrew: `brew install go`).
  - A configured Telegram Bot (follow [this guide](https://medium.com/@Elhazin/creating-a-telegram-bot-55e6ca4e337d)).
  - Your Telegram ID (can be found using `@userinfobot`, `@getmyid_bot`, or `@raw_data_bot`).

## Installation & Configuration

1. **Clone the repository**:
   ```bash
   git clone https://github.com/kostya253/LanTrainer.git
   cd LanTrainer
   ```

2. **Setup Environment Variables**:
   Create a `.env` file in the root directory with the following content:
   ```env
   BOT_TOKEN=your_bot_token
   BOT_NAME="@your_bot_name_bot"
   AUTHORIZED_USER=Your_telegram_id

   # LLM Model to use
   TRANSLATE_MODEL=translategemma:12b

   # Learn Czech
   SAY_VOICE="Zuzana (Premium)"
   SAY_SPEED=175
   TRANSLATE_PROMPT="You are a professional English (en) to Czech (cs) translator. Your goal is to accurately convey the meaning and nuances of the original English text while adhering to Czech grammar, vocabulary, and cultural sensitivities. Produce only the Czech translation, without any additional explanations or commentary. Please translate the following English text into Czech:"
   ```

## Compilation & Execution

### Compile
Build the binary using the Go toolchain:
```bash
go build -o LanTrainer
```

### Run
Start the bot server:
```bash
./LanTrainer
```

## Bot Commands

- `/start`: Initialize the bot.
- `/es`: Add a new sentence to your personal dictionary.
- `/ts:my`: Switch to your personal dictionary.
- `/ts:g`: Switch to the global dictionary.
