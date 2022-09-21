# MedeiroTech Bot

A bot created with golang to serve [Medeiro.Tech Discord Community](https://discord.gg/k6hFV5HxMv).

# Tecnologies

- Golang 1.19
- [Dotenv](github.com/joho/godotenv)
- [DiscordGo](github.com/bwmarrin/discordgo)
- [Docker](https://docker.com)

# How to Run

You have to [setup a golang environment](https://go.dev) in your local machine and after this, do the following steps:

1. Install the dependencies of the module

```
$ go mod download
```

2. Create a `.env` file based on the `.env.example`, making the appropriate replacement

3. Start the process

```
go run main.go
```