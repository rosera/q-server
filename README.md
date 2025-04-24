# Quizzrr Server

A simple websocket server to enable multi-player functionality.

## Backend API

Supported functionality:

- [ ] Create a Room
- [ ] Join a Room
- [ ] Start Game
- [ ] End Game
- [ ] Next Game Question


## Create Room

* Role: Admin
* Description: Create a game room defined by the room_id.

| Role | JSON |
|------|------|
| Admin | {"type": "create_room", "room_id": "room123"} |

## Join Room

Join a game room defined by the room_id.

| Role | JSON |
|------|------|
| User | {"type": "join_room", "room_id": "room123", "name": "Alice"} |

## Start Game 

Start a game room defined by the room_id.

| Role | JSON |
|------|------|
| Admin | {"type": "start_room", "room_id": "room123"} |

## End Game 

End a game room defined by the room_id.

| Role | JSON |
|------|------|
| Admin | {"type": "end_game", "room_id": "room123"} |


## Next Question

Move to the next question.

| Role | JSON |
|------|------|
| Admin | {"type": "next_question", "room_id": "room123"} |


