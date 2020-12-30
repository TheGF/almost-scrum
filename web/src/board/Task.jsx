import { Box, Editable, EditableInput, EditablePreview, HStack, Select, Spacer, Tab, TabList, TabPanel, TabPanels, Tabs } from "@chakra-ui/react";
import { React, useContext, useEffect, useState } from "react";
import T from "../core/T";
import Server from '../server';
import UserContext from '../UserContext';
import TaskEditor from './TaskEditor';
import TaskViewer from './TaskViewer';
import Properties from './Properties';


function Task(props) {
    const { project } = useContext(UserContext);
    const { info } = props;
    const { board, name, modTime } = info;
    const [task, setTask2] = useState(null)

    const owner = task && task.features && task.features['owner']
    const points = task && task.features && task.features['points']
    const [id, title] = name && name.split(/\.(.+)/) || ['', 'Something went wrong']

    function getTask() {
        Server.getTask(project, board, name)
            .then(setTask2)
    }
    useEffect(getTask, [])

    function setTask(task) {
        Server.setTask(project, board, name, task)
    }

    return task && <Box p={1} w="100%" minHeight="100px" borderWidth="3px" overflow="hidden">
        <HStack spacing={3}>
            <label>{id}.</label>
            <Editable defaultValue={title} borderWidth="1px" minW="300px"
                borderColor="blue">
                <EditablePreview />
                <EditableInput />
            </Editable>
            <Spacer />
            <Select placeholder={board || ''} w="10em"></Select>
            <Select placeholder={owner || ''} w="10em"></Select>
        </HStack>
        <HStack spacing={3}>
        <Tabs w="100%">
            <TabList>
                <Tab><T>view</T></Tab>
                <Tab><T>edit</T></Tab>
                <Tab><T>properties</T></Tab>
                <Tab><T>progress</T></Tab>
                <Tab><T>attachments</T></Tab>
            </TabList>

            <TabPanels>
                <TabPanel padding={0}>
                    <TaskViewer task={task} setTask={setTask}  />
                </TabPanel>
                <TabPanel padding={0}>
                    <TaskEditor task={task} setTask={setTask} />
                </TabPanel>
                <TabPanel>
                    <Properties task={task} setTask={setTask} />
                </TabPanel>
                <TabPanel>
                </TabPanel>
                <TabPanel>
                </TabPanel>
            </TabPanels>
        </Tabs>
        </HStack>
    </Box>

}

export default Task;