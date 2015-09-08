# eth-client
Command line interfaces for low-level conversations with ethereum chains

# Overview

In current form, the `eth-client` has three tools: 

```
ethgen - for creating genesis.json files
ethinfo - for querying an ethereum chain for information
ethtx - for crafting and broadcasting transactions
```

# Dependencies

eth-client depends on `eris-keys` for signing transactions. It's also useful for generating keys.

I will presume you have [go installed](https://golang.org/doc/install),
and that you have set your `$GOPATH` and put `$GOPATH/bin` on your `$PATH`. 

Grab and start the keys server with

```bash
go get github.com/eris-ltd/eris-keys
eris-keys server &
```

# Ethereum Test Chains

Let's boot up a test chain and send some transactions to it. First we need to generate a key:

```bash
ADDR=`eris-keys gen --type=secp256k1,sha3`
```

This will create a new ethereum key for you and return the address (`echo $ADDR`). 

Now let's create a directory for our new chain and make a genesis.json file in it using our address:

```bash
mkdir ~/.myethereum
ethgen $ADDR > ~/.myethereum/genesis.json
```

Splendid! Take a look at that genesis file if you like. 
Note that `ethgen` chooses reasonable default values for the difficulty and the gas-limit for running test chains. 
If you want to customize, just use the flags. See `ethgen --help` for more.

Now let's boot the chain (we presume you already have [geth](https://github.com/ethereum/go-ethereum/) installed):

```bash
geth --datadir ~/.myethereum --rpc --mine --genesis ~/.myethereum/genesis.json --maxpeers 0 --etherbase $ADDR --verbosity 7
```

This will start a private ethereum chain on your machine. 
Note the unfortunate reality of ethereum's proof of work means you need at least 1GB of free RAM for mining to work. 
It may take some time for mining to get started (on the order of a few minutes). 
Note on my 2013 macbook with 4GB RAM I cannot mine a test chain and have Firefox open at the same time.

# Ethereum Transactions

Now that you have a chain setup and are mining it, let's send a transaction!

Of course, since `geth` is now running, we should use a new window. Make sure to copy the ADDR variable and set it again in the new window:

```
export ADDR=<address we used>
```

First, we'll generate another account to send funds to:

```bash
ADDR2=`eris-keys gen --type=secp256k1,sha3`
```

Now let's craft a transaction:

```bash
ethtx send --addr=$ADDR --to=$ADDR2 --amt=10 --gas=21000 --price=100000000000 --sign --broadcast
```

The command should execute successfuly, returning the transaction ID.

We can now check the account:

```bash
ethinfo account $ADDR2
```

The balance should be `0xa` (ie. `10`)!

Ok, let's break down the `ethtx` command a little bit. Ethereum only has one official transaction type, but it serves three distinct purposes. 
You can simply send funds from one account to another, or you can create a contract, or you can call a contract.
So `ethtx` has three main commands: `send`, `call`, and `create`.
Over time, `ethtx` will incorporate commands for talking to the major dapps, to facilitate interactions with them. The first such dapp will be the name reg, and you'll be able to use `ethtx name` to register a new name there.

I should note, the `ethtx` flags accept both hex and base 10 numbers. If you are using hex, make sure to prefix with `0x`.

You can also specify a transaction's nonce with the `--nonce` flag. If no nonce is specified, the correct one is fetched from the blockchain.

The `--sign` and `--broadcast` flags allow you to specify exactly what you want to do. Maybe you only want to craft the bytes for transaction and sign it later? Maybe you only want to sign it and broadcast it later? Or maybe you want to do everything now, in which case both `--sign` and `--broadcast` are appropriate. Soon, we will add a `--wait` feature so you can wait until the transaction is actually committed in a block.

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

If you want to compile solidity, check out the lovely-little-languages compiler server at `github.com/eris-ltd/lllc-server`. 

We will also we wrapping the eth-client tools at a slightly higher level so that solidity compiling and abi formatting may be done for you.

The purpose of the `eth-client` itself however is a simple, low-level interface for developers.

Enjoy!

# More

Use the `--help` flag to learn more about the various commands and subcommands.

If you feel something is missing, or would like to see a feature added, please open an issue, or better yet a pull request.

:)
