class LiveLogs
	opts: {
		# Root container of all elements
		container: "#liveLogsContainer"
	}
	# Root container
	container: null
	autoScroll: true

	constructor: ->
		@container = $(@opts.container)

	init: ->
		# On browser resize keep root container size equal
		$(window).resize =>
			@container.width $(window).width()
			@container.height $(window).height()
		@container.scroll @_detectAutoScroll
		@container.show()
		# Subscribe to new log event
		PubSub.subscribe "log_message", (type, data) =>
			@appendMessage data.message

	appendMessage: (message) ->
		@container.html(@container.html() + "<p>" + message + "</p>")
		@container.scrollTop(@container.prop('scrollHeight')) if @autoScroll

	_detectAutoScroll: (e) =>
		@autoScroll = (@container.height() + @container.scrollTop()) == @container.prop('scrollHeight')

$ ->
	$("#live_logs").click ->
		live_logs = new LiveLogs()
		live_logs.init()