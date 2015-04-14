# LogVoyage - fast and simple open-source logging service

LogVoyage allows you to store and explore your logs in real-time with friendly web ui.

Note: This is only beta version.

![Dashboard](https://raw.githubusercontent.com/firstrow/logvoyage/master/screenshots/dashboard.png)
![Live logs](https://raw.githubusercontent.com/firstrow/logvoyage/master/screenshots/live-logs.png)

* [![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/firstrow/logvoyage?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
* <a href='https://pledgie.com/campaigns/28740'><img alt='Click here to lend your support to: LogVoyage and make a donation at pledgie.com !' src='https://pledgie.com/campaigns/28740.png?skin_name=chrome' border='0' ></a>
* ![TravisCI](https://api.travis-ci.org/firstrow/logvoyage.svg?branch=master)


<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Installation](#installation)
  - [Pre-Requirements.](#pre-requirements)
  - [Installing](#installing)
- [Usage](#usage)
  - [Sending data to storage](#sending-data-to-storage)
    - [Telnet](#telnet)
    - [Curl](#curl)
- [Third-party clients](#third-party-clients)
- [Submitting a Pull Request](#submitting-a-pull-request)
- [Front-end development](#front-end-development)
  - [Bower](#bower)
  - [Building](#building)
  - [Auto rebuild](#auto-rebuild)
  - [WebSocket messages](#websocket-messages)
- [Roadmap v0.1](#roadmap-v01)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Installation

### Pre-Requirements.
- [ElasticSearch](https://gist.github.com/firstrow/f57bc873cfd6839b6ea8)
- [Redis](http://redis.io/topics/quickstart)

### Installing
Installing LogVoyage is as easy as installing any other go package:
``` bash
go get github.com/firstrow/logvoyage
logvoyage create_users_index
```

## Usage
Once you installed LogVoyage you need to start backend and web servers.
``` bash
logvoyage start-all
```
Or you can start/stop servers separately
``` bash
logvoyage backend
logvoyage web
```
Once server started you can access it at [http://localhost:3000](http://localhost:3000).
Execute `logvoyage help` for more info about available commands.

### Sending data to storage
By default LogVoyage opens two backend ports accesible to the outsise world.

1. 27077 - TCP port
2. 27078 - HTTP port

#### Telnet

```
NOTE: Keep in mind to change `API_KEY` and `LOG_TYPE`.
You can find your api key at http://localhost:3000/profile page.
```

``` bash
telnet 127.0.0.1 27077
API_KEY@LOG_TYPE {"message": "login", "user_id": 1}
API_KEY@LOG_TYPE simple text message
```

Now you can see your messages at http://localhost:3000 and try some queries

#### Curl

Or we can use curl POST request to send messages. Each message should be separated by new line.

``` bash
echo 'This is simple text message' | curl -d @- http://localhost:27078/bulk\?apiKey\=API_KEY\&type\=LOG_TYPE
echo '{"message": "JSON format also supported", "action":"test"}' | curl -d @- http://localhost:27078/bulk\?apiKey\=API_KEY\&type\=LOG_TYPE
```

## Search data
Refer to [Query String Syntax](http://www.elastic.co/guide/en/elasticsearch/reference/1.x/query-dsl-query-string-query.html#query-string-syntax)
for more info about text queries available.

Examples:

``` bash
user_id:1
simple*
amount:>10 and status:completed
```

## Third-party clients
If you know any programming language, you can join our project and implement
LogVoyage client.

## Submitting a Pull Request

1. Propose a change by opening an issue.
2. Fork the project.
3. Create a topic branch.
4. Implement your feature or bug fix.
5. Commit and push your changes.
6. Submit a pull request.

## Front-end development
### Bower
To manage 3rd-party libraries simply add it to static/bower.json and run
```
bower install
```

### Building
We are using grunt to build project js and css files.
Execute next commands to setup environment:
```
npm install
grunt
```
After grunt is done, you can find result files in static/build directory.

### Auto rebuild
To automatically rebuild js, css, coffee, less files simply run in console
```
grunt watch
```

### WebSocket messages
``` coffee
// Sample coffescript code
PubSub.subscribe "log_message", (type, data) ->
  console.log data.message
```

Sample messages:

``` json
{
	"type": "log_message",
	"log_type": "nginx_access",
	"message": "test received log message goes here..."
}
```

``` json
{
	"type": "logs_per_second",
	"count": 5
}
```

## Roadmap v0.1
- Daemons
- Zero-downtime deployment
- Finish web ui
- Docker image
- Docs

## License
LogVoyage is available without any costs under an MIT license. See LICENSE file
for details.
