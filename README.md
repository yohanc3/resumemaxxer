ResumeMaxxer creates tailored resumes and alerts you whenever a new job is posted on Simplify.
You are notified through Telegram with a `Tailored PDF resume`, `Overleaf link to your tailored
resume`, `Job posting info` such as: company, role, link to apply.

> [!IMPORTANT]
> This project is still in development, without a stable version deployed.

## How to get started
1. Ping the Telegram Bot by clicking here, or search "ResumeMaxxerBot" and type
`/start`
2. It will give you a unique link, where you will only input your name, email,
   and resume. Click "Done" and exit the browser.
3. The Bot will let you you're registered, and you will start getting alerts
   from then on.

## Purpose
Applying for jobs takes me a lot of time, given I have to:

1. Modify my resume to match the job description, so the ATS doesn't rank me lower if the keywords
don't match.
2. Constantly visiting Simplify's internship repo is a time sink as well.

ResumeMaxxer aims to fix these pain points. It's easy to set up, and the generated resumes only
modify as little as needed for your resume to 

This is a **huge** time sink.

ResumeMaxxer fixes this.

## Jake's Template & NO fully AI generated resumes
I don't like fully AI generated resumes, nor do employers.

Do they work sometimes? Yes. Are they easy to spot, and looked upon? Also yes.

So how does ResumeMaxxer avoid AI slop? 
1. The harness instructs the LLM to follow Google's XYZ method as far as possible. Meaning,
   if your bullet-points don't have metrics, then the LLM will not make them up. The harness
   gives the LLM a few options when refactoring bullet-points:

2. The harness gives the LLM a strict set of options when refactoring bullet-points:
    - Format the bullet-point according to Google's XYZ method.
    - Mention technologies or skills that already exist elsewhere on your resume (typically in your skills section)
      when they're implied but omitted from a bullet point. This helps ensure important ATS keywords appear
      in the experience that demonstrates them, without hallucinating technologies or accomplishments. For example,
      if you have xDS in your resume, but the job description specifically mentions "Control Plane", we will add this.
    - Reshape the cohesiveness of your bullet-points, so that it's way easier to 
      spot the work you have done, and the work you are looking for. 
        - Think of resumes with strong ML experiences, but if the role is SWE, you would not seem like a good match.
          ResumeMaxxer fixes this by ensuring your ML experiences are a good match for SWE roles, without lying or hallucinating.
    - Fix typos, etc.

3. ResumeMaxxer uses Jake's Template. The LLM doesn't concern itself with this.
   It is only in charge of the bullet-points. We construct and compile the LaTex into PDF.

4. Static validation checks. There are validators that will reject
   bullet-points if they:

    - Are above 180 characters long (close to 2 full lines when compiled)
    - This is the only validation check as of now. More will be added later,
      after the MVP is out.
