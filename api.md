## Create feed
```
POST /feed
{
    "url": "https://news.ycombinator.com/rss",
}
```

---

## Delete feed
```
DELETE /feed/:id
```

----

## Get all feeds
```
GET /feed
```

```
[
    {"url": "https://news.ycombinator.com/rss"},
    {"url": "https://dou.ua/feed/"}
]
```

---

## Get feed by ID
```
GET /feed/:id
```

Response sample:

```
{"url": "https://news.ycombinator.com/rss"}
```
---

## Get news from the feed

```
GET /feed/:id/news
```


[
{},
]
