import React, { useState, useEffect } from 'react';
import './css/create_popup.css';
import { Box, Button, Typography, Modal, IconButton, TextField, Select, MenuItem, InputAdornment, SelectChangeEvent, FormControlLabel, Checkbox,Tooltip } from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import Icon from '@mdi/react';
import { mdiHomePlus } from '@mdi/js';

const CreatePopup: React.FC = () => {
    const [isPopupVisible, setIsPopupVisible] = useState(false);
    const [timeUnits, setTimeUnits] = useState<string>('時間');
    const [expiry, setExpiry] = useState<string>('1');
    const [roomName, setRoomName] = useState<string>(''); 
    const [userName, setUserName] = useState<string>(''); 
    const [requiresAuth, setRequiresAuth] = useState<boolean>(false); 
    const [error, setError] = useState<string>('');
    const [errorSentence, setErrorSentence] = useState<string>('');

    const openPopup = () => setIsPopupVisible(true);
    const closePopup = () => {
        setIsPopupVisible(false);
        setExpiry(expiry);
        setTimeUnits(timeUnits);
    };

    const validateExpiry = (value: string) => {
        const numValue = Number(value);
        if (isNaN(numValue) || numValue <= 0) {
            return "1以上の数を入力してください";
        }
        if (timeUnits === '分　' && numValue > 60) {
            return "60分以内で入力してください";
        }
        if (timeUnits === '時間' && numValue > 24) {
            return "24時間以内で入力してください";
        }
        if (timeUnits === '日　' && numValue > 5) {
            return "5日以内で入力してください";
        }
        return "";
    };

    const handleExpiryChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const value = event.target.value;
        setExpiry(value);
        setError(validateExpiry(value));
    };

    const handleUnitChange = (event: SelectChangeEvent<string>) => {
        setTimeUnits(event.target.value);
    };

    const handleRoomNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRoomName(event.target.value);
    };

    const handleUserNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setUserName(event.target.value);
    };

    const handleAuthChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRequiresAuth(event.target.checked);
    };

    const calculateExpiration = () => {
        const currentDate = new Date();
        const expirationDate = new Date(currentDate);
        if (timeUnits === '時間') {
            expirationDate.setHours(currentDate.getHours() + Number(expiry));
        }
        if (timeUnits === '日　') {
            expirationDate.setDate(currentDate.getDate() + Number(expiry));
        }
        if (timeUnits === '分　') {
            expirationDate.setMinutes(currentDate.getMinutes() + Number(expiry));
        }
        
        return expirationDate.toISOString();
    };



    const createRoom = async () => {
        console.log("ルーム名:", roomName);
        console.log("ユーザー名:", userName);
        console.log("期限:", expiry, timeUnits);
        console.log("入室許可:", requiresAuth);

        const APIURL = process.env.REACT_APP_API_URL;
        const url = `${APIURL}/room`;
        console.log("API URL:", url);

        const expirationTime = calculateExpiration();

        console.log("期限:", expirationTime);

        const data = {
            name: roomName,
            owner: userName,
            expires: expirationTime,
            requiresAuth: requiresAuth,
        };

        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
            credentials: 'include',

        });

        if (response.ok) {
            const roomData = await response.json();
            console.log('Room creation response:', roomData);
            console.log(roomData.ID);
           
            window.location.href = `/chat`;
            
        }else{
            alert('ルームの作成に失敗しました');
        }
        


    }

    useEffect(() => {
        if (timeUnits === '分　' && Number(expiry) > 60) {
            setExpiry(String(60));
        }
        if (timeUnits === '時間' && Number(expiry) > 24) {
            setExpiry(String(24));
        }
        if (timeUnits === '日　' && Number(expiry) > 5) {
            setExpiry(String(5));
        }
        setError(validateExpiry(expiry));
    }, [timeUnits, expiry]);

    useEffect(() => {
        let messages: string[] = [];
        if (roomName === "") {
            messages.push("ルーム名を入力してください");
        }
        if (userName === "") {
            messages.push("ユーザー名を入力してください");
        }
        if (error) {
            messages.push("期限フィールドには"+error);
        }
        setErrorSentence(messages.join("\n"));
    }, [error, roomName, userName]);
    


    return (
        <div>
            <Button className="home-button" onClick={openPopup}>
                                <Icon path={mdiHomePlus} size={2} />
                                <Typography variant="h6">ルーム作成</Typography>
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
                        ルームを作成
                    </Typography>

                    <TextField
                        className="popup-input"
                        label="ルーム名 (必須)"
                        variant="outlined"
                        value={roomName}
                        onChange={handleRoomNameChange}
                    />

                    <TextField
                        className="popup-input"
                        label="ユーザー名 (必須)"
                        variant="outlined"
                        value={userName}
                        onChange={handleUserNameChange}
                    />

                    <Box className="expiry-container">
                    <TextField
                        className="popup-input"
                        label="期限"
                        variant="outlined"
                        type="number"
                        value={expiry}
                        onChange={handleExpiryChange}
                        error={!!error}
                        helperText={error || " "} 
                        InputProps={{
                            endAdornment: <InputAdornment position="end">{timeUnits}</InputAdornment>,
                        }}
                        inputProps={{
                            min: 1,
                            max: timeUnits === '分　' ? 60 : timeUnits === '時間' ? 24 : 5,
                            style: { textAlign: "right" },
                        }}
                    />


                        <Select
                            className="popup-select"
                            value={timeUnits}
                            onChange={handleUnitChange}
                            displayEmpty
                        >
                            <MenuItem value="分　">分　</MenuItem>
                            <MenuItem value="時間">時間</MenuItem>
                            <MenuItem value="日　">日　</MenuItem>
                        </Select>
                    </Box>

                    <FormControlLabel
                        control={
                            <Checkbox
                                size="large"
                                checked={requiresAuth} // チェックボックスの状態を反映
                                onChange={handleAuthChange} // チェックボックス変更時に状態を更新
                            />
                        }
                        label="チェックすると入室許可が必要になります"
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
                                endIcon={<Icon path={mdiHomePlus} size={1} />}
                                onClick={createRoom}
                                className="popup-button"
                                disabled={!!error || !expiry || !roomName || !userName} // エラーや未入力時は無効化
                            >
                                ルームを作成
                            </Button>
                        </span>
                    </Tooltip>

                </Box>
            </Modal>
        </div>
    );
}

export default CreatePopup;
