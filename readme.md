# LogVoyage - fast and simple logging service.

TODO: Screenshot

## Installation

### Pre-Requirements.
- Redis
- ElasticSearch

### Installing LogVoyage is as easy installing any other go package:
``` bash
> go get github.com/firstrow/logvoyage
```

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

# WebSocket messages
``` coffee
// Sample coffescript code
PubSub.subscribe "log_message", (type, data) ->
  console.log data.message
```
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