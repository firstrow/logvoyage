var Gaze = require('gaze').Gaze;
var gaze = new Gaze('static/**/*.*');
var sys = require('sys')
var exec = require('child_process').exec;

function puts(error, stdout, stderr) {
	sys.puts(stdout)
}

function runGulpRun() {
	exec("gulp && sleep 1 && killall gulp", puts);
}

gaze.on('all', function(event, filepath) {
	runGulpRun();
});
