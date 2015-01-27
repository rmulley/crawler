###How to crawl URLs
To submit URLs to crawl, you'll need to POST a valid JSON to the server like below.

```json
curl http://192.168.59.103:8080 \
-X POST \
-H 'Accept: application/json' \
-H 'Content-Type: application/json' -d '
{
  "urls": [
    "http://www.docker.com",
    "http://www.cnn.com",
    "http://www.yahoo.com/"
  ]
}'
```

###To check the status of jobs running

```curl http://192.168.59.103:8080/status/1```

```json
{
  "job_id": 1,
  "completed": 61,
  "in_progress": 0
}
```


###To view crawling results

```curl http://192.168.59.103:8080/result/1```

```json
{
    "job_id": 1,
    "results": 
    [
      {
        "url": "http://www.docker.com",
        "images": 
        [
          "https://d3oypxn00j2a10.cloudfront.net/0.14.0/img/nav/docker-logo-loggedout.png",
          "https://d3oypxn00j2a10.cloudfront.net/0.14.0/img/homepage/dhe-shot.png",
          "https://d3oypxn00j2a10.cloudfront.net/assets/img/Gilt/Gilt_Logo.jpg",
          "https://d3oypxn00j2a10.cloudfront.net/assets/img/Yelp/Yelp-Logo.jpg",
          "https://d3oypxn00j2a10.cloudfront.net/0.14.0/img/universal/official-repository-icon.png",
          "https://d3oypxn00j2a10.cloudfront.net/0.14.0/img/universal/trusted-icon.svg",
          "https://d3oypxn00j2a10.cloudfront.net/0.14.0/img/universal/official-repository-icon.png",
          "https://d3oypxn00j2a10.cloudfront.net/0.14.0/img/universal/trusted-icon.svg",
          "https://d3oypxn00j2a10.cloudfront.net/0.14.0/img/universal/official-repository-icon.png",
          "https://d3oypxn00j2a10.cloudfront.net/0.14.0/img/universal/trusted-icon.svg",
          "https://d3oypxn00j2a10.cloudfront.net/0.14.0/img/universal/official-repository-icon.png",
          "https://d3oypxn00j2a10.cloudfront.net/0.14.0/img/universal/official-repository-icon.png",
          "https://d3oypxn00j2a10.cloudfront.net/0.14.0/img/universal/official-repository-icon.png",
          "https://d3oypxn00j2a10.cloudfront.net/0.14.0/img/universal/official-repository-icon.png",
          "https://d3oypxn00j2a10.cloudfront.net/0.14.0/img/universal/official-repository-icon.png",
          "https://pbs.twimg.com/profile_images/2620852811/owa81mf2uyho1cioeccz_normal.jpeg",
          "https://pbs.twimg.com/profile_images/2738465953/27e50cf1e0e0c332bbe649db0985fb21_normal.jpeg",
          "https://pbs.twimg.com/profile_images/529629187506585600/kzlDlPED_normal.jpeg",
          "https://pbs.twimg.com/profile_images/529295351656624128/_FuL-Jzb_normal.png"
        ]
      },
      {
        "url": "http://www.docker.com/whatisdocker",
        "images": 
        [
          "https://d3oypxn00j2a10.cloudfront.net/0.14.0/img/nav/docker-logo-loggedout.png"
        ]
      },
      .
      .
      .
```
