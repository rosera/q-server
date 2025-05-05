# Quizzrr Server

A websocket server to enable simple multi-player functionality.

## Backend API

Supported functionality:


- [ ] Create a Room

- [ ] Join a Room
- [ ] Leave a Room
- [ ] List Room

- [ ] Start Game
- [ ] Reset Game
- [ ] End Game

- [ ] Next Game Question

## Architecture

The backend server perform admin and user functions using websockets.

![q-server](https://github.com/rosera/q-server/blob/main/screenshots/q-server-api.png "q-server")


## Running the q-server

1. Initiate the quizzrr websocket server
   ```bash
   ./q-server
   ```

2. Open a new terminal to accept `websocket` commands.

3. To interact with the application use the `websocat` tool to send message to the api.

   Example: To run the `websocat` server on `localhost`

   ```bash
   websocat ws://localhost:8080/ws
   ```

   __NOTE:__
   Use `nix` to run websocat
   ```bash
   nix-shell -p websocat
   ```

With the websocket available, send the required command to the backend api.
Available commands are detailed below:

## Create Room

* Role: Admin
* Description: Create a game room defined by the room_id.

| Role | JSON |
|------|------|
| Admin | {"type": "create_room", "room_id": "room123"} |

1. Admins can create a room on the server
   ```json
   {"type": "create_room", "room_id": "room123"}
   ```
2. Admins can create a room on the server
   ```json
   {"type": "create_room", "room_id": "room113"}
   
3. Admins can create a room on the server
   ```json
   {"type": "create_room", "room_id": "room103"}
   ```
   
Players can join the room by using the `room_id`

## Join Room

Users can join a game room defined by entering the room_id.

| Role | JSON |
|------|------|
| User | {"type": "join_room", "room_id": "room123", "name": "Alice"} |

1. Add user: `Alice`
   ```json
   {"type": "join_room", "room_id": "room123", "name": "Alice"}
   ```

2. Add user: `Bob`
   ```json
   {"type": "join_room", "room_id": "room123", "name": "Bob"}
   ```

3. Add user: `Carol`
   ```json
   {"type": "join_room", "room_id": "room123", "name": "Carol"}
   ```

3. Add user: `Danny`
   ```json
   {"type": "join_room", "room_id": "room103", "name": "Danny"}
   ```

## Start Game 

Admin can start a game room defined by the room_id.

| Role | JSON |
|------|------|
| Admin | {"type": "start_room", "room_id": "room123"} |

1. Admin start a game in `room_id`
   ```json
   {"type": "start_room", "room_id": "room123"}
   ```


## Reset Game 

Admin can reset a game defined by the room_id.

| Role | JSON |
|------|------|
| Admin | {"type": "reset_game", "room_id": "room123"} |

1. End the game defined by `room_id`
   ```json
   {"type": "reset_game", "room_id": "room123"}
   ```


## End Game 

Admin can end a game defined by the room_id.

| Role | JSON |
|------|------|
| Admin | {"type": "end_game", "room_id": "room123"} |

1. End the game defined by `room_id`
   ```json
   {"type": "end_room", "room_id": "room123"}
   ```

## List Room 

Users can list room and metadata.

| Role | JSON |
|------|------|
| Users | {"type": "list_rooms"} |

1. List rooms
   ```json
   {"type": "list_rooms"}
   ```

## Next Question

Admin can indicate to move to the next question.

| Role | JSON |
|------|------|
| Admin | {"type": "next_question", "room_id": "room123"} |

1. Admin event to move to the next question
   ```json
   {"type": "next_question", "room_id": "room123"}
   ```


