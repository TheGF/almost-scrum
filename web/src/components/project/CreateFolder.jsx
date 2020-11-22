import React, { useState } from 'react';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import Modal from 'react-bootstrap/Modal';
import Server from '../Server';

function CreateFolder(props) {
    const { project, store, folders, show, setShow } = props;
    const [name, setName] = useState('');

    function createFolder() {
        const path = [...folders, name].join('/')
        Server.createFolder(project, store, path)
            .then(_ => setShow(false))
            .then(_ => setName(''));
    }

    return <Modal show={show}>
        <Modal.Header>
            <Modal.Title>Create a new folder</Modal.Title>
        </Modal.Header>
        <Modal.Body>
            <Form>
                <Form.Group controlId="formName">
                    <Form.Label>Name</Form.Label>
                    <Form.Control type="text" placeholder="Folder name"
                        value={name} onChange={e => setName(e.target.value)} />
                </Form.Group>
            </Form>
        </Modal.Body>
        <Modal.Footer>
            <Button variant="primary" onClick={createFolder}>
                Create
            </Button>
            <Button variant="secondary" onClick={_ => {
                setShow(false);
                setName('');
            }}>
                Cancel
            </Button>
        </Modal.Footer>
    </Modal>;
}

export default CreateFolder;