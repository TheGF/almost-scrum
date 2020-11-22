import React from 'react';
import ReactDOM from "react-dom";
import App from './components/App';
import './custom.scss';
import './components/custom.css';

const $container = document.getElementById("app-container");
$container ? ReactDOM.render(<App />, $container) : false;