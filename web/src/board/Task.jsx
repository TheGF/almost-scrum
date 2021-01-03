import { Box, Button, Editable, EditableInput, EditablePreview, HStack, Select, Spacer, Tab, TabList, TabPanel, TabPanels, Tabs } from "@chakra-ui/react";
import { React, useContext, useEffect, useState } from "react";
import T from "../core/T";
import Server from '../server';
import UserContext from '../UserContext';
import TaskEditor from './TaskEditor';
import TaskViewer from './TaskViewer';
import Properties from './Properties';
import Progress from './Progress';
import { MdTouchApp, MdVerticalAlignTop, RiFilterLine } from "react-icons/all";


function Task(props) {
    const { project, info } = useContext(UserContext);
    const { board, name, modTime } = props.info;
    const { compact, boards, users } = props;
    const [task, setTask] = useState(null)
    const [progress, setProgress] = useState('')

    function updateProgress(task) {
        const progress = task && task.parts && task.parts.length ?
            `${Math.round(100 * task.parts.filter(p => p.done).length / task.parts.length)}%` : ''
        setProgress(progress)
    }

    function getTask() {
        Server.getTask(project, board, name)
            .then(setTask)
            .then(updateProgress)
    }
    useEffect(getTask, [])

    function touchTask() {
        Server.touchTask(project, board, name)
    }

    function saveTask(task) {
        Server.setTaskLater(project, board, name, task)
    }

    function renameTask(title) {
        Server.moveTask(project, board, name, board, title)
        props.onBoardChanged && props.onBoardChanged()
    }

    function onBoardChanged(evt) {
        const newBoard = evt && evt.target && evt.target.value;
        newBoard && Server.moveTask(project, board, name, newBoard)
        props.onBoardChanged && props.onBoardChanged()
    }

    if (!task) return null;

    const owner = task && task.properties && task.properties['Owner']
        && task.properties['Owner'].substring(1)
    const [id, title] = name && name.split(/\.(.+)/) || ['', 'Something went wrong']
    const userList = users && users.map(u => <option key={u} value={u}>
        {u}
    </option>)
    const boardList = boards && boards.map(b => <option key={b} value={b}>
        {b}
    </option>)

    const mtime = new Date(modTime).toUTCString();
    const header = task && <HStack spacing={3}>
        <label>{id}.</label>
        <Editable defaultValue={title} borderWidth="1px" minW="300px"
            borderColor="blue" onSubmit={title => renameTask(title)}>
            <EditablePreview />
            <EditableInput />
        </Editable>
        <Spacer />
        <Select value={board} w="10em" onChange={onBoardChanged}>
            {boardList}
        </Select>
        <Select value={owner} w="10em">
            {userList}
        </Select>
        <span>{progress}</span>
        <Button size="sm"><RiFilterLine onClick={touchTask} /></Button>
        <Button size="sm" title={mtime}><MdVerticalAlignTop onClick={touchTask} /></Button>
    </HStack>

    function onChange(index) {
        if (index == 0) {
            setTask({ ...task, description: task.description });
        }
    }


    const body = task && !compact ? <HStack spacing={3}>
        <Tabs w="100%" onChange={onChange}>
            <TabList>
                <Tab><T>view</T></Tab>
                <Tab><T>edit</T></Tab>
                <Tab><T>properties</T></Tab>
                <Tab><T>progress</T></Tab>
                <Tab><T>attachments</T></Tab>
            </TabList>

            <TabPanels>
                <TabPanel padding={0}>
                    <TaskViewer task={task} saveTask={saveTask} />
                </TabPanel>
                <TabPanel padding={0}>
                    <TaskEditor task={task} saveTask={saveTask} />
                </TabPanel>
                <TabPanel>
                    <Properties task={task} saveTask={saveTask} />
                </TabPanel>
                <TabPanel>
                    <Progress task={task} saveTask={task => {
                        saveTask(task);
                        updateProgress(task);
                    }} />
                </TabPanel>
                <TabPanel>
                </TabPanel>
            </TabPanels>
        </Tabs>
    </HStack> : ''

    return task ? <Box p={1} w="100%" borderWidth="3px" overflow="hidden">
        {header}
        {body}
    </Box> : ''

}

export default Task;