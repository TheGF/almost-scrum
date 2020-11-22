import React, { useState } from 'react';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Card from 'react-bootstrap/esm/Card';
import Row from 'react-bootstrap/Row';
import Store from './Store';
import Story from '../story/Story';

function Refinement(props) {
    const { project, reload } = props;
    const [story, setStory] = useState(null);
    const [store, setStore] = useState(null);

    function handleSelect(store, story) {
        setStore(store);
        setStory(story);
    }

    return <Card>
        <Card.Body style={{ background: 'whitesmoke', padding: '0.2em 0em' }}>
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


export default Refinement;