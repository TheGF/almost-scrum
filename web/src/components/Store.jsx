import React, { useEffect, useState } from 'react';
import Badge from 'react-bootstrap/Badge';
import Button from 'react-bootstrap/Button';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';
import Card from 'react-bootstrap/Card';
import Form from 'react-bootstrap/Form';
import ListGroup from 'react-bootstrap/ListGroup';
import Modal from 'react-bootstrap/Modal';
import { AiOutlineLeftSquare, AiOutlineRightSquare, AiOutlineUpSquare } from 'react-icons/ai';
import { VscNewFile, VscNewFolder } from 'react-icons/vsc';
import Server from './Server';

function CreateStory(props) {
    const { project, store, show, setShow } = props;
    const [title, setTitle] = useState('');

    function postStory() {
        Server.createStory(project, store, title)
            .then(setShow)
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
                setShow();
                setTitle('');
            }}>
                Cancel
            </Button>
        </Modal.Footer>
    </Modal>;
}

function Store(props) {
    const { project, store, story, onSelect, toLeft, toRight, reload } = props;
    const [stories, setStories] = useState([]);
    const [showCreateStory, setShowCreateStory] = useState(false);

    function fetchStories() {
        Server.getStoriesList(project, store)
            .then(setStories);
    }
    useEffect(fetchStories, [showCreateStory]);

    const storyList = stories.map(s =>
        <ListGroup.Item active={s == story} style={{ padding: '0.3em', cursor: 'pointer' }}
            key={s} onClick={_ => onSelect && onSelect(s)}>
            {s}
            <span className="float-right" onClick={e => e.stopPropagation} >
                {toLeft &&
                    <a href="#" title={`Move the story to ${toLeft}`}
                        onClick={_ =>
                            Server.moveStory(project, store, s, toLeft)
                                .then(reload)}>
                        <AiOutlineLeftSquare />
                    </a>}
            &nbsp;
            <a href="#" title="Move the story to the top of the list"
                    onClick={_ =>
                        Server.touchStory(project, store, s)
                            .then(fetchStories)}>
                    <AiOutlineUpSquare />
                </a>
            &nbsp;
            {toRight &&
                    <a href="#" title={`Move the story to ${toRight}`}
                        onClick={_ =>
                            Server.moveStory(project, store, s, toRight)
                                .then(reload)}>
                        <AiOutlineRightSquare />
                    </a>}
            </span>
        </ListGroup.Item>)

    return <Form>
        <Card>
            <Badge>{project}/{store}</Badge>
            <CreateStory show={showCreateStory} setShow={setShowCreateStory}
                project={project} store={store} />
            <ButtonToolbar className="d-flex">
                <Button className="ml-1 mb-1" variant="primary" onClick={_ => setShowCreateStory(true)}>
                    <VscNewFile /> Create Story
                </Button>
                <Button className="ml-1 mb-1" variant="primary" onClick={_ => setShowCreateStory(true)}>
                    <VscNewFolder /> Create Folder
                </Button>
            </ButtonToolbar>
            <ListGroup variant="flush">
                {storyList}
            </ListGroup>
        </Card>
    </Form>
}

export default Store;