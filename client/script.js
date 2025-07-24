document.addEventListener("DOMContentLoaded", () => {
  const loginView = document.getElementById("login-view");
  const gameView = document.getElementById("game-view");
  const loginForm = document.getElementById("login-form");
  const playerNameInput = document.getElementById("player-name-input");
  const statusDiv = document.getElementById("status");

  const canvas = document.getElementById("game-canvas");
  const ctx = canvas.getContext("2d");

  const chatForm = document.getElementById("chat-form");
  const chatInput = document.getElementById("chat-input");
  const chatLog = document.getElementById("chat-log");

  const TILE_SIZE = 20; // Kích thước của mỗi ô trong game để vẽ
  const CHAT_BUBBLE_DURATION = 5000; // 5 giây (giống server)

  let ws = null; // Đối tượng WebSocket
  let playerID = null; // ID của client này, được server cấp
  let gameState = { players: [], chat: [] }; // Trạng thái game mới nhất từ server
  let lastChatMessages = {}; // Lưu trữ tin nhắn cuối cùng của mỗi người chơi để hiển thị bong bóng chat

  /**
   * Gửi một đối tượng message đến server dưới dạng JSON.
   * @param {object} message - Đối tượng message cần gửi.
   */
  function sendMessage(message) {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(message));
    } else {
      console.error("WebSocket is not connected.");
    }
  }

  // Xử lý sự kiện từ Server
  function setupWebSocketListeners() {
    ws.onopen = () => {
      console.log("Connected to server.");
      statusDiv.textContent = "Connected. Joining game...";

      // Gửi tin nhắn "init" để tham gia game
      const initMessage = {
        type: "init",
        player_name: playerNameInput.value.trim(),
      };
      sendMessage(initMessage);
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);

      // Xử lý message từ server dựa trên 'type'
      switch (data.type) {
        case "init":
          // Server đã chấp nhận, lưu lại Player ID
          playerID = data.player_id;
          console.log("Game joined. My ID:", playerID);

          // Chuyển đổi giao diện
          loginView.style.display = "none";
          gameView.style.display = "flex";
          gameView.style.flexDirection = "column";

          // Bắt đầu vòng lặp game
          requestAnimationFrame(gameLoop);
          break;

        case "game_state":
          // Cập nhật trạng thái game và vẽ lại
          gameState = data.game_state;

          // Cập nhật bong bóng chat
          updateChatBubbles();
          // Cập nhật khung chat log
          updateChatLog();
          break;

        case "error":
          // Hiển thị lỗi từ server
          console.error("Server error:", data.message);
          statusDiv.textContent = `Error: ${data.message}`;
          alert(`Error: ${data.message}`);
          ws.close(); // Đóng kết nối khi có lỗi nghiêm trọng
          break;
      }
    };

    ws.onclose = () => {
      console.log("Disconnected from server.");
      statusDiv.textContent = "Disconnected. Please refresh to rejoin.";
      // Reset giao diện về màn hình đăng nhập
      loginView.style.display = "block";
      gameView.style.display = "none";
      ws = null;
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
      statusDiv.textContent = "Connection error.";
    };
  }

  // Xử lý Input từ Người dùng

  // 1. Form Đăng nhập
  loginForm.addEventListener("submit", (e) => {
    e.preventDefault();
    const playerName = playerNameInput.value.trim();
    if (playerName && !ws) {
      statusDiv.textContent = "Connecting...";
      // Kết nối đến server, thay đổi địa chỉ nếu cần
      ws = new WebSocket("ws://localhost:8080/ws");
      setupWebSocketListeners();
    }
  });

  // 2. Form Gửi Chat
  chatForm.addEventListener("submit", (e) => {
    e.preventDefault();
    const messageText = chatInput.value.trim();
    if (messageText) {
      const chatMessage = {
        type: "chat",
        message: messageText,
      };
      sendMessage(chatMessage);
      chatInput.value = ""; // Xóa input sau khi gửi
    }
  });

  // 3. Di chuyển bằng bàn phím
  document.addEventListener("keydown", (e) => {
    if (!playerID) return; // Chưa vào game thì không di chuyển

    let direction = null;
    switch (e.key) {
      case "ArrowUp":
      case "w":
      case "W":
        direction = "up";
        break;
      case "ArrowDown":
      case "s":
      case "S":
        direction = "down";
        break;
      case "ArrowLeft":
      case "a":
      case "A":
        direction = "left";
        break;
      case "ArrowRight":
      case "d":
      case "D":
        direction = "right";
        break;
    }

    if (direction) {
      e.preventDefault(); // Ngăn cuộn trang
      const moveMessage = {
        type: "move",
        direction: direction,
      };
      sendMessage(moveMessage);
    }
  });

  //  Logic Vẽ (Rendering)

  /**
   * Vòng lặp chính của game, được gọi liên tục để vẽ lại màn hình.
   */
  function gameLoop() {
    if (!ws) return; // Dừng vòng lặp nếu mất kết nối

    // Xóa toàn bộ canvas
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // Vẽ người chơi
    if (gameState.players) {
      gameState.players.forEach(drawPlayer);
    }

    requestAnimationFrame(gameLoop);
  }

  /**
   * Vẽ một người chơi lên canvas.
   * @param {object} player - Đối tượng người chơi từ gameState.
   */
  function drawPlayer(player) {
    const x = player.position.x * TILE_SIZE;
    const y = player.position.y * TILE_SIZE;

    // Vẽ thân người chơi (hình vuông)
    ctx.fillStyle = player.id === playerID ? "blue" : "red"; // Phân biệt bản thân và người khác
    ctx.fillRect(x, y, TILE_SIZE, TILE_SIZE);

    // Vẽ tên người chơi
    ctx.fillStyle = "black";
    ctx.textAlign = "center";
    ctx.fillText(player.name, x + TILE_SIZE / 2, y - 5);

    // Vẽ bong bóng chat nếu có
    const chatInfo = lastChatMessages[player.id];
    if (chatInfo && Date.now() - chatInfo.time < CHAT_BUBBLE_DURATION) {
      drawChatBubble(ctx, chatInfo.message, x + TILE_SIZE / 2, y - 20);
    }
  }

  /**
   * Vẽ bong bóng chat phía trên người chơi.
   */
  function drawChatBubble(ctx, text, x, y) {
    ctx.font = "12px sans-serif";
    const textWidth = ctx.measureText(text).width;

    const bubblePadding = 5;
    const bubbleWidth = textWidth + bubblePadding * 2;
    const bubbleHeight = 20;

    // Vẽ hộp thoại
    ctx.fillStyle = "white";
    ctx.globalAlpha = 0.8;
    ctx.fillRect(
      x - bubbleWidth / 2,
      y - bubbleHeight,
      bubbleWidth,
      bubbleHeight
    );
    ctx.globalAlpha = 1.0;

    // Vẽ viền
    ctx.strokeStyle = "black";
    ctx.strokeRect(
      x - bubbleWidth / 2,
      y - bubbleHeight,
      bubbleWidth,
      bubbleHeight
    );

    // Vẽ text
    ctx.fillStyle = "black";
    ctx.textAlign = "center";
    ctx.fillText(text, x, y - bubbleHeight / 2 + 5);
  }

  // Cập nhật Giao diện Chat

  /**
   * Cập nhật danh sách tin nhắn trong chat log.
   */
  function updateChatLog() {
    chatLog.innerHTML = ""; // Xóa log cũ
    const playersMap = new Map(gameState.players.map((p) => [p.id, p.name]));

    gameState.chat.forEach((chatMsg) => {
      const playerName = playersMap.get(chatMsg.user_id) || "Unknown";
      const p = document.createElement("p");
      p.innerHTML = `<strong>${playerName}:</strong> ${chatMsg.message}`;
      chatLog.appendChild(p);
    });
    // Tự động cuộn xuống dưới
    chatLog.scrollTop = chatLog.scrollHeight;
  }

  /**
   * Kiểm tra và cập nhật tin nhắn mới nhất để hiển thị bong bóng chat.
   */
  function updateChatBubbles() {
    if (!gameState.chat || gameState.chat.length === 0) return;

    const latestMessage = gameState.chat[gameState.chat.length - 1];

    // Chỉ cập nhật nếu tin nhắn chưa được xử lý hoặc là tin mới hơn
    const existingMessage = lastChatMessages[latestMessage.user_id];
    if (
      !existingMessage ||
      new Date(latestMessage.time_sent) > new Date(existingMessage.sentAt)
    ) {
      lastChatMessages[latestMessage.user_id] = {
        message: latestMessage.message,
        time: Date.now(), // Lưu thời gian client nhận được
        sentAt: new Date(latestMessage.time_sent), // Lưu thời gian server gửi
      };
    }
  }
});
