import React, { useState } from 'react';
import './App.css';
import { io, Socket } from 'socket.io-client';

interface Message {
  sender: string;
  content: string;
}

const App: React.FC = () => {
  const [roomName, setRoomName] = useState<string>('');
  const [clientName, setClientName] = useState<string>('');
  const [messages, setMessages] = useState<Message[]>([]);
  const [socket, setSocket] = useState<Socket | null>(null);
  const [sessionID, setSessionID] = useState<string>('');
  const [roomID, setRoomID] = useState<string>('');  // 追加: ルームID用の状態

  const connectToRoom = async () => {
    const response = await fetch('http://localhost:8080/room', {
      method: 'POST',
      body: JSON.stringify({ name: roomName, owner: clientName }),
      headers: { 'Content-Type': 'application/json' },
    });
  
    if (response.ok) {
      const roomData = await response.json();
      setSessionID(roomData.sessionID); // セッションIDを保存
      setRoomID(roomData.id); // ルームIDを保存
      // WebSocket接続の開始
      const socketInstance = io('http://localhost:8080', {
        query: { room_id: roomID, client_name: clientName, session_id: roomData.sessionID },
        transports: ['websocket'], // WebSocketで接続
      });
  
      socketInstance.on('message', (message: string) => {
        setMessages((prevMessages) => [...prevMessages, { sender: 'Server', content: message }]);
      });
  
      setSocket(socketInstance);
    } else {
      alert('Room creation failed');
    }
  };
  
  const joinRoom = async () => {
    const response = await fetch(`http://localhost:8080/room/${roomID}`, {
      method: 'POST',
      body: JSON.stringify({ client_name: clientName }),
      headers: { 'Content-Type': 'application/json' },
    });
  
    if (response.ok) {
      // ここでセッションIDを設定する必要がある場合
      console.log('Join room response:', response);
      const roomData = await response.json();
      console.log(roomData);
      setSessionID(roomData.sessionID); // セッションIDを保存
      setRoomID(roomData.roomID); // ルームIDを保存
  
      const socketInstance = io('http://localhost:8080', {
        query: { room_id: roomID, client_name: clientName, session_id: roomData.sessionID },
        transports: ['websocket'],
      });
  
      socketInstance.on('message', (message: string) => {
        setMessages((prevMessages) => [...prevMessages, { sender: 'Server', content: message }]);
      });
  
      setSocket(socketInstance);
    } else {
      alert('Failed to join the room');
    }
  };
  
  

  // メッセージ送信
  const sendMessage = (message: string) => {
    if (socket) {
      socket.emit('message', message);
      setMessages((prevMessages) => [...prevMessages, { sender: clientName, content: message }]);
    }
  };

  return (
    <div className="App">
      <h1>OneTime Chat</h1>

      {/* ルーム名とクライアント名の入力 */}
      <input
        type="text"
        placeholder="Enter Room Name"
        value={roomName}
        onChange={(e) => setRoomName(e.target.value)}
      />
      <input
        type="text"
        placeholder="Enter Your Name"
        value={clientName}
        onChange={(e) => setClientName(e.target.value)}
      />
      <button onClick={connectToRoom}>Create Room</button>

      {/* 既存のルームIDを入力して参加 */}
      <input
        type="text"
        placeholder="Enter Room ID to Join"
        value={roomID}
        onChange={(e) => setRoomID(e.target.value)}
      />
      <button onClick={joinRoom}>Join Room</button>

      {/* セッションID表示 */}
      {sessionID && <p>ROOM ID: {roomID}</p>}

      {/* メッセージ表示 */}
      <div>
        {messages.map((msg, index) => (
          <p key={index}><strong>{msg.sender}:</strong> {msg.content}</p>
        ))}
      </div>

      {/* メッセージ送信フォーム */}
      <input
        type="text"
        placeholder="Type a message"
        onKeyDown={(e) => {
          if (e.key === 'Enter') {
            sendMessage(e.currentTarget.value);
            e.currentTarget.value = '';
          }
        }}
      />
    </div>
  );
};

export default App;
