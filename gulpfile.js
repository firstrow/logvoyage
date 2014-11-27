var gulp = require('gulp');
var concat = require('gulp-concat');
var less = require('gulp-less');
var path = require('path');
var del = require('del')
var watch = require('gulp-watch');

var js_files = [
	'static/bower_components/jquery/dist/jquery.min.js',
	'static/bower_components/bootstrap/dist/js/bootstrap.min.js',
	'static/bower_components/jquery-jsonview/dist/jquery.jsonview.js',
	'static/bower_components/datetimepicker/jquery.datetimepicker.js',
	'static/bower_components/bootstrap-multiselect/dist/js/bootstrap-multiselect.js',
	'static/bower_components/d3/d3.min.js',
	'static/bower_components/epoch/epoch.min.js',
	'static/bower_components/chosen/chosen.jquery.js',
	'static/js/*',
]

var css_files = [
	'static/bootstrap.css',
	'static/bower_components/jquery-jsonview/dist/jquery.jsonview.css',
	'static/bower_components/datetimepicker/jquery.datetimepicker.css',
	'static/bower_components/bootstrap-multiselect/dist/css/bootstrap-multiselect.css',
	'static/bower_components/epoch/epoch.min.css',
	'static/bower_components/chosen/chosen.min.css',
	'static/css/custom.css',
]

var buildCssTask = function() {
	gulp.src(css_files)
		.pipe(concat('all.css'))
		.pipe(gulp.dest('static/build'));
}

var defaultTask = function() {
	// Compile less
	gulp.src(js_files)
		.pipe(concat('all.js'))
		.pipe(gulp.dest('static/build'));

	// Less and css
	gulp.src(['static/less/bootstrap.less'])
		.pipe(less())
		.pipe(gulp.dest('static'))
		.on('end', buildCssTask);
}

gulp.task('default', defaultTask);
