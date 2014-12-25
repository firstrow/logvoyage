class LiveLogs
	opts: {
		# Root container of all elements
		container: "#liveLogsContainer"
	}
	# Root container
	container: null

	constructor: ->
		@container= $(@opts.container)

	init: ->
		# On browser resize keep root container size equal
		$(window).resize =>
			@container.width $(window).width()
			@container.height $(window).height()
		@container.show()
		# Subscribe to new log event
		PubSub.subscribe "log_message", (type, data) =>
			@appendMessage data.message

	appendMessage: (message) ->
		@container.html(message)

$ ->
	$("#live_logs").click ->
		live_logs = new LiveLogs()
		live_logs.init()