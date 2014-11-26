var gulp = require('gulp');
var concat = require('gulp-concat');
var less = require('gulp-less');
var path = require('path');

var js_files = [
	'static/bower_components/jquery/dist/jquery.min.js',
	'static/bower_components/bootstrap/dist/js/bootstrap.min.js',
	'static/bower_components/jquery-jsonview/dist/jquery.jsonview.js',
	'static/bower_components/datetimepicker/jquery.datetimepicker.js',
	'static/bower_components/bootstrap-multiselect/dist/js/bootstrap-multiselect.js',
	'static/bower_components/d3/d3.min.js',
	'static/bower_components/epoch/epoch.min.js',
	'static/bower_components/chosen/chosen.jquery.js',
	'static/js/init.js',
	'static/js/view.js',
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

gulp.task('default', function() {
	gulp.src(js_files)
		.pipe(concat('all.js'))
		.pipe(gulp.dest('static/build'));

	// Less and css
	gulp.src(['static/less/bootstrap.less'])
		.pipe(less())
		.pipe(gulp.dest('static'));

	gulp.src(css_files)
		.pipe(concat('all.css'))
		.pipe(gulp.dest('static/build'));
});