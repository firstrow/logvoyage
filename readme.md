LogVoyage - logging service you always wanted.

## Working with frontend
###Bower
To manage 3rd-party libraries simply add it to static/bower.json and run
```
bower install
```

###Building
We are using gulp to build project js and css files.
Execute next commands to setup environment:
```
npm install --global gulp
npm install --global gaze
npm install --save-dev gulp
npm install --save-dev gulp-concat
npm install --save-dev gulp-less 
npm install --save-dev gulp-coffee
npm install --save-dev gulp-add-src
npm install --save-dev gulp-clean
npm install --save-dev path
npm install --save-dev del
gulp
```
After gulp is done you can find result files in static/build directory.

### Auto rebuild  
To automatically rebuild js & css files simply run in console
```
node watch.js
```

# WebSocket messages
``` json
{
	"type": "log_message",
	"message": "test received log message goes here..."
}
```

``` json
{
	"type": "logs_per_second",
	"count": 5
}
```