class LiveLogs
	opts: {
		# Root container of all elements
		container: "#liveLogsContainer"
	}

	init: ->
		# On browser resize keep root container size equal
		$(window).resize =>
			@getContainer().width $(window).width()
			@getContainer().height $(window).height()
		@getContainer().show()
		# Subscribe to new log event
		PubSub.subscribe "log_message", (type, data) =>
			console.log data
			@appendMessage data.message

	appendMessage: (message) ->
		@getContainer().html(message)

	getContainer: ->
		$(@opts.container)

$ ->
	$("#live_logs").click ->
		live_logs = new LiveLogs()
		live_logs.init()