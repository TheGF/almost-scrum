import React, { useState } from 'react';
import Desktop from './Desktop';
import Login from './Login';

globalThis.axiosHeaders = {}

const App = () => {

    function handleLogin(username, token) {
        globalThis.axiosHeaders['Authorization'] = `Bearer ${token}`;
        setUsername(username)
    }

    const [username, setUsername] = useState(null);
    return !username ? <Login onLogin={handleLogin}/> : <Desktop username={username}/>;
};

export default App;