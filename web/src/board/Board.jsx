import { StackDivider, VStack } from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from "react";
import "react-mde/lib/styles/css/react-mde-all.css";
import Server from '../server';
import UserContext from '../UserContext';
import Task from './Task';
import FilterPanel from '../desktop/FilterPanel';


function Board(props) {
    const { project } = useContext(UserContext);
    const { name, boards } = props;
    const [start, setStart] = useState(0);
    const [end, setEnd] = useState(5);
    const [searchKeys, setSearchKeys] = useState([]);
    const [infos, setInfos] = useState([]);
    const [users, setUsers] = useState([]);
    const [compact, setCompact] = useState(false);
    const [tags, setTags] = useState([]);

    function loadTags() {
        Server.getSuggestions(project, '%23', 64)
            .then(setTags)
    }
    useEffect(loadTags, []);

    function loadTaskList() {
        const filter = searchKeys.join(',')
        Server.listTasks(project, name, filter, start, end)
            .then(setInfos)
    }
    useEffect(loadTaskList, [name, searchKeys]);

    function loadUserList() {
        Server.listUsers(project)
            .then(setUsers)
    }
    useEffect(loadUserList, []);

    const tasks = infos && infos.map(info =>
        <Task key={info.id} info={info} compact={compact}
            boards={boards} onBoardChanged={loadTaskList}
            users={users} tags={tags}/>
    );
    return <VStack
        spacing={4}
        align="stretch"
        w="100%"
    >
        <FilterPanel compact={compact} setCompact={setCompact} setSearchKeys={setSearchKeys} />
        {tasks}
    </VStack>

}


export default Board