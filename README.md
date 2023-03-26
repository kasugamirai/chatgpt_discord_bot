# chatgpt_discord_bot
Discord Bot with GPT-3.5-turbo Integration

A Discord bot that uses the OpenAI GPT-3.5-turbo API to answer user queries and provide assistance in real-time.

1.Clone the repository:
git clone https://github.com/kasugamirai/chatgpt_discord_bot.git

2.Install required dependencies:
go get github.com/bwmarrin/discordgo

3.Set up the environment variables for your API keys:
export OPENAI_API_KEY=<your_openai_api_key>
export DISCORD_BOT_TOKEN=<your_discord_bot_token>

4.Build and run the bot:
go build
./discord-gpt-bot

Usage
Invite the bot to your Discord server using the Discord Developer Portal.

Use the bot in any text channel by typing commands with the appropriate prefix.

Commands
.<query>: Sends a query to GPT-3.5-turbo and returns the response.
Examples
.What is the capital of France?: The bot will reply with "The capital of France is Paris."
