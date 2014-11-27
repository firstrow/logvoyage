var Gaze = require('gaze').Gaze;
var gaze = new Gaze('static/**/*.*');
var sys = require('sys')
var exec = require('child_process').exec;

function puts(error, stdout, stderr) {
	sys.puts(stdout)
}

gaze.on('all', function(event, filepath) {
	exec("gulp && sleep 1 && killall gulp", puts);
});
