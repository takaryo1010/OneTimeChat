import React, { useState } from 'react';
import './App.css';
import { Socket } from 'socket.io-client';

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
  const [roomID, setRoomID] = useState<string>('');

  const getCookie = (name: string) => {
    const cookies = document.cookie.split('; ');
    const cookie = cookies.find(row => row.startsWith(`${name}=`));
    return cookie ? cookie.split('=')[1] : null;
  };

  const connectToRoom = async () => {
    const response = await fetch('http://localhost:8080/room', {
      method: 'POST',
      body: JSON.stringify({ name: roomName, owner: clientName }),
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
    });

    if (response.ok) {
      const roomData = await response.json();
      console.log('Room creation response:', roomData);
      console.log(roomData.ID);
      const cookiesSessionID = getCookie('session_id');
      setSessionID(cookiesSessionID); // セッションIDを保存
      setRoomID(roomData.ID); // ルームIDを保存

      const ws = new WebSocket(`ws://localhost:8080/ws?room_id=${roomData.ID}&client_name=${clientName}&session_id=${sessionID}`);
      ws.onopen = () => {
        console.log('WebSocket connected');
      };

      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        setMessages((prevMessages) => [...prevMessages, { sender: data.sender, content: data.sentence }]);
      };

      ws.onclose = () => {
        console.log('WebSocket disconnected');
      };

      setSocket(ws);
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
      const roomData = await response.json();
      setSessionID(roomData.sessionID); // セッションIDを保存
      setRoomID(roomData.roomID); // ルームIDを保存

      const ws = new WebSocket(`ws://localhost:8080/ws?room_id=${roomID}&client_name=${clientName}&session_id=${sessionID}`);
      ws.onopen = () => {
        console.log('WebSocket connected');
      };

      ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        setMessages((prevMessages) => [...prevMessages, { sender: data.sender, content: data.sentence }]);
      };

      ws.onclose = () => {
        console.log('WebSocket disconnected');
      };

      setSocket(ws);
    } else {
      alert('Failed to join the room');
    }
  };

  const sendMessage = (message: string) => {
    if (socket) {
      socket.send(message);
    }
  };

  return (
      <div className="App">
    <h1>OneTime Chat</h1>

    {/* ルーム名とクライアント名の入力 */}
    <div className="input-group">
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
    </div>

    {/* 既存のルームIDを入力して参加 */}
    <div className="input-group">
      <input
        type="text"
        placeholder="Enter Room ID to Join"
        value={roomID}
        onChange={(e) => setRoomID(e.target.value)}
      />
      <button onClick={joinRoom}>Join Room</button>
    </div>

    {/* セッションID表示 */}
    <div className="room-info">
      {sessionID && <p>ROOM ID: {roomID}</p>}
    </div>

    {/* メッセージ表示 */}
    <div className="messages">
      {messages.map((msg, index) => (
        <p key={index}><strong>{msg.sender}:</strong> {msg.content}</p>
      ))}
    </div>

    {/* メッセージ送信フォーム */}
    <div className="message-form">
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
  </div>

  );
};

export default App;
