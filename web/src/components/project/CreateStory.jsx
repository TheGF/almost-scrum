import React, { useState } from 'react';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import Modal from 'react-bootstrap/Modal';
import { emptyStory } from '../consts';
import Server from '../Server';

function CreateStory(props) {
    const { project, store, show, setShow, folders } = props;
    const [title, setTitle] = useState('');

    function postStory() {
        const path = folders.join('/');
        Server.createStory(project, store, path, title, emptyStory)
            .then(_ => setShow(false))
            .then(_ => setTitle(''));
    }

    return <Modal show={show}>
        <Modal.Header>
            <Modal.Title>Create a new story</Modal.Title>
        </Modal.Header>
        <Modal.Body>
            <Form>
                <Form.Group controlId="formTitle">
                    <Form.Label>Title</Form.Label>
                    <Form.Control type="text" placeholder="Short title for the story"
                        value={title} onChange={e => setTitle(e.target.value)} />
                </Form.Group>
            </Form>
        </Modal.Body>
        <Modal.Footer>
            <Button variant="primary" onClick={postStory}>
                Create
            </Button>
            <Button variant="secondary" onClick={_ => {
                setShow(false);
                setTitle('');
            }}>
                Cancel
            </Button>
        </Modal.Footer>
    </Modal>;
}

export default CreateStory;