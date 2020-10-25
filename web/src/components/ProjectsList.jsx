import React, { useEffect, useState } from 'react';
import ListGroup from 'react-bootstrap/ListGroup';
import axios from 'axios';


function ProjectsList() {
    const [projectsList, setProjectsList] = useState([]);

    function fetch() {
        axios.get('/api/v1/projects', {headers: globalThis.axiosHeaders})
            .then(r => setProjectsList(r.data))
    }
    useEffect(fetch, []);

    const groups = projectsList.map( p => <ListGroup.Item>{p}</ListGroup.Item> );
    return <ListGroup>{groups}</ListGroup>
}

export default ProjectsList;