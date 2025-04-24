# Quizzrr Server

A websocket server to enable simple multi-player functionality.

## Backend API

Supported functionality:

- [ ] Create a Room
- [ ] Join a Room
- [ ] Start Game
- [ ] End Game
- [ ] Next Game Question

## Architecture

The backend server perform admin and user functions using websockets.

![q-server](https://github.com/rosera/q-server/blob/main/screenshots/q-server-api.png "q-server")


## Running the q-server

```bash
./q-server

```

To interact with the application use the `websocat` tool to send message to the api.

Example: To run the `websocat` server on `localhost`

```bash
websocat ws://localhost:8080/ws
```

With the websocket available, send the required command to the backend api.
Available commands are detailed below:

## Create Room

* Role: Admin
* Description: Create a game room defined by the room_id.

| Role | JSON |
|------|------|
| Admin | {"type": "create_room", "room_id": "room123"} |


```json
{"type": "create_room", "room_id": "room123"}
```

## Join Room

Join a game room defined by the room_id.

| Role | JSON |
|------|------|
| User | {"type": "join_room", "room_id": "room123", "name": "Alice"} |

```json
{"type": "join_room", "room_id": "room123", "name": "Alice"} 
```

## Start Game 

Start a game room defined by the room_id.

| Role | JSON |
|------|------|
| Admin | {"type": "start_room", "room_id": "room123"} |

```json
{"type": "start_room", "room_id": "room123"}
```

## End Game 

End a game room defined by the room_id.

| Role | JSON |
|------|------|
| Admin | {"type": "end_game", "room_id": "room123"} |

```json
{"type": "end_room", "room_id": "room123"}
```

## Next Question

Move to the next question.

| Role | JSON |
|------|------|
| Admin | {"type": "next_question", "room_id": "room123"} |


```json
{"type": "end_room", "room_id": "room123"}
```
