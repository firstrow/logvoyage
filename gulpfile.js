var gulp = require('gulp'),
	concat = require('gulp-concat'),
	less = require('gulp-less'),
	coffee = require('gulp-coffee'),
	path = require('path'),
	del = require('del'),
	watch = require('gulp-watch'),
	addsrc = require('gulp-add-src'),
	clean = require('gulp-clean');

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
	 gulp.src('static/build', {read: false})
        .pipe(clean());

	gulp.src('static/js/*.coffee')
		.pipe(coffee({
			bare: true
		})).on('error', function(e){console.log(e)})
		.pipe(concat("app.js"))
		.pipe(gulp.dest('static/build'));

	gulp.src(js_files)
		.pipe(concat("vendors.js"))
		.pipe(gulp.dest('static/build'));

	gulp.src(['static/less/bootstrap.less'])
		.pipe(less())
		.pipe(gulp.dest('static'))
		.on('end', buildCssTask);

	// Copy chosen image
	gulp.src('static/bower_components/chosen/chosen-sprite@2x.png')
		.pipe(gulp.dest('static/build'));
}

gulp.task('default', defaultTask);