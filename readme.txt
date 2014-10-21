LogVoyage - logging service you always wanted.

Directories:

backend - Heart of LogVoyage server. Main task is to server tons of TCP/IP
		connections, receive logs and sent it to ElasticSearch.
client - Daemon wich runs on user machine. Collects logs and send them to `backend` server.
		Also, can accept logs by TCP/IP or HTTP api. If LogVoyage server is down, `client` daemon will
		store logs to backup file and resend them.
web - http server and UI stuff, etc...