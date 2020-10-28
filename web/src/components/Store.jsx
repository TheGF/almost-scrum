import React, { useEffect, useState } from 'react';
import ListGroup from 'react-bootstrap/ListGroup';
import Button from 'react-bootstrap/Button';
import ButtonGroup from 'react-bootstrap/ButtonGroup';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';
import InputGroup from 'react-bootstrap/InputGroup';
import Badge from 'react-bootstrap/Badge';
import Form from 'react-bootstrap/Form';
import Card from 'react-bootstrap/Card';
import Modal from 'react-bootstrap/Modal';
import axios from 'axios';
import { getConfig, loginWhenUnauthorized } from './axiosUtils';
import { GrNewWindow } from 'react-icons/gr';
import { BiArchive} from 'react-icons/bi';
import { AiOutlineToTop } from 'react-icons/ai';
import { emptyStory } from './consts';


function CreateStory(props) {
    const { project, store, show, afterCreate } = props;
    const [title, setTitle] = useState('');

    function postStory() {
        axios.post(`/api/v1/projects/${project}/${store}?title=${title}`, emptyStory, getConfig())
            .then(() => { afterCreate(); setTitle('') })
            .catch(loginWhenUnauthorized);
    }

    return <Modal show={show}>
        <Modal.Header>
            <Modal.Title>Create a new story</Modal.Title>
        </Modal.Header>
        <Modal.Body>
            <Form>
                <Form.Group controlId="formBasicEmail">
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
        </Modal.Footer>
    </Modal>;
}

function Store(props) {
    const { project, store } = props;
    const [stories, setStories] = useState([]);
    const [selected, setSelected] = useState(null);
    const [showNewStory, setShowNewStory] = useState(false);

    function afterCreate() {
        setShowNewStory(false);
    }

    function touchStory(s) {
        axios.post(`/api/v1/projects/${project}/${store}/${s}?touch`, null, getConfig())
            .then(r => setStories(r.data))
            .catch(loginWhenUnauthorized)
    }

    function archiveStory(s) {
        axios.post(`/api/v1/projects/${project}/${store}/${s}?archive`, null, getConfig())
            .then(r => setStories(r.data))
            .catch(loginWhenUnauthorized)
    }

    function selectStory(s) {
        props.onSelect && props.onSelect(s);
        setSelected(s);
    }

    function fetchStories() {
        axios.get(`/api/v1/projects/${project}/${store}`, getConfig())
            .then(r => setStories(r.data))
            .catch(loginWhenUnauthorized)
    }
    useEffect(fetchStories, [showNewStory]);

    const storyList = stories.map(s => <ListGroup.Item 
            key={s} active={s==selected} onClick={() => selectStory(s)}>
        {s}
        <span className="float-right" >
            <a href="#" onClick={() => archiveStory(s)}>
                <BiArchive />
            </a>
            <a href="#" onClick={() => touchStory(s)}>
                <AiOutlineToTop />
            </a>
        </span>
    </ListGroup.Item>)

    return <Form>
        <Card>
            <Badge>{store}</Badge>
            <CreateStory show={showNewStory} project={project} store={store}
                afterCreate={afterCreate} />
            <Button variant="secondary" onClick={() => setShowNewStory(true)}>
                <GrNewWindow />
                New
            </Button>
            <ListGroup>
                {storyList}
            </ListGroup>
        </Card>
    </Form>
}

export default Store;