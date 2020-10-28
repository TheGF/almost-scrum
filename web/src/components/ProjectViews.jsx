import React, { useEffect, useState } from 'react';
import Store from './Store';
import Story from './Story';
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/esm/Card';


function Refinement(props) {
    const { project } = props;
    const [storyName, setStoryName] = useState(null);
    const [store, setStore] = useState(null);

    return <Container>
        <Card>
            <Row>
                <Col md="3">
                    <Store project={project} store="backlog" onSelect={
                        (s) => setStore('backlog') || setStoryName(s)
                    } />;
            </Col>
                <Col md="6">
                    <Story project={project} store={store} storyName={storyName} />
                </Col>
                <Col md="3">
                    <Store project={project} store="sandbox" onSelect={
                        (s) => setStore('sandbox') || setStoryName(s)
                    } />;
            </Col>
            </Row>
        </Card>
    </Container>
}


function ProjectViews(props) {
    const { project } = props;
    const views = [Refinement]; //','Planning','Daily','Review', 'All', 'Users', 'Tools'];
    const tabs = views.map(v => <Tab key={v.name} eventKey={v.name} title={v.name}>
        <Refinement project={project} />
    </Tab>);
    const [key, setKey] = useState(views[0].name);

    return (
        <Tabs
            activeKey={key}
            onSelect={(k) => setKey(k)}
        >
            {tabs}
        </Tabs>
    );
}

export default ProjectViews;