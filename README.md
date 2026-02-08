# DMARC-ANALYSER

A quick tool to analyze your tons of DMARC report files

## Requirements
Docker installed

## How to use

1. Clone the repository
2. Set you DMARC reports folder (input) on `.env` file (see env.example)
3. From project root run

```bash
docker compose up -d
```
4. Navigate to `http://localhost:3000`