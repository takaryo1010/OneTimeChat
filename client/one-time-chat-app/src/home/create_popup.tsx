import React, { useState, useEffect } from 'react';
import './css/create_popup.css';
import { Box, Button, Typography, Modal, IconButton, TextField, Select, MenuItem, InputAdornment, SelectChangeEvent, FormControlLabel, Checkbox } from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import Icon from '@mdi/react';
import { mdiHomePlus } from '@mdi/js';

const CreatePopup: React.FC = () => {
    const [isPopupVisible, setIsPopupVisible] = useState(false);
    const [timeUnits, setTimeUnits] = useState<string>('時間');
    const [expiry, setExpiry] = useState<string>('1');
    const [roomName, setRoomName] = useState<string>(''); 
    const [userName, setUserName] = useState<string>(''); 
    const [requiresAuth, setRequiresAuth] = useState<boolean>(false); // チェックボックスの状態を管理
    const [error, setError] = useState<string>('');

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
        if (timeUnits === '分' && numValue > 60) {
            return "60分以内で入力してください";
        }
        if (timeUnits === '時間' && numValue > 24) {
            return "24時間以内で入力してください";
        }
        if (timeUnits === '日' && numValue > 5) {
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

    useEffect(() => {
        setError(validateExpiry(expiry));
    }, [timeUnits, expiry]);

    return (
        <div>
            <Button
                variant="contained"
                color="primary"
                endIcon={<Icon path={mdiHomePlus} size={1} />}
                onClick={openPopup}
            >
                ルームを作成
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
                        label="ルーム名"
                        variant="outlined"
                        value={roomName}
                        onChange={handleRoomNameChange}
                    />

                    <TextField
                        className="popup-input"
                        label="ユーザー名"
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
                            helperText={error}
                            InputProps={{
                                endAdornment: <InputAdornment position="end">{timeUnits}</InputAdornment>,
                            }}
                            inputProps={{
                                min: 1,
                            }}
                        />

                        <Select
                            className="popup-select"
                            value={timeUnits}
                            onChange={handleUnitChange}
                            displayEmpty
                        >
                            <MenuItem value="分">分　</MenuItem>
                            <MenuItem value="時間">時間</MenuItem>
                            <MenuItem value="日">日　</MenuItem>
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

                    <Button
                        variant="contained"
                        color="primary"
                        endIcon={<Icon path={mdiHomePlus} size={1} />}
                        onClick={closePopup}
                        className="popup-button"
                        disabled={!!error || !expiry || !roomName || !userName} // エラーや未入力時は無効化
                    >
                        ルームを作成
                    </Button>
                </Box>
            </Modal>
        </div>
    );
}

export default CreatePopup;
