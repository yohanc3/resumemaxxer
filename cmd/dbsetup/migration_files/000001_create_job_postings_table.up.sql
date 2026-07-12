CREATE TABLE IF NOT EXISTS job_postings(
    id VARCHAR (40) PRIMARY KEY,
    source TEXT,
    category TEXT,
    company_name TEXT,
    title TEXT,
    active BOOLEAN DEFAULT false,
    terms TEXT[],
    date_updated TIMESTAMP,
    date_posted TIMESTAMP,
    url TEXT,
    locations TEXT[],
    company_url TEXT,
    is_visible BOOLEAN DEFAULT false,
    sponsorship TEXT,
    degrees TEXT[]
);
