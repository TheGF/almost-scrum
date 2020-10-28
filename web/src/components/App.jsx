import React, { useState } from 'react';
import Desktop from './Desktop';
import Login from './Login';

const App = () => {

    function handleLogin(username, token) {
        localStorage.token = token;
        localStorage.username = username;
        setUsername(username)
    }

    const [username, setUsername] = useState(localStorage.username);
    const content = !username ? <Login onLogin={handleLogin}/> : <Desktop username={username}/>
    return <div className="bg-dark container-fluid min-vh-100" >{content}</div>;
};

export default App;