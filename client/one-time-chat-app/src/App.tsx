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
  const [expires, setExpires] = useState<string>(''); // expires を文字列として保存
  const [requiresAuth, setRequiresAuth] = useState<boolean>(false);
  const [selectedTime, setSelectedTime] = useState<number>(0); // 時間
  const [selectedMinutes, setSelectedMinutes] = useState<number>(0); // 分
  const [selectedDays, setSelectedDays] = useState<number>(0); // 日数

  const getCookie = (name: string) => {
    const cookies = document.cookie.split('; ');
    const cookie = cookies.find(row => row.startsWith(`${name}=`));
    return cookie ? cookie.split('=')[1] : null;
  };

  const calculateExpiration = () => {
    const currentDate = new Date();
    if (selectedTime > 0) {
      // 現在時刻から時間を追加
      currentDate.setHours(currentDate.getHours() + selectedTime);
    }
    if (selectedMinutes > 0) {
      // 現在時刻から分を追加
      currentDate.setMinutes(currentDate.getMinutes() + selectedMinutes);
    }
    if (selectedDays > 0) {
      // 現在時刻から日数を追加
      currentDate.setDate(currentDate.getDate() + selectedDays);
    }
    return currentDate.toISOString(); // ISOフォーマットで返す
  };

  const connectToRoom = async () => {
    const baseUrl = `http://${window.location.hostname}:8080`;

    const expirationTime = calculateExpiration(); // 計算した期限を取得
    setExpires(expirationTime); // expires を更新

    const response = await fetch(`${baseUrl}/room`, {
      method: 'POST',
      body: JSON.stringify({ name: roomName, owner: clientName, expires: expirationTime, requiresAuth }),
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
    });

    if (response.ok) {
      const roomData = await response.json();
      console.log('Room creation response:', roomData);
      console.log(roomData.ID);
      const cookiesSessionID = getCookie('session_id');
      setSessionID(cookiesSessionID);
      setRoomID(roomData.ID);

      const ws = new WebSocket(`ws://${window.location.hostname}:8080/ws?room_id=${roomData.ID}&client_name=${clientName}&session_id=${cookiesSessionID}`);
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
    const baseUrl = `http://${window.location.hostname}:8080`;
    const response = await fetch(`${baseUrl}/room/${roomID}`, {
      method: 'POST',
      body: JSON.stringify({ client_name: clientName }),
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
    });
    const cookiesSessionID = getCookie('session_id');
    setSessionID(cookiesSessionID);

    if (response.ok) {
      const roomData = await response.json();
      setSessionID(roomData.sessionID);
      setRoomID(roomData.roomID);

      const ws = new WebSocket(`ws://${window.location.hostname}:8080/ws?room_id=${roomData.roomID}&client_name=${clientName}&session_id=${cookiesSessionID}`);
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

      {/* expiresとAuth設定 */}
      <div className="input-group">
        {/* 時間のプルダウン */}
        <select value={selectedTime} onChange={(e) => setSelectedTime(Number(e.target.value))}>
          <option value={0}>Select Hours</option>
          {[...Array(24)].map((_, index) => (
            <option key={index} value={index}>{index} hour{index !== 1 ? 's' : ''}</option>
          ))}
        </select>

        {/* 分のプルダウン */}
        <select value={selectedMinutes} onChange={(e) => setSelectedMinutes(Number(e.target.value))}>
          <option value={0}>Select Minutes</option>
          {[...Array(60)].map((_, index) => (
            <option key={index} value={index}>{index} minute{index !== 1 ? 's' : ''}</option>
          ))}
        </select>

        {/* 日数のプルダウン */}
        <select value={selectedDays} onChange={(e) => setSelectedDays(Number(e.target.value))}>
          <option value={0}>Select Days</option>
          {[...Array(31)].map((_, index) => (
            <option key={index} value={index}>{index} day{index !== 1 ? 's' : ''}</option>
          ))}
        </select>

        <label>
          Requires Auth
          <input
            type="checkbox"
            checked={requiresAuth}
            onChange={(e) => setRequiresAuth(e.target.checked)}
          />
        </label>
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
