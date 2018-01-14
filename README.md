# go-net-demo
Intended as a showcase of goroutines and, consequently, connection handling

First commit: a client/server that demonstrate a simple "long pulling" approach.

Second commit should demonstrate the ease of keeping 10k connections alive.
For this, the server provides two methods:
* `POST /anything` will take a connection and keep it, till ...
* `DELETE /anything` makes all the hanging post requests reply with the the number of requests at the time of calling.

The client will simply do 10_000 POST, then a DELETE.

## Running

The first time I tried even 1000 connections, I got "too many open files" error.

These can be checked by `ulimit -Sa`, you'd need to [lift it](https://ro-che.info/articles/2017-03-26-increase-open-files-limit)
to about 11000.

## Results

    2018/01/14 12:28:41 sent 100 POST requests
    ...
    2018/01/14 12:28:43 sent 1100 POST requests
    ...
    2018/01/14 12:28:48 sent 1200 POST requests
    ...
    2018/01/14 12:28:49 sent 10000 POST requests
    2018/01/14 12:28:49 received replies from 100 POST requests
    2018/01/14 12:28:49 received replies from 200 POST requests
    2018/01/14 12:28:49 received replies from 300 POST requests
    ...
    2018/01/14 12:28:49 received replies from 9800 POST requests
    2018/01/14 12:28:49 received replies from 9900 POST requests
    2018/01/14 12:28:49 received replies from 10000 POST requests
    2018/01/14 12:28:49 received replies from all 10000 POST requests
    2018/01/14 12:28:49 done

This is currently my answer to the "Why Go(lang)?" question.
This result is non-straightforward for most other languages.
For many you need to do some trickery to get this far.
