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
