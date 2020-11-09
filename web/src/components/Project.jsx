import React, { useEffect, useState } from 'react';
import ProjectViews from './ProjectViews';
import Server from './Server';

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
        project && Server.getProject(project)
                        .then(setProjectData);
    }
    useEffect(fetch, [project]);


    const content = !project ?
        <NoProjectSelected username={props.username} />
        : <ProjectViews project={project} />
    return content;
}

export default Project;