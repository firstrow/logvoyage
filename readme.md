# LogVoyage - fast and simple open-source logging service

LogVoyage is front-end for ElasticSearch. It allows you to store and explore your logs in real-time with friendly web ui.

Note: This is only beta version.

![Dashboard](https://raw.githubusercontent.com/firstrow/logvoyage/master/screenshots/dashboard.png)
![Live logs](https://raw.githubusercontent.com/firstrow/logvoyage/master/screenshots/live-logs.png)

## Installation

### Pre-Requirements.
- ElasticSearch
- Redis

### Installing
Installing LogVoyage is as easy as installing any other go package:
``` bash
go get github.com/firstrow/logvoyage
cd $GOPATH/src/github.com/firstrow/logvoyage
go get ./...
go install
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
By default LogVoyage opens two backend ports accesible to the world.
27077 - TCP port
27078 - HTTP port

#### Sending test messages via telnet

NOTE: Keep in mind to change `apiKey`. You can find your api key at http://localhost:3000/profile page

``` bash
telnet 127.0.0.1 27077
apiKey@logType {"message": "login", "user_id": 1}
apiKey@logType simple text message
```

Now you can see your messages at http://localhost:3000 and try some queries

``` bash
user_id:1
simple*
```

Refer to [ElasticSearch String Query](http://www.elastic.co/guide/en/elasticsearch/reference/1.x/query-dsl-query-string-query.html)
for more info about text queries available.

#### HTTP POST request

Or we can use curl POST request to send messages. Each message should be separated by new line.

``` bash
echo 'This is simple text message' | curl -d @- http://localhost:27078/bulk\?apiKey\=apiKey\&type\=logType
echo '{"message": "JSON format also supported", "action":"test"}' | curl -d @- http://localhost:27078/bulk\?apiKey\=apiKey\&type\=logType
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


<a href='https://pledgie.com/campaigns/28740'><img alt='Click here to lend your support to: LogVoyage and make a donation at pledgie.com !' src='https://pledgie.com/campaigns/28740.png?skin_name=chrome' border='0' ></a>
