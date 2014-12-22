class LiveLogs
	opts: {
		# Root container of all elements
		container: "#liveLogsContainer"
	}

	show: ->
		$(@opts.container).show()

$ ->
	$("#live_logs").click ->
		live_logs = new LiveLogs()
		live_logs.show()