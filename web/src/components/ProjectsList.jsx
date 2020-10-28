import React, { useEffect, useState } from 'react';
import ListGroup from 'react-bootstrap/ListGroup';
import ToggleButton from 'react-bootstrap/ToggleButton';
import ButtonGroup from 'react-bootstrap/ButtonGroup';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';
import InputGroup from 'react-bootstrap/InputGroup';
import FormControl from 'react-bootstrap/FormControl';
import Card from 'react-bootstrap/Card';
import axios from 'axios';
import { getConfig, loginWhenUnauthorized } from './axiosUtils';
import { MdFavoriteBorder, MdFavorite, MdSearch } from 'react-icons/md';


function ProjectsList(props) {
    const [activeOption, setActiveOption] = useState('All');
    const [projectsList, setProjectsList] = useState([]);
    const [searchFilter, setSearchFilter] = useState('');

    function searchAndActiveFilter(p) {
        return (!searchFilter || p.includes(searchFilter)) &&
            (activeOption == 'All' || localStorage[`fav-${p}`]);
    }

    function FavButtonGroup(props) {
        const options = ['All', 'Favorites'];
        const buttons = options.map(op => <ToggleButton key={op} type="radio"
            variant="primary" checked={activeOption == op}
            onChange={() => setActiveOption(op)}>
            {op}
        </ToggleButton >);

        return <ButtonToolbar>
            <ButtonGroup toggle className="mb-2">
                {buttons}
            </ButtonGroup>
        </ButtonToolbar>
    }

    function SearchInput(props) {
        return <InputGroup key="SearchGroup" className="mb-2">
            <InputGroup.Prepend>
                <InputGroup.Text><MdSearch /></InputGroup.Text>
            </InputGroup.Prepend>
            <FormControl key="SearchForm" placeholder="Search" onChange={e => setSearchFilter(e.target.value)}
                value={searchFilter} />
        </InputGroup>
    }

    function FavIcon(props) {
        const property = `fav-${props.project}`;
        const [flag, setFlag] = useState(!!localStorage[property])
        return flag
            ? <MdFavorite onClick={() => localStorage.removeItem(property) || setFlag(false)} />
            : <MdFavoriteBorder onClick={() => localStorage.setItem(property, 'true') || setFlag(true)} />;
    }

    function ProjectsListGroup(props) {
        const groups = projectsList.filter(searchAndActiveFilter).map(
            p => <ListGroup.Item key={p} action onClick={() => props.onSelect(p)}>
                <FavIcon project={p} />&nbsp;{p}
            </ListGroup.Item>
        );
        return <ListGroup>{groups}</ListGroup>
    }

    function fetch() {
        axios.get('/api/v1/projects', getConfig())
            .then(r => setProjectsList(r.data))
            .catch(loginWhenUnauthorized);
    }
    useEffect(fetch, []);

    return <Card bg="dark">
        <Card.Body>
            <FavButtonGroup />
            <SearchInput />
            <ProjectsListGroup onSelect={props.onSelect}/>
        </Card.Body>
    </Card>


}

export default ProjectsList;