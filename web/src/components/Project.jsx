import React, { useEffect, useState } from 'react';
import ListGroup from 'react-bootstrap/ListGroup';
import axios from 'axios';
import { getConfig, loginWhenUnauthorized } from './axiosUtils';
import Card from 'react-bootstrap/Card';
import Tabs from 'react-bootstrap/Tabs';
import Tab from 'react-bootstrap/Tab';
import ProjectViews from './ProjectViews';

function NoProjectSelected(props) {
    return <div style={{ marginTop: '50px' }}>
        <h1>Hi {props.username}. Welcome!</h1>
        Click on a project on the left
    </div>
}


function Project(props) {
    const { project } = props;
    const [projectData, setProjectData] = useState(null);
    const [forbidden, setForbidden] = useState(null);

    

    function fetch() {
        project && axios.get(`/api/v1/projects/${project}`, getConfig())
            .then(r => setProjectData(r.data))
            .catch(r => {
                loginWhenUnauthorized(r);
                if (r.response.status == 403) {
                    setForbidden(r.response.data);
                }
            })
    }
    useEffect(fetch, [project]);


    const content = !project ?
        <NoProjectSelected username={props.username} />
        : <ProjectViews project={project}/>
    return content;
}

export default Project;