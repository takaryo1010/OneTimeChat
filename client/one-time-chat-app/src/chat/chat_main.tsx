import React, { use, useEffect, useState } from 'react';
import { CircularProgress, Button, Typography, Box, IconButton, Modal,Snackbar } from '@mui/material';
import { Refresh, RemoveCircle,Info,Close } from '@mui/icons-material';  // RemoveCircle アイコンを使用
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
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
    const [isInfoVisible, setIsInfoVisible] = useState<boolean>(false);
    const [openSnackbar, setOpenSnackbar] = useState(false)
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
        };
        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
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
        };
        ws.onerror = (error) => {
            setIsConnectedWS(false);
            console.error('WebSocket error:', error);
        };
        setWs(ws);
        handleMessageUpdateParticipants(ws)

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
        const roomID = getCookie('room_id');
        const APIURL = process.env.REACT_APP_API_URL;
        
            const RoomInfoURL = `${APIURL}/room/${roomID}`;
            const RoomInforesponse = await fetch(RoomInfoURL, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
            });

            const roomData = await RoomInforesponse.json();

            setIsAuthenticated(true);
            setRoomInfo(roomData);
            // if (roomID) {
            //     connectToRoom(roomID);
            // }
       
            const IsAuthURL = `${APIURL}/room/${roomID}/isAuth`;
            const response = await fetch(IsAuthURL, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
            });
            const AuthInfo = await response.json();
            if (AuthInfo.isAuth) {
                setIsAuthenticated(true);

            } else {
                if (isTryingToConnect) {
                    setIsAuthenticated(false);
                } else {
                    setIsAuthenticated(false);
                    setIsTryingToConnect(true);
                }
            }
        if (roomID) {
            connectToRoom(roomID);
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
        setAuthenticatedClients(participants.authenticatedClients);
        setUnauthenticatedClients(participants.unauthenticatedClients);
        handleMessageUpdateParticipants(ws);
    };

    const handleKick = (clientId: string) => {
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
        handleMessageUpdateParticipants(ws);
    };

    const handleApprove = async (clientId: string) => {
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
        handleMessageUpdateParticipants(ws);
    };

    const handleMessageUpdateParticipants = (websocket:WebSocket) => {

        if (websocket &&websocket.readyState === WebSocket.OPEN) {
            console.log('Sending participants_update');
            websocket.send(JSON.stringify({ type: 'participants_update' }));
        }
    };
    const handleCopyURL = () => {
        const roomURL = `${window.location.origin}?room_id=${roomInfo?.ID}`;
        navigator.clipboard.writeText(roomURL).then(() => {
            setOpenSnackbar(true);
        });
    };
    const setupRoom = () => {
        fetchRoomInfo();
        fetchParticipants();
    };

    const optimaizedRoomExpiration = () => {
        return roomInfo?.expires ? new Date(roomInfo.expires).toLocaleString() : '不明';
    };

    useEffect(() => {
        setupRoom();
    }, []);

//ws.readyStateがOPENの時にhandleMessageUpdateParticipantsを実行する

    useEffect(() => {
        if (isConnectedWS) {

            fetchParticipants();
        }
    }
    , [isConnectedWS]);






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

        <div>

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
                {/*詳細ボタン*/}
                <IconButton
                    className="refresh-button"
                    title="ルーム詳細"
                    onClick={() => setIsInfoVisible(!isInfoVisible)}
                    sx={{
                        position: 'absolute',
                        bottom: 10,
                        left: 80,
                        width: 60, // ボタンの幅
                        height: 60, // ボタンの高さ
                        bgcolor: 'primary.main', // 背景色
                        color: 'white', // アイコンの色
                        '&:hover': {
                            bgcolor: 'primary.dark', // ホバー時に色を変更
                        },
                    }}
                >
                    <Info />
                    
                </IconButton>
                <Modal
                    open={isInfoVisible}
                    onClose={() => setIsInfoVisible(false)}
                    aria-labelledby="modal-title"
                    aria-describedby="modal-description"
                >
                    <Box className="info-popup">
                        <Box className="close-button">
                            <IconButton onClick={() => setIsInfoVisible(false)}>
                                    <Close />
                            </IconButton>
                        </Box>
                        <Typography id="modal-title" variant="h4" component="h2" className="popup-title">
                            ルーム詳細
                        </Typography>
                        <Typography variant="h6" className="popup-content">
                            ルーム名: {roomInfo?.name}
                        </Typography>
                        
                        
                        <Box display="flex" alignItems="center" className="popup-content">
                            <Typography variant="h6" >
                                ルームID: {roomInfo?.ID}
                            </Typography>
                            <IconButton onClick={handleCopyURL} sx={{ ml: 1 }}>
                                <Typography>
                                    ルームURLをコピー
                                </Typography>
                                <ContentCopyIcon />
                            </IconButton>
                        </Box>

                        {/* コピー完了の通知 */}
                        <Snackbar
                            open={openSnackbar}
                            autoHideDuration={2000}
                            onClose={() => setOpenSnackbar(false)}
                            message="ルームURLをコピーしました"
                        />
                        <Typography variant="h6" className="popup-content">
                            オーナー: {roomInfo?.owner}
                        </Typography>
                        <Typography variant="h6" className="popup-content">
                            有効期限: {optimaizedRoomExpiration()}
                        </Typography>
                        <Typography variant="h6" className="popup-content">
                            認証: {roomInfo?.requiresAuth ? '必要' : '不要'}
                        </Typography>
                    </Box>
                </Modal>
                

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
        </div>
    );
};

export default Chat;
