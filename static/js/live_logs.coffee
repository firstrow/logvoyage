class window.LiveLogs
	opts: {
		# Root container of all elements
		container: "#liveLogsContainer"
		filterContainer: "#liveLogsSearch"
		stackLimit: 2000
	}
	# Root container
	container: null
	autoScroll: true
	messages: []
	filter: null

	constructor: ->
		@container = $(@opts.container)
		@filterContainer = $(@opts.filterContainer)
		@setTheme $.cookie("livelogstheme")

	init: ->
		# On browser resize keep root container size equal
		@container.height $(window).height() - 36
		$(window).resize =>
			@container.height $(window).height() - 36
		@container.scroll @_detectAutoScroll
		# Subscribe to new log event
		PubSub.subscribe "log_message", (type, data) =>
			# TODO: Find out some way to limit number of displayed messages
			@messages.push data
			@appendMessage data.type, data.message
		# Filter events
		@filterContainer.find("input.query").keyup @_filter

	appendMessage: (type, message) ->
		message = @escapeHtml message
		if @filter
			message = @_filterMessage message
		if message
			@container.append "<p><span class='type'>#{type}</span>#{message}</p>"
		@scrollToBottom() if @autoScroll

	clear: ->
		@container.html ''
		@messages = []

	scrollToBottom: ->
		@container.scrollTop @container.prop('scrollHeight')

	switchTheme: ->
		if @container.hasClass "dark"
			@setTheme "light"
		else
			@setTheme "dark"

	setTheme: (t) =>
		if t == "dark" or t == "light"
			@container.removeClass().addClass(t)
			$.cookie("livelogstheme", t)

	_detectAutoScroll: (e) =>
		@autoScroll = (@container.height() + @container.scrollTop()) == @container.prop('scrollHeight')

	_filter: (e) =>
		wait = =>
			@filter = $(e.target).val()
			@_filterAllMessages()
		setTimeout wait, 300

	_filterAllMessages: =>
		@container.html ''
		for data in @messages
			@appendMessage data.type, data.message

	# Returns highlighted text if mached search
	# or false if not
	_filterMessage: (text) =>
		re = new RegExp("(#{@filter})", 'ig')
		if text.match(re)
			return text.replace(re, '<span class="highlight">$1</span>')
		false

	escapeHtml: (unsafe) =>
		unsafe.replace(/&/g, "&amp;")
		.replace(/</g, "&lt;")
		.replace(/>/g, "&gt;")
		.replace(/"/g, "&quot;")
		.replace(/'/g, "&#039;")


