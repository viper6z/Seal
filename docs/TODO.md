* [ ]: Create `udp-service/` with `README.md`, `server.py`, `Dockerfile` and `compose.yaml`

* [ ]: Add `udp-service` to the root Compose file and connect it only to `backend`

* [ ]: Run `docker compose config` and confirm the service has no host port mapping

* [ ]: Build the first UDP server that binds to `0.0.0.0:9001`

* [ ]: Use `recvfrom()` to receive both the datagram and sender address

* [ ]: Handle `JOIN` and store subscribed client addresses

* [ ]: Handle `UPDATE <text>`

* [ ]: Add a server-generated sequence number for every accepted update

* [ ]: Broadcast `TEXT <sequence> <text>` to all subscribed clients with `sendto()`

* [ ]: Create `toolbox/udp_client.py`

* [ ]: Make the client use one UDP socket for both sending and receiving

* [ ]: Send `JOIN` when the client starts

* [ ]: First make a simple version where `input()` sends `UPDATE <text>` after Enter

* [ ]: Test two clients in separate SSH terminals and make sure both receive the same broadcast

* [ ]: Add client-side sequence checking so older `TEXT` packets are ignored

* [ ]: Upgrade the client from `input()` to raw keypress reading

* [ ]: Keep a local `current_text` buffer and support normal characters, Backspace and Ctrl-C

* [ ]: Send the latest full text every 0.1 seconds only when it changed

* [ ]: Make terminal redraws readable when incoming broadcasts arrive

* [ ]: Test server restart, wrong port, non-subscribed client and multiple clients typing

* [ ]: Add a short logbook entry, commit and push






















































































