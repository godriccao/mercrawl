# mercrawl

mercrawl crawls Mercari items of your search condition and send you the result by email.

## Usage

`mercrawl search_condition`

    search_condition: string after https://www.mercari.com/jp/search/?

example: search on sale PS4 Pro with category of "家庭用ゲーム本体" and price range ¥30,000 ~ ¥45,000, will use string 

    keyword=ps4+pro&category_root=5&category_child=76&category_grand_child%5B701%5D=1&price_min=30000&price_max=45000&status_on_sale=1
