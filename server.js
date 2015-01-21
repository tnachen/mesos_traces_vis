/*
    Node.js server script
    Required node packages: express, redis, socket.io
*/
const PORT = 3000;
const HOST = 'localhost';

var express = require('express'),
    http = require('http'),
    server = http.createServer(app);

var app = express();

const redis = require('redis');
const client = redis.createClient();
log('info', 'connected to redis server');

const io = require('socket.io');

/*
if (!module.parent) {
    server.listen(PORT, HOST);
    const socket  = io.listen(server);

    socket.on('connection', function(client) {
        const subscribe = redis.createClient()
        subscribe.subscribe('realtime');

        subscribe.on("message", function(channel, message) {
            client.send(message);
            log('msg', "received from channel #" + channel + " : " + message);
        });

        client.on('message', function(msg) {
            log('debug', msg);
        });

        client.on('disconnect', function() {
            log('warn', 'disconnecting from redis');
            subscribe.quit();
        });
    });
}
*/

app.get('/traces', function (req, res) {
    log('msg', 'Asking for /traces');
    client.hgetall('traces', function(err, ts) {
        res.send(ts);
    });
})

app.get('/trace/:trace_id', function(req, res) {
    client.get(req.params.trace_id + "-count", function(err, count) {
        log('debug', 'Found ' + count + ' calls for trace ' + req.params.trace_id);
        var calls = [];
        for (var i = 1; i <= count; i++) {
            log('debug', 'Asking for ' + req.params.trace_id + "-" + i);
            client.hgetall(req.params.trace_id + "-" + i, function(err, cs) {
                log('debug', 'Call returned: ' + cs);
                calls.push(cs);
                if (calls.length == count) {
                  res.send(calls);
                }
            });
        }
    });
});

var server = app.listen(3000, function() {
    var host = server.address().address
    var port = server.address().port

    console.log('Example app listening at http://%s:%s', host, port)
});

function log(type, msg) {
    var color   = '\u001b[0m',
        reset = '\u001b[0m';

    switch(type) {
        case "info":
            color = '\u001b[36m';
            break;
        case "warn":
            color = '\u001b[33m';
            break;
        case "error":
            color = '\u001b[31m';
            break;
        case "msg":
            color = '\u001b[34m';
            break;
        default:
            color = '\u001b[0m'
    }

    console.log(color + '   ' + type + '  - ' + reset + msg);
}