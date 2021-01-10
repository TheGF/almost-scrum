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

function GitMessage(props) {
    const { project, info } = useContext(UserContext)
    const [infos, setInfos] = useState([]);

    function fetchTasks() {
        Server.listTasks(project, '~', `@${info.system_user}`, 0, 20)
            .then(setInfos)
    }
    useEffect(fetchTasks, [])

    function TaskComment(props) {
        const { name, board } = props
        const [task, setTask] = useState(null)

        function fetchFromServer() {
            Server.getTask(project, board, name)
                .then(setTask)
        }

        const progress = task && task.parts && task.parts.map(part => <ListItem>
            <ListIcon as={part.done ? BiCheckCircle : BiCircle} color="green.500" />
            {part.description}
        </ListItem>)

        return <AccordionItem>
            <AccordionButton _expanded={{ bg: "tomato", color: "white" }}
                onClick={fetchFromServer}
            >
                <Box flex="1" textAlign="left">
                    {name}
                </Box>
                <AccordionIcon />
            </AccordionButton>
            <AccordionPanel pb={4}>
                <HStack
                    divider={<StackDivider borderColor="gray.200" />}
                >
                    <Flex w="50%">
                        <Textarea
                            placeholder="Enter the comment for this task" />
                    </Flex>
                    <Flex w="50%">
                        <VStack>
                            <label><b>Progress</b></label>
                            <List spacing={3}>
                                {progress}
                            </List>
                        </VStack>
                    </Flex>
                </HStack>
            </AccordionPanel>
        </AccordionItem>
    }

    const accordions = infos ? infos.map(info => <TaskComment {...info} />) : null

    return <VStack >
        <Input placeholder="Message header" ></Input>
        <HStack w="100%">
            <Spacer />
            <Text>Select the task where you progressed during this commit and
            describe the progress</Text>
            <Spacer />
            <Button>New Task</Button>
        </HStack>
        <Accordion defaultIndex={[]} allowMultiple w="100%">
            {accordions}
        </Accordion>
    </VStack>
}
export default GitMessage;