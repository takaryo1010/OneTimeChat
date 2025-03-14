import React, { useState } from 'react';
import { Box, Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, TextField } from '@mui/material';
import { mdiHomePlus, mdiLogin } from '@mdi/js';
import { Typography } from '@mui/material';
import { mdiClose } from '@mdi/js';
import Icon from '@mdi/react';

const JoinPopup: React.FC = () => {
    const [isPopupVisible, setIsPopupVisible] = useState(false);
    const [roomID, setRoomID] = useState<string>('');
    const [userName, setUserName] = useState<string>('');
    const [error, setError] = useState<string>('');
    const [errorSentence, setErrorSentence] = useState<string>('');
    const openPopup = () => setIsPopupVisible(true);
    const closePopup = () => {
        setIsPopupVisible(false);
    };
    const handleRoomIDChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRoomID(event.target.value);
    };
    const handleUserNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setUserName(event.target.value);
    };
    const handleJoin = () => {
        if (roomID === '') {
            setError('ルームIDを入力してください');
            setErrorSentence('');
            return;
        }
        if (userName === '') {
            setError('ユーザー名を入力してください');
            setErrorSentence('');
            return;
        }
        setError('');
        setErrorSentence('');
        console.log("ルーム参加");
    };
    return (
        <div>
            <Button className="home-button" onClick={openPopup}>
                <Icon path={mdiLogin} size={2} />
                <Typography variant="h6">ルーム参加</Typography>
            </Button>
            <Dialog open={isPopupVisible} onClose={closePopup}>
                <DialogTitle>ルーム参加</DialogTitle>
                <DialogContent>

                    <TextField
                        autoFocus
                        margin="dense"
                        id="roomID"
                        label="参加したいルームのID"
                        type="text"
                        fullWidth
                        onChange={handleRoomIDChange}
                    />
                    <TextField
                        margin="dense"
                        id="userName"
                        label="あなたのユーザー名"
                        type="text"
                        fullWidth
                        onChange={handleUserNameChange}
                    />
                    <Typography color="error">{error}</Typography>
                    <Typography>{errorSentence}</Typography>
                </DialogContent>
                <DialogActions>
                    <Button onClick={closePopup}>キャンセル</Button>
                    <Button onClick={handleJoin}>参加</Button>
                </DialogActions>
            </Dialog>
        </div>
    );



}

export default JoinPopup;