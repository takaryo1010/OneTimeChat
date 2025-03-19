import React, { useState } from 'react';
import { Box, Button, Typography, Modal, IconButton, TextField, Tooltip } from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import Icon from '@mdi/react';
import { mdiLogin } from '@mdi/js';
import "./css/join_popup.css";

const JoinPopup: React.FC = () => {
    const [isPopupVisible, setIsPopupVisible] = useState(false);
    const [roomID, setRoomID] = useState<string>('');
    const [userName, setUserName] = useState<string>('');
    const [error, setError] = useState<string>('');
    const [errorSentence, setErrorSentence] = useState<string>('');

    const openPopup = () => setIsPopupVisible(true);
    const closePopup = () => {
        setIsPopupVisible(false);
        setError('');
    };

    // ルームIDのバリデーション関数
    const validateRoomID = (id: string): string | null => {
        const regex = /^[ABCDEFGHJKLMNPQRSTUVWXY0123456789]{5}$/; // 5文字 & 指定された文字のみ
        if (!regex.test(id)) {
            return "ルームIDは5文字の英大文字（I, O, Z 除く）と数字である必要があります";
        }
        return null;
    };

    const handleRoomIDChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRoomID(event.target.value.toUpperCase());
        setError(validateRoomID(event.target.value.toUpperCase()) || '');
    };

    const handleUserNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setUserName(event.target.value);
    };

    const joinRoom = async () => {
        const APIURL = process.env.REACT_APP_API_URL;
        const URL = `${APIURL}/room/${roomID}`;
        const response = await fetch(URL,
            {
                method: 'POST',
                body: JSON.stringify({ client_name: userName }),
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
            }
        );
        if (response.ok) {
            const roomData = await response.json();
            console.log(roomData);
            window.location.href = `/chat`;
        }
        else {
            setErrorSentence('ルームへの参加に失敗しました');
        }
        
    };

    return (
        <div>
            <Button className="home-button" onClick={openPopup}>
                <Icon path={mdiLogin} size={2} />
                <Typography variant="h6">ルーム参加</Typography>
            </Button>

            <Modal
                open={isPopupVisible}
                onClose={closePopup}
                aria-labelledby="modal-title"
                aria-describedby="modal-description"
            >
                <Box className="popup">
                    <Box className="close-button">
                        <IconButton onClick={closePopup}>
                            <CloseIcon />
                        </IconButton>
                    </Box>

                    <Typography id="modal-title" variant="h4" component="h2" className="popup-title">
                        ルームに参加
                    </Typography>

                    <TextField
                        className="popup-input"
                        label="ルームID (必須)"
                        variant="outlined"
                        value={roomID}
                        onChange={handleRoomIDChange}
                        error={!!error}
                        helperText={error || " "}
                    />

                    <TextField
                        className="popup-input"
                        label="ユーザー名 (必須)"
                        variant="outlined"
                        value={userName}
                        onChange={handleUserNameChange}
                    />

                    <Tooltip
                        title={
                            !!errorSentence ? (
                                <Typography component="span">
                                    {errorSentence.split("\n").map((line, index) => (
                                        <React.Fragment key={index}>
                                            {line}
                                            <br />
                                        </React.Fragment>
                                    ))}
                                </Typography>
                            ) : ""
                        }
                        arrow
                    >
                        <span className="tooltip-wrapper">
                            <Button
                                variant="contained"
                                color="primary"
                                endIcon={<Icon path={mdiLogin} size={1} />}
                                onClick={joinRoom}
                                className="popup-button"
                                disabled={!!error || !roomID || !userName}
                            >
                                参加
                            </Button>
                        </span>
                    </Tooltip>
                </Box>
            </Modal>
        </div>
    );
};

export default JoinPopup;
