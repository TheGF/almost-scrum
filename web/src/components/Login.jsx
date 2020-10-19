import React from 'react';
import Modal from 'react-bootstrap/Modal';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';

function Login(logged_in, handle_login) {
    return (
        <Modal show={logged_in} >
            <Modal.Header>
                <Modal.Title>Welcome to Quick Scrum</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <Form>
                    <Form.Group controlId="formBasicEmail">
                        <Form.Label>User</Form.Label>
                        <Form.Control type="text" placeholder="Enter your username" />
                    </Form.Group>
                    <Form.Group controlId="formBasicPassword">
                        <Form.Label>Password</Form.Label>
                        <Form.Control type="password" placeholder="Password" />
                    </Form.Group>
                </Form>
            </Modal.Body>
            <Modal.Footer>
                <Button variant="primary" onClick={handle_login}>
                    Login
                </Button>
            </Modal.Footer>
        </Modal>
    )
}

export default Login;