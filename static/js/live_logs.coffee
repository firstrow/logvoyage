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
	addedMessages: 0
	filter: null

	constructor: ->
		@container = $(@opts.container)
		@filterContainer = $(@opts.filterContainer)

	init: ->
		# On browser resize keep root container size equal
		@container.height $(window).height() - 28
		$(window).resize =>
			@container.height $(window).height() - 28
		@container.scroll @_detectAutoScroll
		# Subscribe to new log event
		PubSub.subscribe "log_message", (type, data) =>
			@appendMessage data.type, data.message
		# Filter events
		@filterContainer.find("input.query").keyup @_filter

	appendMessage: (type, message) ->
		# @filter = if filter then new RegExp "(#{filter})", 'ig' else null
		message = @escapeHtml message
		cls = ""
		if @filter
			_msg = @_filterMessage message
			message = if _msg then _msg else message
			cls = if _msg then "" else "hidden" 

		@container.append("<p class='#{cls}'><span class='type'>#{type}</span>#{message}</p>")

		@addedMessages++
		if @addedMessages == @opts.stackLimit
			console.log "stack limit reached"
			@container.find("p").slice(0, 1).remove()
			@addedMessages--
		@container.scrollTop(@container.prop('scrollHeight')) if @autoScroll

	_detectAutoScroll: (e) =>
		@autoScroll = (@container.height() + @container.scrollTop()) == @container.prop('scrollHeight')

	_filter: (e) =>
		wait = =>
			@filter = $(e.target).val()

			if @filter == ""
				$("#{@opts.container} p").removeClass("hidden")
			else
				$(@opts.container).find("p").each @_filterAllMessages
		$(@opts.container).find(".highlight").each @_removeHighlight
		setTimeout wait, 300

	_removeHighlight: (index, el) =>
		$(el).html($(el).text())

	_filterAllMessages: (index, el) =>
		result = @_filterMessage $(el).html()
		if result 
			$(el).html result
			$(el).removeClass "hidden"
		else
			$(el).addClass "hidden"

	# Returns highlughted text if mached search
	# or false if not
	_filterMessage: (text) =>
		re = new RegExp("(#{@filter})", 'ig')
		if text.match(re)
			return text.replace(re, '<span class="highlight">$1</span>')
		false

	escapeHtml: (unsafe) ->
		unsafe.replace(/&/g, "&amp;")
		.replace(/</g, "&lt;")
		.replace(/>/g, "&gt;")
		.replace(/"/g, "&quot;")
		.replace(/'/g, "&#039;")


