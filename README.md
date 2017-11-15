# mercrawl

mercrawl crawls Mercari items of your search condition and send you the result by email.

## Getting started

1. Install postgresql
2. Set up the dataabase. Refer to `migrate.sql`
3. Set up the environment variables. See **Environment Variables** below.

## Usage

### Crawler

`mercrawl search_condition`
* `search_condition`: string after https://www.mercari.com/jp/search/?

Example: search on sale PS4 Pro with category of "家庭用ゲーム本体" and price range ¥30,000 ~ ¥45,000

    mercrawl "keyword=ps4+pro&category_root=5&category_child=76&category_grand_child%5B701%5D=1&price_min=30000&price_max=45000&status_on_sale=1"

### Mailer

`mermail <mail_addr>`

## Environment Variables

Database related:
* `USER`: database username
* `DBNAME`: database name
* `SSLMODE`: should be `disable` or `verify-full`

Mailer related:
* `INTERVAL`: seconds interval of sending new item info. Default interval is 5 seconds if the variable is not set.


## Requirements

* Postgresql

## Tuning

Tune the `itemWorkers` and `pageWorkers` parameters to achieve a better performance for your environment.