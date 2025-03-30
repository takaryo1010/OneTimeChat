import React, { useEffect, useState } from 'react';
import { CircularProgress, Button, Typography, Box, IconButton } from '@mui/material';
import { Refresh, RemoveCircle } from '@mui/icons-material';  // RemoveCircle アイコンを使用
import './css/chat_main.css';
import ChatArea from './chat_area.tsx';

const Chat: React.FC = () => {
    type Clients = {
        name: string;
        clientid: string;
        isowner?: boolean;
    };
    type RoomInfo = {
        ID: string;
        name: string;
        owner: string;
        expires: string;
        requiresAuth: boolean;
        UnauthenticatedClients: Clients[];
        AuthenticatedClients: Clients[];
        Mu: Record<string, unknown>;
    };

    const [roomInfo, setRoomInfo] = useState<RoomInfo | null>(null);
    const [ws, setWs] = useState<any>(null);
    const [isConnectedWS, setIsConnectedWS] = useState<boolean>(false);
    const [message, setMessage] = useState<{ sender: string; content: string; isMe: boolean }[]>([]);
    const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [authenticatedClients, setAuthenticatedClients] = useState<Clients[]>([]);
    const [unauthenticatedClients, setUnauthenticatedClients] = useState<Clients[]>([]);
    const [isOwner, setIsOwner] = useState<boolean>(false);
    const [isTryingToConnect, setIsTryingToConnect] = useState<boolean>(false);

    const getCookie = (name: string) => {
        const cookies = document.cookie.split('; ');
        const cookie = cookies.find(row => row.startsWith(`${name}=`));
        return cookie ? decodeURIComponent(cookie.split('=')[1]) : null;
    };

    const connectToRoom = async (roomID: string) => {
        const APIURL = process.env.REACT_APP_WSAPI_URL;
        const clientName = getCookie('user_name');
        const URL = `${APIURL}/ws?room_id=${roomID}&client_name=${clientName}`;
        const ws = new WebSocket(URL);
        ws.onopen = () => {
            setIsConnectedWS(true);
            console.log('WebSocket connected');
        };
        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            console.log('Received:', data);
            if (data.type === 'message') {
                setMessage((prevMessages) => [...prevMessages, { sender: data.sender, content: data.sentence, isMe: false }]);
            } else if (data.type === 'participants_update') {
                fetchParticipants();
                if(!isAuthenticated) {
                    setupRoom();
                    setIsTryingToConnect(false);
                }
            }
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

        console.log('WebSocket:', ws);
    };

    const sendMessage = (message: string) => {
        let clientName = getCookie('user_name');
        if (clientName === null) {
            clientName = 'Unknown';
        }
        if (ws && ws.readyState === WebSocket.OPEN) {
            setMessage((prevMessages) => [...prevMessages, { sender: clientName, content: message, isMe: true }]);
            ws.send(JSON.stringify({ type: 'message', content: message }));
        } else {
            console.error('WebSocketはまだ開いていません。現在の状態:', ws?.readyState);
        }
    };

    const fetchRoomInfo = async () => {
        const ownerFlag = getCookie('is_owner') === 'true';
        setIsOwner(ownerFlag);
        console.log('isOwner:', ownerFlag);
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
            console.log('Room Data:', roomData);

            setIsAuthenticated(true);
            setRoomInfo(roomData);
            if (roomID) {
                connectToRoom(roomID);
            }
        } else if (roomID) {
            const URL = `${APIURL}/room/${roomID}/isAuth`;
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
                if (isTryingToConnect) {
                    console.log('Failed to authenticate');
                    setIsAuthenticated(false);
                } else {
                    setIsAuthenticated(false);
                    connectToRoom(roomID);
                    setIsTryingToConnect(true);
                    console.log('Not Authenticated');
                }
            }
        } else {
            console.log('Room ID not found');
        }
        setIsLoading(false);
    };

    const fetchParticipants = async () => {
        const roomID = getCookie('room_id');
        const APIURL = process.env.REACT_APP_API_URL;
        const URL = `${APIURL}/room/${roomID}/participants`;
        const response = await fetch(URL, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
        });
        const participants = await response.json();
        console.log('Participants:', participants);
        setAuthenticatedClients(participants.authenticatedClients);
        setUnauthenticatedClients(participants.unauthenticatedClients);
    };

    const handleKick = (clientId: string) => {
        console.log(`Kick client with ID: ${clientId}`);
        const roomID = getCookie('room_id');
        const APIURL = process.env.REACT_APP_API_URL;
        const URL = `${APIURL}/room/${roomID}/kick?client_id=${clientId}`;
        fetch(URL, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
        });
        handleMessageUpdateParticipants();
    };

    const handleApprove = async (clientId: string) => {
        console.log(`Approve client with ID: ${clientId}`);
        const roomID = getCookie('room_id');
        const APIURL = process.env.REACT_APP_API_URL;
        const URL = `${APIURL}/room/${roomID}/auth?client_id=${clientId}`;

        const response = await fetch(URL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'include',
        });
        if (response.ok) {
            console.log('Approved');
        } else {
            console.error('Failed to approve');
        }
        handleMessageUpdateParticipants();
    };

    const handleMessageUpdateParticipants = () => {
        if (ws && ws.readyState === WebSocket.OPEN) {
            console.log('Sending participants_update');
            ws.send(JSON.stringify({ type: 'participants_update' }));
        }
    };

    const setupRoom = () => {
        fetchRoomInfo();
        fetchParticipants();
    };

    useEffect(() => {
        setupRoom();
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
                <Button variant="contained" onClick={setupRoom} sx={{ marginTop: 2 }}>
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
            {/* 画面左上に固定配置するボタン */}
            <IconButton
                className="refresh-button"
                title="メンバー、リクエストのリストを更新"
                onClick={setupRoom}
                sx={{
                    position: 'absolute',
                    bottom: 10,
                    left: 10,
                    width: 60, // ボタンの幅
                    height: 60, // ボタンの高さ
                    bgcolor: 'primary.main', // 背景色
                    color: 'white', // アイコンの色
                    '&:hover': {
                        bgcolor: 'primary.dark', // ホバー時に色を変更
                    },
                }}
            >
                <Refresh />
            </IconButton>

            <div className="members-section">
                <div className="members-header">メンバー ({authenticatedClients.length}人)</div>
                {authenticatedClients.map((client) => (
                    <div key={client.clientid} className="member-item">
                        {client.name} {client.isowner && '(オーナー)'}
                        {!client.isowner && isOwner && (
                            <IconButton
                                color="error"
                                title="退室させる"
                                onClick={() => handleKick(client.clientid)}
                                sx={{ marginLeft: 2 }}
                            >
                                <RemoveCircle />
                            </IconButton>
                        )}
                    </div>
                ))}
            </div>

            <ChatArea message={message} sendMessage={sendMessage} />

            {isOwner && (
                <div className="requests-section">
                    <div className="requests-header">リクエスト ({unauthenticatedClients.length}人)</div>
                    {unauthenticatedClients.map((client) => (
                        <div key={client.clientid} className="request-item">
                            {client.name}{' '}
                            <Button
                                variant="contained"
                                color="primary"
                                onClick={() => handleApprove(client.clientid)}
                                sx={{ marginLeft: 2 }}
                            >
                                承認
                            </Button>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default Chat;
