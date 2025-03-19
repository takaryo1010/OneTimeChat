import React from 'react';
import CreatePopup from './create_popup.tsx';
import JoinPopup from './join_popup.tsx';
import './css/home.css';
import { Box, Button, Typography } from "@mui/material";
import Icon from '@mdi/react';
import { mdiHomePlus, mdiLogin } from '@mdi/js';

const Home: React.FC = () => {
    return (
        <div className="home">
            <Typography variant="h2" gutterBottom>
                ようこそ OneTimeChat へ!
            </Typography>
            <Typography variant="body1" >
                本サービスは使い捨てのチャットルームを提供します
            </Typography>

            {/* ボタンを横並びに配置 */}
            <Box className="button-container">
                {/* ルーム作成のポップアップ */}
                <CreatePopup />
                <JoinPopup />
            </Box>

            <Typography variant="body2">
                もし部屋を抜けてしまって戻りたい場合は下のボタンをクリックしてください
            </Typography>
            <Button className="sub-button" href="/chat">
                部屋に戻る
            </Button>
            <Typography variant="body2" color="textSecondary">
                (部屋を抜けた後にルーム作成や参加をしていると戻れない場合があります)
            </Typography>
        </div>
    );
}

export default Home;
