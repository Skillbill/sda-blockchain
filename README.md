# sda-blockchain
Blockchain for Developers SDA

Slides are in [present](https://godoc.org/golang.org/x/tools/present) format.

For your convenience the [generated html page](https://skillbill.github.io/sda-blockchain/slides.html) is provided as well.

## toychain demo

### build
Install the go tool using your distro package manager (look for something like golang or go-lang).

cd to the repo's root and run:
```
$ make 
```
If you're running *BSD, you already know you have to use `gmake` instead.
### run
```
$ bin/toychain -v 2>/tmp/toychain.log
```
On another terminal you can monitor the log with:
```
$ tail -f /tmp/toychain.log
```
### commands
```
> help

block last              # print the last block stats
block i                 # print the ith block stats

pool                    # show nodes overview
pool add                # add a node to the pool
pool add n              # add n nodes to the pool
pool del i              # remove the ith node from the pool
pool del all            # remove every node from the pool

pool run|stop i         # run|stop the mining for the ith node in the pool
pool run|stop all       # run|stop the mining for every node in the pool

pay src dst amount      # make a tx of the given amount from address src to address dst

diff                    # show PoW current difficulty
diff n                  # change PoW difficulty to n
stat                    # show toychain status

exit                    # shutdown the toychain and exit
help                    # you already know this one

```
