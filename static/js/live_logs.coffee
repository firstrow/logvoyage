class window.LiveLogs
	opts: {
		# Root container of all elements
		container: "#liveLogsContainer"
	}
	# Root container
	container: null
	autoScroll: true
	addedMessages: 0

	constructor: ->
		@container = $(@opts.container)

	init: ->
		# On browser resize keep root container size equal
		@container.height $(window).height()
		$(window).resize =>
			@container.height $(window).height()
		@container.scroll @_detectAutoScroll
		# Subscribe to new log event
		PubSub.subscribe "log_message", (type, data) =>
			@appendMessage data.type, data.message

	appendMessage: (type, message) ->
		@container.append("<p>" + "<span class='type'>" + type + "</span>"  + message + "</p>")
		@container.scrollTop(@container.prop('scrollHeight')) if @autoScroll
		@addedMessages++
		if @addedMessages >= 120
			@container.find("p").slice(0, 20).remove()
			@addedMessages = 100

	_detectAutoScroll: (e) =>
		@autoScroll = (@container.height() + @container.scrollTop()) == @container.prop('scrollHeight')