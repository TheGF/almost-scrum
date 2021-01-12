import {
    AccordionButton,
    AccordionIcon, AccordionItem,
    AccordionPanel, Box, Text,
    Flex, HStack, Input, List, ListIcon, ListItem,
    Spacer, StackDivider, Textarea, VStack, Select, Editable, EditablePreview, EditableInput
} from "@chakra-ui/react";
import { React, useContext, useState } from "react";
import { BiCheckCircle, BiCircle } from "react-icons/bi";
import Server from '../server';
import UserContext from '../UserContext';

function TaskComment(props) {
    const { project } = useContext(UserContext)
    const { gitMessage, setGitMessage, isNew } = props
    const [name, setName] = useState(props.name)
    const [board, setBoard] = useState(props.board)
    const [task, setTask] = useState(null)

    function fetchFromServer() {
        Server.getTask(project, board, name)
            .then(setTask)
    }

    function changeComment(evt) {
        const comment = evt && evt.target && evt.target.value
        gitMessage.body[name] = comment
        setGitMessage({ ...gitMessage })
    }

    function renameTask(newTitle) {
        const [id, title] = name && name.split(/\.(.+)/) || ['', 'Something went wrong']

        Server.moveTask(project, board, name, board, newTitle)
            .then(_ => setName(`${id}.${newTitle}`))
    }

    function getHeader() {
        if (!isNew) {
            return <HStack>
                <label>{name}</label>
                <Spacer />
                <label>{board}</label>
                <Text w="2em" />
            </HStack>
        }

        const [id, title] = name && name.split(/\.(.+)/) || ['', 'Something went wrong']
        return <HStack>
            <label>{id}</label>
            <Editable defaultValue={title} borderWidth="1px" minW="7em" maxW="300%"
                title="Click to rename" borderColor="blue"
                onSubmit={newTitle => renameTask(newTitle)}>
                <EditablePreview />
                <EditableInput />
            </Editable>
            <Spacer />
            <label>{board}</label>
            <Text w="2em" />
        </HStack>
    }

    const comment = gitMessage.body[name]

    const progress = task && task.parts && task.parts.map(part => <ListItem>
        <ListIcon as={part.done ? BiCheckCircle : BiCircle} color="green.500" />
        {part.description}
    </ListItem>)

    return <AccordionItem>
        <AccordionButton _expanded={{ bg: "tomato", color: "white" }}
            onClick={fetchFromServer}>
            <Box flex="1" textAlign="left">
                {getHeader()}
            </Box>
            <AccordionIcon />
        </AccordionButton>
        <AccordionPanel pb={4}>
            <HStack
                divider={<StackDivider borderColor="gray.200" />}
            >
                <Flex w="50%">
                    <Textarea
                        placeholder="Enter the comment for this task"
                        value={comment} onChange={changeComment} />
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
export default TaskComment