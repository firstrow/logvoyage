class WSocket
	constructor: (apiKey) ->
		@apiKey = apiKey
		@ws = new WebSocket("ws://" + window.location.host + ":12345/ws")
		@ws.onopen = (=> this.register())
		@ws.onmessage = (=> this.onMessage(event))

	onMessage: (event) ->
		data = JSON.parse event.data
		PubSub.publish data.type, data

	register: ->
		@ws.send @apiKey
		console.log "registered user " + @apiKey

$ ->
	new WSocket(options.apiKey)