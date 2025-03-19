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
            <h1>ようこそ OneTimeChat へ!</h1>
            <p>本サービスは使い捨てのチャットルームを提供します</p>

            {/* ボタンを横並びに配置 */}
            <Box className="button-container">

            {/* ルーム作成のポップアップ */}
            <CreatePopup />

            <JoinPopup/>
            </Box>
            <p>もし部屋を抜けてしまって戻りたい場合は下のボタンをクリックしてください</p>
            <Button className="sub-button" href="/chat">部屋に戻る</Button>
            <p>(部屋を抜けた後にルーム作成や参加をしていると戻れない場合があります)</p>
        </div>
    );
}

export default Home;
