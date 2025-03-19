import React, { useEffect, useState } from 'react';
import { CircularProgress, Button, Typography, Box, List, ListItem, ListItemText } from '@mui/material';
import './css/chat_main.css';
import ChatArea from './chat_area.tsx';


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
    const [isConnectedWS, setIsConnectedWS] = useState<boolean>(false);
    const [message, setMessage] = useState<{ sender: string; content: string; isMe:boolean }[]>([]);
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
    const [isLoading, setIsLoading] = useState<boolean>(true); // 初期値を true に設定

    const getCookie = (name: string) => {
        const cookies = document.cookie.split('; ');
        const cookie = cookies.find(row => row.startsWith(`${name}=`));
        return cookie ? decodeURIComponent(cookie.split('=')[1]) : null; // URLデコードを追加
    };
    

    const connectToRoom = async (roomID: string) => {
        const APIURL = process.env.REACT_APP_WSAPI_URL;
        const clientName = getCookie('user_name');
        const cookiesSessionID = getCookie('session_id');
        const URL = `${APIURL}/ws?room_id=${roomID}&client_name=${clientName}&session_id=${cookiesSessionID}`;
        const ws = new WebSocket(URL);
        ws.onopen = () => {
            setIsConnectedWS(true);
            console.log('WebSocket connected');
        }
        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            console.log('Received:', data);
            setMessage((prevMessages) => [...prevMessages, { sender: data.sender, content: data.sentence, isMe: false }]);
        };
        ws.onclose = () => {
            setIsConnectedWS(false);
            console.log('WebSocket closed');
        };
        ws.onerror = (error) => {
            setIsConnectedWS(false);
            console.error('WebSocket error:', error);
        };
        setWs(ws);
    };

    const sendMessage = (message: string) => {
        let clientName = getCookie('user_name');
        if (clientName === null) {
            clientName = 'Unknown';
        }
        if (ws && ws.readyState === WebSocket.OPEN) {
            setMessage((prevMessages) => [...prevMessages, { sender:clientName, content: message, isMe: true }]);
            ws.send(message);
        } else {
            console.error('WebSocketはまだ開いていません。現在の状態:', ws?.readyState);
        }
    };

    const fetchRoomInfo = async () => {
        const isOwner = getCookie('is_owner') === 'true';
        console.log('isOwner:', isOwner);
        const roomID = getCookie('room_id');
        const APIURL = process.env.REACT_APP_API_URL;
        if (isOwner) {
            const URL = `${APIURL}/room/${roomID}`;
            const response = await fetch(URL, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
            });

            const roomData = await response.json();
            if (roomData.ownerSessionID !== getCookie('session_id')) {
                console.log('Owner session ID does not match');
                alert("あなたはオーナーではありません。cookieは変更しないでください。");
                return;
            }
            setIsAuthenticated(true);
            setRoomInfo(roomData);
            if (roomID){
                connectToRoom(roomID);
            }
        } else if (roomID) {
            const sessionID = getCookie('session_id');
            const URL = `${APIURL}/room/${roomID}/isAuth?client_session_id=${sessionID}`;
            const response = await fetch(URL, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
            });
            const roomData = await response.json();
            if (roomData.isAuth) {
                setIsAuthenticated(true);
                connectToRoom(roomID);
            } else {
                setIsAuthenticated(false);
                console.log('Not Authenticated');
            }
        } else {
            console.log('Room ID not found');
        }
        setIsLoading(false); // ローディングを終了
    };

    useEffect(() => {
        fetchRoomInfo();
    }, []);

    if (isLoading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
                <CircularProgress />
            </Box>
        );
    }

    if (!isAuthenticated) {
        return (
            <Box sx={{ textAlign: 'center' }}>
                <Typography variant="h4">認証されていません</Typography>
                <Button variant="contained" onClick={fetchRoomInfo} sx={{ marginTop: 2 }}>
                    再試行
                </Button>
            </Box>
        );
    }

    if (!isConnectedWS) {
        return (
            <Box sx={{ textAlign: 'center' }}>
                <CircularProgress />
            </Box>
        );
    }


    return (
        <div className="chat-container">
            <div className="members-section">
                <div className="members-header">メンバー (5人)</div>
                <div className="member-item">Aさん (オーナー)</div>
                <div className="member-item">0Aさん <span className="kick-button">キック</span></div>
                <div className="member-item">1Aさん <span className="kick-button">キック</span></div>
                <div className="member-item">0Aさん <span className="kick-button">キック</span></div>
                <div className="member-item">0Aさん <span className="kick-button">キック</span></div>
            </div>
            <ChatArea message={message} sendMessage={sendMessage} />
            <div className="requests-section">
                <div className="requests-header">リクエスト</div>
                <div className="request-item">Aさん <span className="approve-button">承認</span></div>
                <div className="request-item">Aさん <span className="approve-button">承認</span></div>
                <div className="request-item">Aさん <span className="approve-button">承認</span></div>
                <div className="request-item">Aさん <span className="approve-button">承認</span></div>
                <div className="request-item">Aさん <span className="approve-button">承認</span></div>
            </div>
        </div>
    );
};

export default Chat;
