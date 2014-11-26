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
npm install --save-dev gulp
npm install --save-dev gulp-concat
npm install --save-dev gulp-less 
npm install --save-dev path
npm install --save-dev del
gulp
```
