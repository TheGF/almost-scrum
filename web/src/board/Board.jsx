import { Button, Center, Spacer, Spinner, StackDivider, VStack } from '@chakra-ui/react';
import { React, useContext, useEffect, useState, useRef } from "react";
import Server from '../server';
import UserContext from '../UserContext';
import Task from './Task';
import FilterPanel from './FilterPanel';
import InfiniteScroll from 'react-infinite-scroll-component';
import './datePicker.css'

function getTop(elem) {
    let top = 0
    while (elem) {
        top += elem.offsetTop
        elem = elem.offsetParent
    }
    return top
}

function Board(props) {
    const { project } = useContext(UserContext)
    const { name, boards } = props
    const [hasMore, setHasMore] = useState(true)
    const [searchKeys, setSearchKeys] = useState([])
    const [infos, setInfos] = useState([])
    const [users, setUsers] = useState([])
    const [compact, setCompact] = useState(false)
    const spacerRef = useRef(null)

    function onNewTask(type) {
        Server.createTask(project, name, 'Click_and_Rename', type)
            .then(_ => loadTaskList())
    }

    function loadMore() {
        const filter = searchKeys.join(',')
        const start = infos.length
        const end = start + 5
        Server.listTasks(project, name, filter, start, end)
            .then(items => {
                items = [...infos, ...items]
                setInfos(items)
                setHasMore(items.length == end)
            })
    }

    function checkSpace() {
        const top = getTop(spacerRef.current)
        if (top < window.innerHeight) {
            loadMore()
        }
    }
    useEffect(_ => { if (hasMore) setTimeout(checkSpace, 100) }, [infos, compact]);

    function loadTaskList() {
        infos.length = 0
        setHasMore(true)
        loadMore()
    }
    useEffect(loadTaskList, [name, searchKeys]);

    function loadUserList() {
        Server.listUsers(project)
            .then(setUsers)
    }
    useEffect(loadUserList, []);

    const tasks = infos && infos.map(info =>
        <Task key={`${info.id}.${info.name}`} info={info} compact={compact}
            boards={boards} onBoardChanged={loadTaskList}
            users={users} searchKeys={searchKeys} />
    );
    return infos ? <VStack
        spacing={4}
        align="stretch"
        w="100%"
    >
        <FilterPanel board={name} compact={compact} setCompact={setCompact} setSearchKeys={setSearchKeys}
            onNewTask={onNewTask} users={users}/>
        <InfiniteScroll
            dataLength={infos.length}
            next={loadMore}
            hasMore={hasMore}
            loader={<Center><Spinner /> Loading...</Center>}
        >
            {tasks}
            <Spacer key={infos.length} ref={spacerRef} />
        </InfiniteScroll>
    </VStack> : null

}


export default Board