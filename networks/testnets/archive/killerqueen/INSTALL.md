# Full Node

### Prerequisites:

Recommended hardware will depend on use cases and desired functionalities. However, the minimum specifications **required** are as follows: 

2+ vCPU

4+ GB RAM

120+ GB SSD

### Installing Node:

Before installing your node, please verify you meet all the prerequisites to participating in Killer Queen above. 

Operators can install the `quicksilverd` binary from the source.

To install the `quicksilverd` binary from source, you will need to have Golang [1.17+](https://golang.org/dl/) installed on your OS.

Once you have your environment setup correctly, clone the Quicksilver repository and install the binary:

```go
$ git clone --branch v0.4.0 https://github.com/ingenuity-build/quicksilver.git

$ cd quicksilver && make install
```

The `quicksilverd` binary will be installed into your `$GOPATH/bin` directory.

Finally, verify your quicksilverd version:

```
quicksilverd version
```

Verify you are using **v0.4.0** for the `killerqueen-1` testnet.
