# TCP Server

This is a TCP server implemented in Go that handles requests for generating and validating proofs of work (PoW).

### Protocol description

App use custom protocol for swift and efficient interactions between a client and a server and and focused on transferring as little data as possible

Structures of requests/responses:
1. **Request to ask a puzzle from a server:**

```[1]```
2. **Response for a request from point 1:**

```
[2, puzzleSize, targetBits, puzzle[1], ..., puzzle[puzzleSize]]
```
3. **After the client processed a response from point 2, it should return response:**

```
[3, nonce[1], ..., nonce[4], hash[1], ..., hash[32]]
```

4. **Response of result of validation of PoW:**
- In case of success of validation:
  ```
  [4, sizeQuote, quote[1], ..., quote[sizeQuote]]
  ```
- In case of failure of validation:
  ```
  [5]
  ```


### Usage

Using docker-compose:

`docker compose -f docker-compose.yml -p pow up`

### Config

Config is customizable through environmental variables. You can see the default values for the config [here](pkg/config/config.go).

Variables and descriptions:
- `HOST`: Host for the server
- `PORT`: Port for the server
- `PUZZLE_SIZE`: Size of puzzle returned from the server to the client
- `TARGET_BITS`: Number of bits for Proof of Work (PoW), affecting the complexity of calculations
- `CONN_TIMEOUT`: Timeout for the server for an authorization session and timeout for a client for an authorization session


### What to improve

- Add sending to a client what precisely went wrong.
- Add the possibility to send PoW from another connection.
- Implement the server and client as interfaces.
- Increase coverage of the project (probably after the previous point, it will be easier).
- More customization through config / consider moving config to YAML/.env.
