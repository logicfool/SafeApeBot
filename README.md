# SafeApeBot

Telegram group management and DeFi tooling bot built in Go.  
SafeApeBot focuses on **group safety**, **automation**, and **on-chain analysis** (honeypot checks, price lookups, and mempool watchers) for EVM networks.

> Status: Work in progress / early alpha.

---

## Table of Contents

- [SafeApeBot](#safeapebot)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Key Features](#key-features)
    - [Group Administration](#group-administration)
    - [DeFi \& On-Chain Tooling](#defi--on-chain-tooling)
    - [Bot Platform \& Multi-Bot Support](#bot-platform--multi-bot-support)
    - [Data Layer \& Caching](#data-layer--caching)
    - [Smart Contract (GoSniper)](#smart-contract-gosniper)
  - [Implemented vs Pending Features](#implemented-vs-pending-features)
    - [Implemented](#implemented)
    - [In Progress](#in-progress)
    - [Planned](#planned)
  - [Architecture](#architecture)
  - [Tech Stack](#tech-stack)
  - [Configuration](#configuration)
  - [Running the Bot](#running-the-bot)
  - [Project Status \& Roadmap](#project-status--roadmap)
  - [License](#license)

---

## Overview

SafeApeBot is a **Golang-based Telegram bot** that combines:

- **Group administration** (moderation, filters, blacklists)
- **DeFi utilities** (honeypot checks, token price queries, mempool watchers – partially implemented)
- **Multi-bot platform** capabilities (root bot + user clones, stored in MongoDB)
- **On-chain smart contract tooling** via a Solidity contract (`GoSniper`) for DEX interactions and honeypot detection.

The goal is to offer **production-ready group safety tooling** for crypto projects and communities while exposing a **clean, typed, and extensible Golang backend**.

---

## Key Features

### Group Administration

- **Admin checks**  
  - Validates that the caller is a group admin before executing sensitive actions (ban/kick/unban).
- **Member management**  
  - `/ban` – ban a user by replying to one of their messages.  
  - `/unban` – unban a previously banned user.  
  - `/kick` – temporary ban (kick) a user by reply.
- **Content filtering system**  
  - `/filter` – add keyword-based filters to a group.  
    - Supports *text responses* and *media filters* (photo, animation, video, audio, document, sticker, voice).  
    - Filters are persisted per-chat in MongoDB.
  - `/removefilter` – remove a specific filter by keyword.  
  - `/filters` – list all filters configured for the current chat.  
  - `/removeallfilters` – clear all filters for the current chat.
  - Runtime message hook automatically scans all non-command messages and fires matching filters.
- **Service message cleanup**  
  - Detects and auto-deletes “service” messages (joins, leaves, etc.) to keep groups clean.
- **Basic utilities**  
  - `/start` – welcome message with inline keyboard (Help / Add to group / Join SafeApe chat).  
  - `/help` – detailed help text and command list.  
  - `/hi` – simple connectivity test command.  
  - `/testphoto` – sends a test photo to verify media sending capabilities.

### DeFi & On-Chain Tooling

From the current codebase and constants:

- **Goal / design**  
  - Provide Telegram commands for:
    - Honeypot checks  
    - Token price queries  
    - Token trade notifications (buy/sell events)  
    - Wallet watchlists and gains tracking
- **Current status**  
  - Exposed command: `/scan` → `ScanContract` (currently a stub; real logic to be wired).  
  - Solidity contract `GoSniper` is included for DEX interactions and honeypot checks.
  - High-level DeFi commands are defined in the help text but not all implemented in Go yet (see [Pending Features](#implemented-vs-pending-features)).

### Bot Platform & Multi-Bot Support

- **Multi-bot architecture**  
  - `DBBotStruct` stores multiple bot instances, with:
    - `BotToken`  
    - `SudoUser` (owner of the bot)  
    - `Chats` (per-chat settings incl. filters/blacklist)  
    - `RootBot` and `Enabled` flags
- **Commands for bot registration**  
  - `/addbottoroot` – promotes the current bot to a root bot in the DB.  
  - `/addbotindb` – registers the current bot as a non-root “clone” with persistence.
- **Webhook entrypoint for multiple bots**  
  - HTTP endpoint `POST /bot/:token` receives Telegram webhook updates for arbitrary bot tokens.  
  - For each webhook call:
    - A `bothook.Bot` is created for the token.  
    - It internally delegates to the shared `bot.Bot` command and message logic.  
  - Enables running a **platform of bots** with shared logic and per-token configuration.

### Data Layer & Caching

- **MongoDB persistence** (see `db/` and `db_structure.txt`):
  - Users (`DBUserStruct`): user state, premium flags, sudo status, etc.
  - Bots (`DBBotStruct`): token, owner, enabled/root flags, per-chat settings.
  - Chats (`DBChatStruct`):  
    - `Filters`  
    - `Blacklist` (structure present; logic WIP)  
    - `KickSpamMode`, `Warnings`, `VerificationType`, `ETHAddressAutoRemove`, etc.
  - Mempool (`DBMempoolStruct`): structures for tracking mempool subscriptions per network/chat/bot.
- **MongoDB client management**  
  - Centralized `InitDb` with connection and ping checks and a typed accessor layer.
- **In-memory cache (planned/partial)**  
  - `utils/cachehandler.go` defines a global cache (`go-cache`) for:
    - User objects (`DBUserStruct`)  
    - Bot objects (`DBBotStruct`)
  - Wiring to the rest of the code is still in progress.

### Smart Contract (GoSniper)

- **Solidity contract**: `resources/contract.sol`
  - Compiler: `pragma solidity ^0.8.6;`
- **DEX tooling**  
  - Generic interfaces for:
    - `IERC20`  
    - DEX Router (`IRouter`)  
    - DEX Factory (`IFactory`)  
    - Liquidity Pair (`ILiquidityPair`)  
  - `SwapHelper` library for:
    - Token sorting and reserve fetching  
    - Manual `getAmountOut` computations  
- **Contract features (GoSniper)**  
  - Multi-chain WETH handling (Ethereum, BSC, Polygon, Avalanche) via `setChainData`.  
  - Gas optimization using `ChiGasToken` with `discountCHI` modifier.  
  - `RouterSwapBySplit` / `RouterSwapBySplitWithGasRefund`  
    - Splits a single ETH order into multiple swaps for better execution.  
  - `TradeDirectlyByPair` / `TradeDirectlyByPairWithGasRefund`  
    - Direct interaction with a DEX pair with custom path and split logic.  
  - `HoneyPotCheck`  
    - Simulates a buy and sell to detect honeypot behavior by comparing:
      - Expected tokens vs. actual tokens received  
      - Expected ETH vs. returned ETH after selling
  - Utility function `checkContractAllowance` to verify ERC-20 allowances.

Integration into the Go layer is partially planned but not fully wired yet.

---

## Implemented vs Pending Features

### Implemented

- **Core bot**
  - `/start`, `/help`, `/hi`
  - Inline keyboards for help, add-to-group, and community chat links
  - Dual mode: long polling + HTTP webhook (`/bot/:token`)

- **Group administration**
  - `/ban`, `/unban`, `/kick`
  - Admin-only enforcement via `IsChatAdmin`
  - Automatic service message cleanup
  - Rich filter system:
    - `/filter` – add filters (text + media types)  
    - `/removefilter`  
    - `/filters`  
    - `/removeallfilters`  
    - Runtime `SearchAdnSendFilter` hook that responds when filter keywords are seen

- **Bot platform**
  - Multi-bot data model (`DBBotStruct`, `DBChatStruct`)
  - `/addbottoroot` and `/addbotindb` commands to register bots in the DB
  - Webhook handler forwarding to shared bot/session logic

- **Persistence & infra**
  - MongoDB initialization and typed CRUD operations for:
    - Users  
    - Bots  
    - Mempool entries (add/update/find by user/chat)
  - Basic utilities: Ethereum address validation, deep copy helper, reflection helper

- **Smart contract**
  - Full Solidity contract (`GoSniper`) implementing DEX trade utilities and honeypot checks

### In Progress

These have **visible code paths / structures**, but are not fully implemented or wired:

- **Blacklist system**
  - `/addblacklist` exists with parsing logic and TODO markers:
    - Needs completion to store and enforce blacklist separately from generic filters.
  - `DBChatStruct.Blacklist` present but underused.

- **Spam / verification / warnings**
  - `DBChatStruct` includes:
    - `KickSpamMode`  
    - `Warnings`  
    - `VerificationType`  
    - `ETHAddressAutoRemove`
  - `/spamban` handler stub exists but contains no enforcement logic yet.
  - No implemented commands yet for:
    - Warnings count per user  
    - Post-join / pre-join verification flows

- **Mempool watcher plumbing**
  - MongoDB structs and methods (`DBMempoolStruct`, `FindMempoolByUser`, `FindMempoolByChat`, `AddMempool`, `UpdateMempool`) are implemented.  
  - User-facing mempool commands (add/remove tokens, showsells, etc.) are not yet implemented.

- **DeFi bot commands**
  - `/scan` command is registered, but `ScanContract` is a stub.  
  - Go code still needs to integrate with:
    - The `GoSniper` contract  
    - Network configuration from `config.json`
  - DeFi commands mentioned in help but missing implementations (see “Planned” below).

- **Caching layer**
  - `utils.Cache` is defined (using `go-cache`), with helpers for `GetUser`, `SetUser`, `GetBot`, `SetBot`.  
  - Integration into DB accessors is not yet fully wired.

### Planned

From `TODO`, help text, and data structures:

- **Group commands**
  - `/kickme` – allow users to self-kick.  
  - `/permaban` – permanent ban shortcut command.  
  - `/mute` / `/muteall` – mute users or whole groups.  
  - `/warn` – track warnings per user (`UserWarns` in `db_structure.txt`).  
  - `/pin`, `/unpin` – pin/unpin messages.  
  - `/setwelcome`, `/welcome on|off` – welcome system configuration.  
  - `/joinverify on|off`, `/preverify on|off` – join/pre-join verification flows.  
  - `/postverify on|off` – post-action verification toggle.  
  - `/removeblacklist` – management for blacklist entries.

- **Mempool / token tracking**
  - `/addtoken` – subscribe a group/channel to token buys/sells.  
  - `/removetoken` – remove subscription.  
  - `/modifymessage` – customize token notification templates.  
  - `/showsells` – toggle sell notifications.  
  - `/check` – inspect which tokens are tracked in a group.

- **DeFi commands**
  - `/hp` – honeypot check, likely leveraging `GoSniper.HoneyPotCheck`.  
  - `/price` – token price lookup.  
  - `/networks` – list supported networks from `config.Networks`.  
  - `/tokens <address>` – list token balances and values for a wallet.  
  - `/gains <address>` – compute realized gains across tokens for a wallet.

- **Wallet watchlist**
  - `/setwatchlist <address>` – track a wallet’s on-chain activity.  
  - `/removewatchlist <address>` – remove from watchlist.

- **Clone / self-service bots**
  - `/clone` (mentioned in help text for root bot) – allow users to create their own bot instance reusing the same infrastructure.

- **Additional DeFi features (from TODO)**
  - Enumerate and post **all pairs for a given chain + exchange** into a channel using an external “pair saver” mechanism.

- **Extended DB schema**
  - Collections listed in `db_structure.txt` such as:
    - `Bot_channels`, `Networks`, `TokensLinked`, `BotAdded`, `ChannelSettings`, `UserWarns`, `ChannelFilters`, etc.,
  - Planned for future expansion of channel-level settings and network configuration.

---

## Architecture

- **Entry point**: `main.go`
  - Reads `config.json` → `structs.SConfig`
  - Initializes MongoDB (`db.InitDb`)
  - Creates and configures `bot.Bot`
  - Starts:
    - Long-polling bot (`bot.Start()`)  
    - HTTP server (`server.Server`) for webhook-based updates

- **Bot layer (`bot/`)**
  - `Bot` struct wraps:
    - `*tgbotapi.BotAPI`  
    - Bot metadata (`Me`)  
    - Config (`*SConfig`)  
    - DB handle (`*db.DB`)  
    - `CommandHandler`
  - Command routing via `CommandHandler` (`map[string]func(tgbotapi.Update)`)  
  - Message routing:
    - Commands: `HandleCommand`  
    - Normal messages: `SearchAdnSendFilter` + potential future handlers
    - Callback queries: `ActionHandler` and `Helphandler`

- **HTTP server (`server/`)**
  - `github.com/julienschmidt/httprouter` router.  
  - Route: `POST /bot/:token` → `HandleBotWebHook`.
  - `server/bothook` package:
    - Creates a `bothook.Bot` for the given token and update.  
    - Bridges HTTP updates to the core `bot.Bot` logic (so polling and webhooks share code).

- **Data layer (`db/` + `structs/`)**
  - Typed Go structs for users, bots, chats, mempool, and config.
  - Centralized MongoDB client with per-collection operations.
  - Config-driven target DB name and MongoDB URI.

- **Smart contract (`resources/contract.sol`)**
  - EVM-side utility for trading and honeypot checking.
  - Designed to be called from an off-chain component (e.g., Go or other services).

---

## Tech Stack

- **Language**: Go 1.18  
- **Messaging**: Telegram Bot API (`github.com/go-telegram-bot-api/telegram-bot-api/v5`)  
- **Web framework**: `github.com/julienschmidt/httprouter`  
- **Database**: MongoDB via `go.mongodb.org/mongo-driver`  
- **Caching**: `github.com/patrickmn/go-cache`  
- **Debug utilities**: `github.com/davecgh/go-spew/spew`  
- **Smart contracts**: Solidity ^0.8.6 (EVM chains: Ethereum, BSC, Polygon, Avalanche, etc.)

---

## Configuration

Configuration is loaded from a JSON file (`config.json`) into `structs.SConfig`:

```json
{
  "privatekey": "YOUR_PRIVATE_KEY",
  "networks": [
    {
      "node_url": "https://rpc.your-node.example",
      "chain_id": 1,
      "name": "Ethereum",
      "explorer": "https://etherscan.io"
    }
  ],
  "BotToken": "YOUR_TELEGRAM_BOT_TOKEN",
  "BotGroup": "safeapechat", 
  "MongoDburl": "mongodb://localhost:27017",
  "DB": "safeapebot",
  "ServerPort": 8080
}
```

Fields (from `structs.SConfig`):

- `Privatekey` – EVM private key for on-chain operations.  
- `Networks` – array of supported networks (`NodeURL`, `ChainID`, `Name`, `Explorer`).  
- `BotToken` – Telegram bot token.  
- `BotGroup` – Telegram username for the main support group (used in `/start` and `/help` CTAs).  
- `MongoDBURI` (`MongoDburl`) – Mongo connection string.  
- `DatabaseToUse` (`DB`) – DB name.  
- `ServerPort` – HTTP server port.

---

## Running the Bot

1. **Prerequisites**
   - Go 1.18+  
   - MongoDB running and reachable at `MongoDburl`  
   - A Telegram bot token (via [@BotFather](https://t.me/BotFather))

2. **Clone and build**

```bash
git clone <repo-url>
cd SafeApeBot
go build
```

(Windows helper script: `run.ps1` runs `go build` and then the compiled `SafeApeBot.exe`.)

3. **Create `config.json`**
   - Place it in the project root.
   - Use the structure described in [Configuration](#configuration).

4. **Run**

```bash
go run main.go -token=<YOUR_TELEGRAM_BOT_TOKEN> -port=<SERVER_PORT>
```

- The bot will:
  - Connect to MongoDB.  
  - Start a long-polling Telegram bot session.  
  - Start an HTTP server on `:<ServerPort>` for webhook-based bots at `/bot/:token`.

5. **(Optional) Configure Webhooks**
   - Point Telegram webhook for additional bots to:  
     `https://<your-domain>/bot/<bot_token>`  
   - The HTTP server will route updates to the shared bot logic.

---

## Project Status & Roadmap

SafeApeBot is an **active work in progress**.

- **Stable / usable today**
  - Core Telegram bot, group admin features (ban/kick/unban), filter system, multi-bot DB model.
- **Under development**
  - Blacklist enforcement, spam mode, verification flows, and mempool/subscription wiring.
  - Integration between the Go service and `GoSniper` contract for `/hp`, `/price`, and related DeFi commands.
- **Planned**
  - Full DeFi command suite (`/hp`, `/price`, `/tokens`, `/gains`, mempool tools).  
  - Wallet watchlists and notification channels.  
  - Cloneable bots via `/clone`.  
  - Comprehensive cache-backed read paths and more granular channel settings.

---

## License

This project is licensed under the **GNU General Public License v3.0 (GPL-3.0)**.  
See [`LICENSE`](./LICENSE) for full details.

