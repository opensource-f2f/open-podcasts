import React from 'react';
import ReactDOM from "react-dom";
import App from './App'
import RSSList from './RSSList'

ReactDOM.render(
    <React.StrictMode>
        <App />
        <RSSList />
    </React.StrictMode>,
    document.getElementById('root')
);