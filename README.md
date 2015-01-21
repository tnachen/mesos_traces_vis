This assumes all the tracing data is installed into a local redis instance.

To ingest local data into redis:

tail -f -c +1 <path> | python -u populate_trace_db.py

To install the Node.js server part:

$ npm install
```

- Start the redis server (default localhost on port 6379)
- Start node with server.js

```shell
    $ node server.js
```

- Open client.html in browser code

```shell
    $ open client.html
```





