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

	getContainer: ->
		$(@opts.container)

$ ->
	$("#live_logs").click ->
		live_logs = new LiveLogs()
		live_logs.init()