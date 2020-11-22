import React, { useEffect, useState } from 'react';
import Badge from 'react-bootstrap/Badge';
import Button from 'react-bootstrap/Button';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';
import Card from 'react-bootstrap/Card';
import Form from 'react-bootstrap/Form';
import ListGroup from 'react-bootstrap/ListGroup';
import { AiOutlineLeftSquare, AiOutlineRightSquare, AiOutlineUpSquare } from 'react-icons/ai';
import { BiRefresh } from 'react-icons/bi';
import { FaFileAlt, FaFolder } from 'react-icons/fa';
import { VscNewFile, VscNewFolder } from 'react-icons/vsc';
import Server from '../Server';
import CreateFolder from './CreateFolder';
import CreateStory from './CreateStory';

function Store(props) {
    const { project, store, story, onSelect, toLeft, toRight, reload } = props;
    const [items, setItems] = useState([]);
    const [showCreateStory, setShowCreateStory] = useState(false);
    const [showCreateFolder, setShowCreateFolder] = useState(false);
    const [folders, setFolders] = useState([]);

    function getPath(name=null) {
        const p = folders.length ? `/${folders.join('/')}` : '';
        return name ? `${p}/${name}` : p;
    }

    function fetchItems() {
        Server.getStoreList(project, store, getPath())
            .then(setItems);
    }
    useEffect(fetchItems, [showCreateStory, folders, showCreateFolder]);

    function onItemClick(item) {
        const { name, dir } = item;
        if (dir) {
            const fs = name == '..' ? folders.slice(0, -1) : [...folders, name];
            setFolders(fs);
        } else if (onSelect) {
            onSelect(getPath(name));
        }
    }

    function renderIcons(name, dir) {
        const icons = [];

        if (toLeft && !dir) {
            icons.push(
                <a key="left" href="#" title={`Move the story to ${toLeft}`}
                    onClick={_ =>
                        Server.moveStory(project, store, getPath(name), toLeft)
                            .then(reload)}>
                    <AiOutlineLeftSquare />
                </a>);
        }
        if (name != '..') {
            icons.push(
                <a key="top" href="#" title="Move to the top of the list"
                    onClick={_ =>
                        Server.touchStory(project, store, getPath(name))
                            .then(fetchStories)}>
                    <AiOutlineUpSquare />
                </a>);
        }
        if (toRight && !dir) {
            icons.push(
            <a key="right" href="#" title={`Move the story to ${toRight}`}
                onClick={_ => 
                    Server.moveStory(project, store, getPath(name), toRight)
                        .then(reload)}>
                <AiOutlineRightSquare />
            </a>);
        }
        return icons;
    }

    const is = folders.length ? [{ name: '..', dir: true }, ...items] : items;
    const storyList = is.map(s => {
        const { name, dir } = s;

        return <ListGroup.Item active={name == story} style={{ padding: '0.3em', cursor: 'pointer' }}
            key={name} onClick={_ => onItemClick(s)}>
            {dir ? <FaFolder /> : <FaFileAlt />} &nbsp;
            {name}
            <span className="float-right" onClick={e => e.stopPropagation} >
                {renderIcons(name, dir)}
            </span>
        </ListGroup.Item>
    });

    return <Form>
        <Card>
            <Badge>{project}/{store}</Badge>
            <CreateStory show={showCreateStory} setShow={setShowCreateStory}
                project={project} store={store} folders={folders} />
            <CreateFolder show={showCreateFolder} setShow={setShowCreateFolder}
                project={project} store={store} folders={folders} />

            <ButtonToolbar className="d-flex">
                <Button variant="primary" onClick={_ => setShowCreateStory(true)}>
                    <VscNewFile /> Create Story
                </Button>
                <Button 
                    variant="primary" onClick={_ => setShowCreateFolder(true)}>
                    <VscNewFolder /> Create Folder
                </Button>
                <Button 
                    variant="primary" onClick={fetchItems}>
                    <BiRefresh />
                </Button>
            </ButtonToolbar>
            <ListGroup variant="flush">
                {storyList}
            </ListGroup>
        </Card>
    </Form>
}

export default Store;