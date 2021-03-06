# mercrawl

`mercrawl` crawls Mercari items of your search condition and send you the result by email.

## Getting started

1. Install Postgresql
2. Set up database. Refer to `migrate.sql`
3. Set up environment variables. See **Environment Variables** below. You can utilize `setenv.sh` template for your convenience.

## Usage

### Try it ASAP!

After set up the environment,

    go run mercrawl/mercrawl.go "sort_order=&keyword=iphone+x&category_root=7&category_child=100&category_grand_child%5B859%5D=1&brand_name=&brand_id=&size_group=&price_min=60000&price_max=&item_condition_id%5B1%5D=1&item_condition_id%5B2%5D=1&status_on_sale=1" & # start crawler
    go run mermail/mermail.go your_mail_addr & # start mailer
    go run rest-api/merest.go & # start rest api server

To quickly quit `mercrawl`, `mermail` and `merest` background processes,

    kill %1 %2 %3

### Crawler

Usage:

`mercrawl <search_condition>`
* `search_condition`: string after https://www.mercari.com/jp/search/?

Example: search on sale PS4 Pro with category of "家庭用ゲーム本体" and price range ¥30,000 ~ ¥45,000

    mercrawl "keyword=ps4+pro&category_root=5&category_child=76&category_grand_child%5B701%5D=1&price_min=30000&price_max=45000&status_on_sale=1"

**WARNING** A too generic search condition that have too many pages of result may cause your IP address banned by Mercari. Please make your search condition be as precise as possible.

### Mailer

Usage:

`mermail <mail_addr>`

You will receive email like this if `mercrawl` successfully scrape new items

![email sent from mermail](img/mermail-example.png "email sent from mermail")

### RESTful API Server

Usage:

`merest`

After the server started, you can access the following resources in JSON:

    GET /items
    GET /item/{id}

## Environment Variables

Global configurations:
* `USER`: database username
* `SSLMODE`: `disable` or `verify-full`

Crawler configurations:
* `PAGE_WORKERS`(optional): max goroutine number for crawling a search result page. Default value is 5
* `ITEM_WORKERS`(optional): max goroutine number for crawling an item page. Default value is 20
* `RECRAWL_INTERVAL`(optional): interval of re-crawl with the same search condition. Default interval is 30 seconds if the variable is not set.

Mailer configurations:
* `INTERVAL`(optional): interval of sending new item info in seconds. Default interval is 30 seconds if the variable is not set.
* `SMTP_SERVER`: mail server address
* `SMTP_PORT`: mail server port
* `SMTP_USER`: mail server login user name
* `SMTP_PWD`: mail server login password

RESTful API Server configurations:
* `REST_PORT`(optional): rest-api server port. default is 8000

## Dependency

* Postgresql

## Performance Tuning

Tune the `itemWorkers` and `pageWorkers` parameters to achieve a better performance for your environment.
