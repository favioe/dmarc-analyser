# DMARC-ANALYSER

A quick tool to analyze your tons of DMARC report files

## Requirements
Docker installed

## How to use

1. Clone the repository
2. Create a `.env` file. You can create it based on [`env.example`](env.example) 
3. Set you `DMARC_DIR` reports folder (input) on `.env` file (see env.example)
4. Download your tons of DMARC reports sent by Gmail, Hotmail, etc. and place them inside `DMARC_DIR` folder
5. From project root run

```bash
docker compose up -d
```
6. Navigate to `http://localhost:3000` (set your ports on `.env` file)

## Adding more input files

1. Add the new DMARC report files to `DMARC_DIR` folder
2. Go to [http://localhost:3000](http://localhost:3000) and press "Refresh Report"


# Screenshots

On the example, multiple `.gz` and `.zip` files downloaded from DAMRC report files sent by Google, Hotmail, etc. were added to `DMARC_DIR` folder

You don't need to uncompress them, just download them from your email client, drop them in the `DMARC_DIR` folder and press "Refresh report" button

![Report Output example](./docs/screenshot1.jpg)
