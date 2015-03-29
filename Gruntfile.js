module.exports = function(grunt) {

  grunt.initConfig({
    pkg: grunt.file.readJSON('package.json'),
    concat: {
      options: {
        separator: "\n"
      },
      javascripts: {
        src: [
          'static/bower_components/jquery/dist/jquery.min.js',
          'static/bower_components/bootstrap/dist/js/bootstrap.min.js',
          'static/bower_components/jquery-jsonview/dist/jquery.jsonview.js',
          'static/bower_components/jquery-cookie/jquery.cookie.js',
          'static/bower_components/datetimepicker/jquery.datetimepicker.js',
          'static/bower_components/bootstrap-multiselect/dist/js/bootstrap-multiselect.js',
          'static/bower_components/ladda-bootstrap/dist/spin.min.js',
          'static/bower_components/ladda-bootstrap/dist/ladda.min.js',
          'static/bower_components/jquery.hotkeys/jquery.hotkeys.js',
          'static/bower_components/sockjs-client/dist/sockjs.js',
          'static/bower_components/pubsub-js/src/pubsub.js',
          'static/js/*.js',
        ],
        dest: 'static/build/all.js'
      },
      css: {
        src: [
          'static/bower_components/bootstrap/dist/css/bootstrap.css',
          'static/bower_components/jquery-jsonview/dist/jquery.jsonview.css',
          'static/bower_components/datetimepicker/jquery.datetimepicker.css',
          'static/bower_components/bootstrap-multiselect/dist/css/bootstrap-multiselect.css',
          'static/bower_components/epoch/epoch.min.css',
          'static/bower_components/chosen/chosen.min.css',
          'static/bower_components/ladda-bootstrap/dist/ladda-themeless.min.css',
        ],
        dest: 'static/build/all.css'
      }
    },
    less: {
      development: {
        files: {
          "static/build/app.css": "static/less/*.less"
        }
      }
    },
    coffee: {
      compileJoined: {
        options: {
          join: true
        },
        files: {
          'static/build/app.js': ['static/js/*.coffee']
        }
      },
    },
    watch: {
      scripts: {
        files: ['static/js/*.*'],
        tasks: ['js'],
        options: {
          spawn: false
        }
      },
      css: {
        files: ['static/less/*.less'],
        tasks: ['css'],
        options: {
          spawn: false
        }
      }
    }
  });

  grunt.loadNpmTasks('grunt-contrib-concat');
  grunt.loadNpmTasks('grunt-contrib-coffee');
  grunt.loadNpmTasks('grunt-contrib-less');
  grunt.loadNpmTasks('grunt-contrib-watch');

  // Default task(s).
  grunt.registerTask('default', ['all']);
  grunt.registerTask('all', ['js', 'css']);
  grunt.registerTask('js', ['concat:javascripts', 'coffee']);
  grunt.registerTask('css', ['concat:css', 'less']);

};
