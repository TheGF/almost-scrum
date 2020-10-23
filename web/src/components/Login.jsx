import React, { useState } from 'react';
import Modal from 'react-bootstrap/Modal';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import axios from 'axios';


function Login(logged_in, success_login) {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');

    function handle_login() {
        const data = {username: username, password: password};
        axios.post("/api/v1/login", data)
            .then(response => success_login(username));
    }

    return (
        <Modal show={logged_in} >
            <Modal.Header>
                <Modal.Title>Welcome to A Scrum</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <Form>
                    <Form.Group controlId="formBasicEmail">
                        <Form.Label>User</Form.Label>
                        <Form.Control type="text" placeholder="Enter your username"
                            value={username} onChange={e => setUsername(e.value)} />
                    </Form.Group>
                    <Form.Group controlId="formBasicPassword">
                        <Form.Label>Password</Form.Label>
                        <Form.Control type="password" placeholder="Password"
                            value={password} onChange={e => setPassword(e.value)} />
                    </Form.Group>
                </Form>
            </Modal.Body>
            <Modal.Footer>
                <Button variant="primary" onClick={
                    () => handle_login(username, password)
                }>
                    Login
                </Button>
            </Modal.Footer>
        </Modal>
    )
}

export default Login;