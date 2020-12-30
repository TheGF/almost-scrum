import { StackDivider, VStack } from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from "react";
import "react-mde/lib/styles/css/react-mde-all.css";
import Server from '../server';
import UserContext from '../UserContext';
import Task from './Task';
import FilterPanel from '../desktop/FilterPanel';


function Board(props) {
    const { project } = useContext(UserContext);
    const { name } = props;
    const [filter, setFilter] = useState('');
    const [start, setStart] = useState(0);
    const [end, setEnd] = useState(5);
    const [infos, setInfos] = useState([]);

    function loadTaskList() {
        Server.listTasks(project, name, filter, start, end)
            .then(setInfos)
    }
    useEffect(loadTaskList, [name]);

    const tasks = infos && infos.map(info => <Task key={info.id} info={info} />);
    return <VStack
        spacing={4}
        align="stretch"
        w="100%"
    >
        <FilterPanel />

        {tasks}
    </VStack>
}

export default Board