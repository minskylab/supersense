# Supersense
Supersense is a simple customizable event board.



### Getting Started

Write a config file as a env variables in a .env file like:
```dotenv
CONSUMER_KEY=<YOUR_TWITTER_CONSUMER_KEY>
CONSUMER_SECRET=<YOUR_TWITTER_CONSUMER_SECRET_KEY>
ACCESS_TOKEN=<YOUR_TWITTER_ACCESS_TOKEN_KEY>
ACCESS_SECRET=<YOUR_TWITTER_ACCESS_SECRET_KEY>

GITHUB_TOKEN=<GITHUB_TOKEN>
```

Event Dispatcher Service:
```shell script
$ PORT=4000 go run cmd/main.go
```

Observer:
```shell script
$ cd observer
$ yarn start
```