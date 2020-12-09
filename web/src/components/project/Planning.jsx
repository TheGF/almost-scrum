import React, { useState } from 'react';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Dropdown from 'react-bootstrap/Dropdown';
import Card from 'react-bootstrap/esm/Card';
import Row from 'react-bootstrap/Row';
import ButtonGroup from 'react-bootstrap/ButtonGroup';
import Story from '../story/Story';
import Store from './Store';

function Planning(props) {
    const { project, reload } = props;
    const [story, setStory] = useState(null);
    const [store, setStore] = useState(null);
    const [spring, setSpring] = useState(null);

    function handleSelect(store, story) {
        setStore(store);
        setStory(story);
    }



    return <Card>
        <Card.Body style={{ background: 'whitesmoke', padding: '0.2em 0em' }}>
            <Row>
                <Col md="9">
                </Col>
                <Col md="3">
                    <Dropdown style={{width: '90%'}}>
                        <Dropdown.Toggle className="btn-block"
                         variant="success" id="dropdown-basic" >
                            Sprint-1
                            </Dropdown.Toggle>

                        <Dropdown.Menu className="btn-block">
                            <Dropdown.Item href="#/action-1">sprint-2</Dropdown.Item>
                            <Dropdown.Item href="#/action-2">sprint-2</Dropdown.Item>
                            <Dropdown.Item href="#/action-3">sprint-3</Dropdown.Item>
                        </Dropdown.Menu>
                    </Dropdown>
                </Col>
            </Row>
            <Container fluid>
                <Row>
                    <Col md="3" style={{ padding: '0em 0.2em' }}>
                        <Store key={story} project={project} store="backlog" story={story}
                            toRight="sandbox"
                            onSelect={s => handleSelect('backlog', s)}
                            reload={reload} />
                    </Col>
                    <Col md="6" style={{ padding: '0em 0.2em' }}>
                        <Story project={project} store={store} story={story}
                            moveTo={store == 'backlog' ? 'sandbox' : 'backlog'}
                            reload={reload} />
                    </Col>
                    <Col md="3" style={{ padding: '0em 0.2em' }}>
                        <Store key={story} project={project} store="sandbox" story={story}
                            toLeft="backlog"
                            onSelect={s => handleSelect('sandbox', s)}
                            reload={reload} />
                    </Col>
                </Row>
            </Container>
        </Card.Body>
    </Card>
}


export default Planning;