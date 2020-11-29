import React, { useState } from 'react';
import Modal from 'react-bootstrap/Modal';
import Form from 'react-bootstrap/Form';
import Alert from 'react-bootstrap/Alert';
import Button from 'react-bootstrap/Button';
import axios from 'axios';

 
function Login(props) {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [errorMessage, setError] = useState('');

    function handleAuthenticate() {
        const data = { username: username, password: password };
        axios.post("/auth/login", data)
            .then(r => props.onLogin(username, r.data))
            .catch(r => setError(`Invalid Credentials: ${r.message}`));
    }
    return <Modal show={true}>
        <Modal.Header>
            <Modal.Title>Welcome to A Scrum</Modal.Title>
        </Modal.Header>
        <Modal.Body>
            <Form>
                <Form.Group controlId="formBasicEmail">
                    <Form.Label>User</Form.Label>
                    <Form.Control type="text" placeholder="Enter your username"
                        value={username} onChange={e => setUsername(e.target.value)} />
                </Form.Group>
                <Form.Group controlId="formBasicPassword">
                    <Form.Label>Password</Form.Label>
                    <Form.Control type="password" placeholder="Password"
                        value={password} onChange={e => setPassword(e.target.value)} />
                </Form.Group>
            </Form>
            {errorMessage && <Alert variant="warning">{errorMessage}</Alert>}
        </Modal.Body>
        <Modal.Footer>
            <Button variant="primary" onClick={handleAuthenticate}>
                Login
                </Button>
        </Modal.Footer>
    </Modal>;
}

export default Login;