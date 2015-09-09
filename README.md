# eth-client
Command line interfaces for low-level conversations with ethereum chain(s)

Inspired by https://github.com/eris-ltd/mint-client

# Overview

In current form, the `eth-client` has three tools: 

```
ethgen - for creating genesis.json files
ethinfo - for querying an ethereum chain for information
ethtx - for crafting and broadcasting transactions
```

# Install

I will presume you have [go installed](https://golang.org/doc/install),
and that you have set your `$GOPATH` and put `$GOPATH/bin` on your `$PATH`. 

To install the tools, just run 

```bash
go get github.com/eris-ltd/eth-client/...
```

If you want to use the tools on a local node, you should [install `geth`](https://github.com/ethereum/go-ethereum/wiki/Building-Ethereum).

You will also need the `eris-keys` server for generating keys and signing transactions.

```bash
go get github.com/eris-ltd/eris-keys
```

# eris-keys

The key server is meant to run as a persistent process in the background, and play much the same role as a GPG daemon (eventually, hopefully, we can replace it with a simple wrapper around GPG).

To start the server, run

```bash
eris-keys server &
```

It will bind to `localhost:4767`

Now, let's generate a key.

```bash
ADDR=`eris-keys gen --type=secp256k1,sha3`
echo $ADDR
```

We save the `$ADDR` variable for later convenience. 
Note we are also going to need it in a new terminal window soon, so either copy it to your clipboard or save it to your ~/.bashrc

# Ethereum Test Chains

Let's boot up a test chain and send some transactions to it. 
First we need to create a root directory for the chain, and give it a genesis.json with our address:

```bash
mkdir ~/.myethereum
ethgen $ADDR > ~/.myethereum/genesis.json
```

Splendid! Take a look at that genesis file if you like with `cat ~/.myethereum/genesis.json` .
Note that `ethgen` chooses reasonable default values for the difficulty and the gas-limit for running test chains. 
If you want to customize, just use the flags. See `ethgen --help` for more.

Now it's time to boot the chain.

But first note the unfortunate reality of ethereum's proof of work means you need at least 1GB of free RAM for mining to work. 
It may take some time for mining to get started (on the order of a few minutes). 
If you'd like to try an EVM compatible blockchain that isn't so resource intensive, checkout [erisdb](https://github.com/eris-ltd/eris-db) ;)

```bash
geth --datadir ~/.myethereum --rpc --mine --genesis ~/.myethereum/genesis.json --maxpeers 0 --etherbase $ADDR --verbosity 7
```

This will start a private ethereum chain on your machine. 

Check the chain's status with

```bash
ethinfo status
```

# Ethereum Transactions

Typically, one would now run `geth attach` in another window, and be presented with a javascript console for querying and transacting on the blockchain.

Let's instead use our golang command line client!

Of course, since `geth` is now running, we should use a new window. Make sure to copy the ADDR variable and set it again in the new window:

```
export ADDR=<address we used>
```

First, we'll generate another account to send funds to:

```bash
ADDR2=`eris-keys gen --type=secp256k1,sha3`
```

Now we can craft a transaction:

```bash
ethtx send --addr=$ADDR --to=$ADDR2 --amt=10 --gas=21000 --price=100000000000 --sign --broadcast
```

The command should execute successfuly, returning the transaction ID.

Now wait a few seconds for the block to be committed. You can monitor the block number with `ethinfo status`.

Now we can get the transaction receipt

```bash
ethinfo receipt <transaction ID>
```

and check on the account:

```bash
ethinfo account $ADDR2
```

The balance should be `0xa` (ie. `10`)!

Ok, let's break down the `ethtx` command a little bit. Ethereum only has one official transaction type, but it serves three distinct purposes. 
You can simply send funds from one account to another, or you can create a contract, or you can call a contract.
To reflect this, `ethtx` has three main commands: `send`, `create`, and `call`.
Over time, `ethtx` will incorporate commands for talking to the major dapps, to facilitate interactions with them. The first such dapp will be the name reg, and you'll be able to use `ethtx name` to register a new name there.

I should note, the `ethtx` flags accept both hex and base 10 numbers. If you are using hex, make sure to prefix with `0x`.

You can also specify a transaction's nonce with the `--nonce` flag. If no nonce is specified, the correct one is fetched from the blockchain. Note however that at the current time this mechanism only supports sending one transaction from a given address per block.

The `--sign` and `--broadcast` flags allow you to specify exactly what you want to do. 
Maybe you only want to craft the bytes for the transaction now and sign it later, or maybe sign it now and broadcast later? 
Or maybe you want to do everything now, in which case both `--sign` and `--broadcast` are appropriate. 
Soon, we will add a `--wait` feature so you can wait until the transaction is actually committed in a block.

You can also add the `--binary` flag to print the hex encoded rlp serialization of the transaction. 
For example, if you are signing the transaction offline, you might do:

```bash
ethtx send --addr=$ADDR --to=$ADDR2 --nonce=3 --amt=10 --gas=21000 --price=100000000000 --sign --binary
```

Note how we specify the nonce. Using `--binary` gets us the transaction bytes, which we can copy and broadcast once we're back online with 

```bash
ethinfo broadcast <transaction bytes>
```

# Ethereum Contracts

Time to deploy a contract. You will need some ethereum byte code. Here is the bytecode for the simplest transaction imagineable:

```bash
0x6005600055
```

In assembly, this would look like:

```asm
PUSH1 0x5 PUSH1 0x0 SSTORE
```

The net result of this contract is that the number `0x5` gets stored at position `0x0`.

Let's deploy it:

```bash
ethtx create --addr=$ADDR --code=0x6005600055 --amt=0 --gas=50000 --price=100000000000 --sign --broadcast
```

Note how it prints the address of the newly created contract on the last line.

Until we implement `--wait`, we'll have to wait a few seconds to ensure a block is mined. Then we can check the storage:

```bash
ethinfo storage <new address>
```

Lo and behold, we see `0x5` stored at `0x0` !

Alternatively, we can ask for just the value at `0x0`:

```bash
ethinfo storage <new address> 0x0
```

Of course this is a trivial contract that is now useless, but it demonstrates the basics of using these tools.

If you want to compile solidity, check out the lovely-little-languages compiler server at https://github.com/eris-ltd/lllc-server.

We will also be wrapping the eth-client tools at a slightly higher level so that solidity compiling and abi formatting may be done for you.

The purpose of the `eth-client` itself however is a simple, low-level interface for developers.

Enjoy!

# Tips

By setting the `ETHTX_ADDR` environment variable, you can avoid passing the `--addr` flag.

There are also `ETHTX_SIGN_ADDR` and `ETHTX_NODE_ADDR` environment variables to set the address of the signing daemon and the node itself (since keys are managed by the signing daemon, it becomes reasonable to set up an ethereum node whose rpc is bound to the public internet - ethereum rpc as a service, if you will).

# Live Ethereum Network

This tool works perfectly well on the live ethereum network (I have sent a couple transactions with it).

However, as it is still in development, we urge you to excercise caution. 
For example, use `--log=3` so you can see the transaction in human readable form before you issue the same command again with `--sign --broadcast`.

If you find any issues, please report them.

# More

Use the `--help` flag to learn more about the various commands and subcommands.

If you feel something is missing, or would like to see a feature added, please open an issue, or better yet a pull request.

:)

This tool was inspired by a similar suite of tools built for working with tendermint chains at https://github.com/eris-ltd/mint-client
