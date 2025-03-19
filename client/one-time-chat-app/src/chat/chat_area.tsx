import React, { useState } from 'react';
import './css/chat_area.css';

interface Message {
  sender: string;
  content: string;
  isMe: boolean;
}

interface ChatAreaProps {
  message: Message[];
  sendMessage: (content: string) => void;
}

const ChatArea: React.FC<ChatAreaProps> = ({ message, sendMessage }) => {
  const [inputMessage, setInputMessage] = useState<string>('');
  const [expandedMessages, setExpandedMessages] = useState<boolean[]>(message.map(() => false));

  const changeInputMessage = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInputMessage(e.target.value);
  };

  const handleSendMessage = () => {
    if (inputMessage.trim() !== '') {
      console.log('send message:', inputMessage);
      sendMessage(inputMessage);
      setInputMessage('');
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault(); // デフォルトの改行を防ぐ
      handleSendMessage();
    }
  };

  // テキストエリアの高さを動的に調整
  const handleInput = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInputMessage(e.target.value);
    const textarea = e.target;
    textarea.style.height = 'auto'; // 高さをリセット
    textarea.style.height = `${Math.min(textarea.scrollHeight, 10 * 24)}px`; // 10行分以上にならないように設定
  };

  // もっと見る/閉じるボタンをクリックしたときの処理
  const toggleExpandMessage = (index: number) => {
    setExpandedMessages((prev) => {
      const newState = [...prev];
      newState[index] = !newState[index];
      return newState;
    });
  };

  // 5行を超えたメッセージの判定
  const getMessageContent = (content: string, index: number) => {
    const maxLines = 5;
    const maxChars = 200;
    const lines = content.split('\n');
    const charCount = content.length;

    if (lines.length > maxLines || charCount > maxChars) {
      const displayedContent = expandedMessages[index]
        ? content
        : lines.slice(0, maxLines).join('\n').substring(0, maxChars);
      return (
        <div>
          <div className="message-content">{displayedContent}</div>
          <button className="expand-button" onClick={() => toggleExpandMessage(index)}>
            {expandedMessages[index] ? '閉じる' : '続きを見る'}
          </button>
        </div>
      );
    }
    return <div className="message-content">{content}</div>;
  };

  return (
    <div className="chat-area">
      <div className="chat-messages">
        {message.map((m, index) => (
          <div key={index} className={`message ${m.isMe ? 'me' : 'other'}`}>
            <span className="sender">{m.sender}:</span>
            {getMessageContent(m.content, index)}
          </div>
        ))}
      </div>
      <div className="message-input-area">
        <textarea
          className="message-input"
          placeholder="メッセージを入力...（Shift + Enter で改行、Enter で送信）"
          value={inputMessage}
          onChange={changeInputMessage}
          onKeyDown={handleKeyDown}
          onInput={handleInput} // テキストエリアの高さを調整
          rows={1} // 最初は1行分
        />
        <button className="send-button" onClick={handleSendMessage}>送信</button>
      </div>
    </div>
  );
};

export default ChatArea;