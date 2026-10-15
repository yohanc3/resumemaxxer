ResumeMaxxer alerts you whenever a new job is posted on Simplify, and tailors your resume to the job description.

You are notified through Telegram. Every alert includes:
1. A tailored resume (PDF)
2. Overleaf link to the tailored resume so you can edit it
3. Job posting info such as: company, role, and the link to apply

> [!IMPORTANT]
> This project is still in development, without a stable version deployed.

## How to get started
1. Ping the Telegram Bot by clicking here, or search "ResumeMaxxerBot" and type
`/start`
2. A unique link will be issued, so you can input your resume.
3. The Bot will let you know you're registered, meaning you will start getting alerts
   from then on.

## Purpose
In my experience, tailoring resumes can take quite some time. ResumeMaxxer helps
with: 

1. Tailoring resumes to job descriptions without hallucinating the user's experience/skills.
ATS systems will rank applications not matching the keywords lower.
2. Saving time visiting Simplify's internship repo, which can be a huge time sink.

## How your resume is modified  

The LLM tailoring your resume is allowed to make only 3 types of modifications:

1. Google's XYZ Formatting: Aligning text to the standard XYZ structure.
2. Keyword Alignment: Adding explicit terms/skills found in the job description 
    - **Example: "Kubernetes" where your resume may only have "K8s"**
3. Cohesive Reshaping: Pivot framing of specific experiences to match the target role without lying
    - **Example: "Setup Dockerfile to deploy app." -> "Containerized application environments to standardize production deployment pipelines."**

A static check automatically rejects bullet points exceeding 180 characters (roughly 1.5 lines when compiled), triggering a retry automatically. This is the sole validator for the MVP.

## How to setup locally 

1. Create a .env file, and fill it out using .env.example for guidance.

2. Run the containers:

```sh
docker compose up 
```

## Send one Telegram text message

This standalone command sends `Hello from ResumeMaxxer!` once. It requires Go
(the Docker builds use Go 1.25), but no database or running containers.

1. Create a bot with [@BotFather](https://t.me/BotFather) using `/newbot`, or use
   the token for your existing bot. Keep the token private.
2. Open a private chat with **your bot** and send `/start`.
3. Load your token into your shell without putting it in shell history:

   ```bash
   read -r -s -p 'Telegram bot token: ' TELEGRAM_BOT_TOKEN
   echo
   export TELEGRAM_BOT_TOKEN
   ```

4. Run the following lookup yourself. It calls Telegram's `getUpdates` endpoint
   and prints only private chat IDs and message text, never the token. It requires
   Python 3 and uses the exported token:

   ```bash
   python3 - <<'PY'
   import json
   import os
   import urllib.request

   endpoint = "https://api.telegram.org/bot" + os.environ["TELEGRAM_BOT_TOKEN"] + "/getUpdates"
   try:
       with urllib.request.urlopen(endpoint, timeout=15) as response:
           data = json.load(response)
       if not data.get("ok"):
           raise ValueError("Telegram API rejected the lookup")
       for update in data.get("result", []):
           message = update.get("message", {})
           chat = message.get("chat", {})
           if chat.get("type") == "private":
               print("chat_id:", chat["id"], "text:", message.get("text", ""))
   except Exception:
       raise SystemExit("Chat ID lookup failed; check your token and network connection.")
   PY
   ```

   Copy the `chat_id` associated with your `/start` message. If there are no
   results, send `/start` again and repeat the lookup. Updates are retained for
   at most 24 hours. `getUpdates` cannot be used while a webhook is configured;
   another polling process can also consume the updates.
   See the [Telegram Bot API documentation](https://core.telegram.org/bots/api#getupdates).

5. Export that ID, then run the sender from the repository root:

   ```bash
   export TELEGRAM_CHAT_ID='123456789'
   go run ./cmd/telegram-send
   ```

The Go command reads the process environment directly; it does **not**
automatically load `.env`. Alternatively, add the two Telegram values to your
local `.env` using shell-compatible assignments, then run this single command
from the repository root in Bash:

```bash
(set -a; source ./.env; set +a; go run ./cmd/telegram-send)
```

`set -a` exports the variables loaded by `source`; parentheses keep them scoped
to the subshell. Only source a trusted local file. Keep token values out of commits.

The command exits unsuccessfully for missing configuration, network failures,
HTTP failures, or Telegram API errors. It has a 15-second HTTP timeout and does
not retry. It does not print the token or raw HTTP errors/responses.
