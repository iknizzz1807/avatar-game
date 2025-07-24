# Go Avatar Game

A simple 2D multiplayer game project built with Go and WebSockets, allowing players to move their avatars and chat with each other in real-time.

## Features

- **Real-time Communication:** Utilizes WebSockets for bidirectional communication between the client and server.
- **Player Movement:** Players can move their avatars around the map using the arrow keys.
- **Chat System:**
  - **Chat Log:** A shared chat box displays the history of the conversation.
  - **Chat Bubbles:** A player's latest message appears above their avatar's head for a short duration.
- **Multiplayer:** Supports multiple players connecting and interacting within the same game world.

## Technology Stack

- **Backend:**
  - **Language:** Go
  - **WebSocket:** `gorilla/websocket` library
- **Frontend:**
  - **Languages:** HTML, CSS, JavaScript (Vanilla)
  - **Graphics:** HTML5 Canvas
  - **Communication:** Browser's native WebSocket API

## Setup and Running

1.  **Run the Server (Backend):**

    - Open a terminal and navigate to the `server` directory.
    - Run the following command to start the server:
      ```bash
      go run main.go
      ```
    - The server will listen for connections at `ws://localhost:8080/ws`.

2.  **Run the Client (Frontend):**
    - Serve `frontend` folder (recommend).
    - Open the `client/index.html` file in your web browser (you can open multiple tabs to simulate multiple players). (not recommends).
    - Enter a player name and click "Join".
    - Use the arrow keys to move.
