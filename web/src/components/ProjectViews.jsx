import React, { useEffect, useState } from 'react';
import Store from './Store';
import Story from './Story';
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/esm/Card';
import Server from './Server';

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


function ProjectViews(props) {
    const { project } = props;
    const [refresh, setRefresh] = useState(false);

    function reload() {
        setRefresh(!refresh);
    }
    const views = [Refinement]; //','Planning','Daily','Review', 'All', 'Users', 'Tools'];
    const tabs = views.map(v =>
        <Tab key={v.name} eventKey={v.name} title={v.name} style={{ background: 'whitesmoke' }}>
            <Refinement project={project} reload={reload} />
        </Tab>);
    const [key, setKey] = useState(views[0].name);

    return (
        <Tabs key={refresh}
            activeKey={key}
            onSelect={(k) => setKey(k)}
            variant="pills"
        >
            {tabs}
        </Tabs>
    );
}

export default ProjectViews;