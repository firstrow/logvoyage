var gulp = require('gulp'),
	concat = require('gulp-concat'),
	less = require('gulp-less'),
	coffee = require('gulp-coffee'),
	path = require('path'),
	del = require('del'),
	watch = require('gulp-watch'),
	addsrc = require('gulp-add-src');

var js_files = [
	'static/bower_components/jquery/dist/jquery.min.js',
	'static/bower_components/bootstrap/dist/js/bootstrap.min.js',
	'static/bower_components/jquery-jsonview/dist/jquery.jsonview.js',
	'static/bower_components/datetimepicker/jquery.datetimepicker.js',
	'static/bower_components/bootstrap-multiselect/dist/js/bootstrap-multiselect.js',
	'static/bower_components/d3/d3.min.js',
	'static/bower_components/epoch/epoch.min.js',
	'static/bower_components/chosen/chosen.jquery.js',
	'static/bower_components/ladda-bootstrap/dist/spin.min.js',
	'static/bower_components/ladda-bootstrap/dist/ladda.min.js',
	'static/bower_components/jquery.hotkeys/jquery.hotkeys.js',
	'static/bower_components/sockjs-client/dist/sockjs.js',
	'static/bower_components/pubsub-js/src/pubsub.js',
	'static/js/*.js',
]

var css_files = [
	'static/bootstrap.css',
	'static/bower_components/jquery-jsonview/dist/jquery.jsonview.css',
	'static/bower_components/datetimepicker/jquery.datetimepicker.css',
	'static/bower_components/bootstrap-multiselect/dist/css/bootstrap-multiselect.css',
	'static/bower_components/epoch/epoch.min.css',
	'static/bower_components/chosen/chosen.min.css',
	'static/bower_components/ladda-bootstrap/dist/ladda-themeless.min.css',
	'static/css/custom.css',
]

var buildCssTask = function() {
	gulp.src(css_files)
		.pipe(concat('all.css'))
		.pipe(gulp.dest('static/build'));
}

var defaultTask = function() {
	// Compile coffee-script files
	gulp.src('static/js/*.coffee')
		.pipe(coffee({
			bare: true
		}))
		.pipe(addsrc(js_files))
		.pipe(concat("all.js"))
		.pipe(gulp.dest('static/build'));

	// Less and css
	gulp.src(['static/less/bootstrap.less'])
		.pipe(less())
		.pipe(gulp.dest('static'))
		.on('end', buildCssTask);

	gulp.src('static/bower_components/chosen/chosen-sprite@2x.png')
		.pipe(gulp.dest('static/build'));
}

gulp.task('default', defaultTask);