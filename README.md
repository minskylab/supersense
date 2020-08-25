# supersense
Supersense is a simple customizable events board, at the moment only offer two event sources: Twitter and Github
test PR

### Getting Started

Write a config file as a env variables in a .env file like:
```dotenv
SS_PORT=4000
SS_GRAPHQL_PLAYGROUND=true

SS_TWITTER_CONSUMER_KEY=<YOUR_TWITTER_CONSUMER_KEY>
SS_TWITTER_CONSUMER_SECRET=<YOUR_TWITTER_CONSUMER_SECRET_KEY>
SS_TWITTER_ACCESS_TOKEN=<YOUR_TWITTER_ACCESS_TOKEN_KEY>
SS_TWITTER_ACCESS_SECRET=<YOUR_TWITTER_ACCESS_SECRET_KEY>
SS_TWITTER_QUERY="#peru"

SS_GITHUB_TOKEN=<GITHUB_TOKEN>
SS_GITHUB_REPOS=minskylab/supersense,minskylab/figport,minskylab/base

SS_DUMMY_PERIOD=1m
```

Event Dispatcher Service:
```shell script
$ PORT=4000 go run cmd/main.go
```

Observer:
```shell script
$ cd observer
$ yarn # in order to download dependencies
$ yarn start
```
