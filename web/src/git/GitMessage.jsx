import {
    Accordion,
    AccordionButton,
    AccordionIcon, AccordionItem,
    AccordionPanel, Box, Button, Center, Spacer,
    HStack, List, ListIcon, ListItem, StackDivider, Textarea, VStack, Flex, Input, Text,
} from "@chakra-ui/react";
import { React, useContext, useEffect, useState } from "react";
import { BiCheckCircle, BiCircle } from "react-icons/bi";
import Server from '../server';
import UserContext from '../UserContext';
import TaskComment from './TaskComment';

function GitMessage(props) {
    const { project, info } = useContext(UserContext)
    const [infos, setInfos] = useState([]);
    const { gitMessage, setGitMessage } = props;

    function fetchTasks() {
        Server.listTasks(project, '~', `@${info.systemUser}`, 0, 20)
            .then(setInfos)
    }
    useEffect(fetchTasks, [])

    function newTask() {
        Server.createTask(project, info.currentBoard, 'rename me')
            .then(n => {
                const newTask = {
                    board: info.currentBoard,
                    name: n,
                    isNew: true,
                }
                setInfos([newTask, ...infos])
            })
    }

    function changeHeader(evt) {
        const header = evt && evt.target && evt.target.value
        gitMessage.header = header
        setGitMessage({ ...gitMessage })
    }

    const comments = infos ? infos.map(info => <TaskComment key={info.name}
        gitMessage={gitMessage} setGitMessage={setGitMessage} {...info} />) : null
    return <VStack >
        <Input placeholder="Message header" value={gitMessage.header} isRequired
            onChange={changeHeader}></Input>
        <HStack>
            <Spacer />
            <Text>Select the task where you progressed during this commit and
            describe the progress</Text>
            <Spacer />
            <Button onClick={newTask}>New Task</Button>
        </HStack>
        <Flex overflow="auto" h="20em" w="100%">
            <Accordion defaultIndex={[]} allowMultiple w="100%">
                {comments}
            </Accordion>
        </Flex>
    </VStack>
}
export default GitMessage;