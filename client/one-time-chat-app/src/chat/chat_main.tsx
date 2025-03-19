import React, { use, useEffect, useState } from 'react';
import { connect } from 'socket.io-client';

const Chat: React.FC = () => {
    type Clients = {
        Name: string;
        SessionID: string;
        Ws: any;
    };
    type RoomInfo = {
        ID: string;
        name: string;
        owner: string;
        ownerSessionID: string;
        expires: string;
        requiresAuth: boolean;
        UnauthenticatedClients: Clients[];
        AuthenticatedClients: Clients[];
        Mu: Record<string, unknown>;
    };

    const [roomInfo, setRoomInfo] = useState<RoomInfo | null>(null);
    const [ws, setWs] = useState<any>(null);
    const [message, setMessage] = useState<{ sender: string; content: string }[]>([]);

    const getCookie = (name: string) => {
        const cookies = document.cookie.split('; ');
        const cookie = cookies.find(row => row.startsWith(`${name}=`));
        return cookie ? cookie.split('=')[1] : null;
    };

    const connectToRoom = async (roomID: string) => {
        const APIURL = process.env.REACT_APP_WSAPI_URL;
        const clientName = getCookie('user_name');
        const cookiesSessionID = getCookie('session_id');
            const URL = `${APIURL}/ws?room_id=${roomID}&client_name=${clientName}&session_id=${cookiesSessionID}`;
            const ws = new WebSocket(URL);
            ws.onopen = () => {
                console.log('WebSocket connected');
            }
            ws.onmessage = (event) => {
                const data = JSON.parse(event.data);
                setMessage((prevMessages) => [...prevMessages, { sender: data.sender, content: data.sentence }]);
              };
            ws.onclose = () => {
                console.log('WebSocket closed');
            }
            ws.onerror = (error) => {
                console.error('WebSocket error:', error);
            }
            setWs(ws);

        };
        const sendMessage = (message: string) => {
            if (ws && ws.readyState === WebSocket.OPEN) {
                ws.send(message);
            } else {
                console.error('WebSocketはまだ開いていません。現在の状態:', ws?.readyState);
            }
        };
        
        



    const fetchRoomInfo = async () => {
        const isOwner = getCookie('is_owner') === 'true';
        const roomID = getCookie('room_id');
        const APIURL = process.env.REACT_APP_API_URL;
        const URL = `${APIURL}/room/${roomID}`;
        if (isOwner) {
            const response = await fetch(URL, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
            });

            const roomData = await response.json();
            setRoomInfo(roomData);

        } else {
            console.log('User is not the owner');
        }

        if (roomID) {
            connectToRoom(roomID);
        }else{
            console.log('Room ID not found');
        }
    };
    useEffect(() => {
        fetchRoomInfo();
    }, []); 

   
        
    if (!roomInfo) {
        return <div>Loading...</div>;
    }

    return (
        <div>
            <h1>Chat Page</h1>
            <p>Room ID: {roomInfo.ID}</p>
            <p>Room Name: {roomInfo.name}</p>
            <p>Owner: {roomInfo.owner}</p>
            <p>Expires: {roomInfo.expires}</p>
            {/* You can add more fields from roomInfo to display here */}
            <button 
                onClick={() => {
                    sendMessage('Hello, World!');
                }}
            >
                Send Message
            </button>
        </div>
    );
};

export default Chat;
