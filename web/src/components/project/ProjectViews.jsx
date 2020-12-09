import React, { useState } from 'react';
import Tab from 'react-bootstrap/Tab';
import Tabs from 'react-bootstrap/Tabs';
import Refinement from './Refinement';
import Planning from './Planning';
import Library from '../library/Library';
import Team from './Team';


function ProjectViews(props) {
    const { project } = props;
    const [refresh, setRefresh] = useState(false);
    const views = [Refinement, Planning, Library, Team]; 
    const [key, setKey] = useState(localStorage.getItem('currentView') || views[0].name);

    function selectView(key) {
        localStorage.setItem('currentView', key);
        setKey(key);
    }


    function reload() {
        setRefresh(!refresh);
    }
        //','Planning','Daily','Review', 'All', 'Users', 'Tools'];
    const tabs = views.map(V =>
        <Tab key={V.name} eventKey={V.name} title={V.name} style={{ background: 'whitesmoke' }}>
            <V project={project} reload={reload} />
        </Tab>);


    return (
        <Tabs key={refresh}
            activeKey={key}
            onSelect={k => selectView(k)}
            variant="pills"
        >
            {tabs}
        </Tabs>
    );
}

export default ProjectViews;