import React, { useState } from 'react';
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';
import Refinement from './Refinement';
import Library from '../library/Library';


function ProjectViews(props) {
    const { project } = props;
    const [refresh, setRefresh] = useState(false);

    function reload() {
        setRefresh(!refresh);
    }
    const views = [Refinement, Library]; 
        //','Planning','Daily','Review', 'All', 'Users', 'Tools'];
    const tabs = views.map(V =>
        <Tab key={V.name} eventKey={V.name} title={V.name} style={{ background: 'whitesmoke' }}>
            <V project={project} reload={reload} />
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