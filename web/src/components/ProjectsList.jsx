import React, { useEffect, useState } from 'react';
import Button from 'react-bootstrap/Button';
import ButtonGroup from 'react-bootstrap/ButtonGroup';
import ButtonToolbar from 'react-bootstrap/ButtonToolbar';
import Card from 'react-bootstrap/Card';
import FormControl from 'react-bootstrap/FormControl';
import InputGroup from 'react-bootstrap/InputGroup';
import ListGroup from 'react-bootstrap/ListGroup';
import ToggleButton from 'react-bootstrap/ToggleButton';
import { CgPushLeft, CgPushRight } from 'react-icons/cg';
import { MdFavorite, MdFavoriteBorder, MdSearch } from 'react-icons/md';
import Server from './Server';

function ProjectsList(props) {
    const [activeOption, setActiveOption] = useState('All');
    const [projectsList, setProjectsList] = useState([]);
    const [searchFilter, setSearchFilter] = useState('');
    const [collapsed, setCollapsed] = useState(false);

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
                <Button onClick={_ => setCollapsed(true)} className="float-right">
                    <CgPushLeft />
                </Button>
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
            ? <MdFavorite onClick={e => {
                localStorage.removeItem(property);
                setFlag(false);
                e.stopPropagation();
            }} />
            : <MdFavoriteBorder onClick={e => {
                localStorage.setItem(property, 'true');
                setFlag(true);
                e.stopPropagation();
            }} />;
    }

    function ProjectsListGroup(props) {
        const groups = projectsList.filter(searchAndActiveFilter).map(
            p => <ListGroup.Item key={p} action onClick={() => {
                props.onSelect(p);
                setCollapsed(true);
            }}>
                <FavIcon project={p} />&nbsp;{p}
            </ListGroup.Item>
        );
        return <ListGroup>{groups}</ListGroup>
    }

    function fetch() {
        Server.getProjectsList().then(setProjectsList);
    }
    useEffect(fetch, []);

    const content = <>
        <FavButtonGroup />
        <SearchInput />
        <ProjectsListGroup onSelect={props.onSelect} />
    </>
    const expand = <Button style={{ writingMode: 'tb-rl', minHeight: '60%' }}
        onClick={_ => setCollapsed(false)}>
        <CgPushRight /> Show Project List
        </Button>

    return <Card bg="dark">
        <Card.Body>
            {collapsed ? expand : content}
        </Card.Body>
    </Card>


}

export default ProjectsList;