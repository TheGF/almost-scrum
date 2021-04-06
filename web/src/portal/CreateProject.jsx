import {
    Button, ButtonGroup, Divider, FormControl, FormLabel,
    Input, List, ListIcon, ListItem, VStack,
} from "@chakra-ui/react";
import { React, useEffect, useState } from "react";
import { FaRegCheckCircle, FaRegCircle } from 'react-icons/all';
import Server from '../server';

function CreateProject(props) {
    const { onClose, onCreate } = props
    const [templates, setTemplates] = useState([])
    const [selectedTemplates, setSelectedTemplates] = useState([])
    const [name, setName] = useState(null)

    function loadTemplates() {
        Server.getTemplatesList()
            .then(setTemplates)
    }
    useEffect(loadTemplates, [])

    function createProject() {
        Server.createProject(name, selectedTemplates)
            .then(_ => onCreate(name))
    }

    function switchTemplate(template) {
        const idx = selectedTemplates.indexOf(template)
        if (idx !== -1) {
            setSelectedTemplates([
                ...selectedTemplates.slice(0, idx),
                ...selectedTemplates.slice(idx + 1),
            ])
        } else {
            setSelectedTemplates([
                ...selectedTemplates,
                template
            ])
        }
    }

    const templatesList = templates.sort().map(t => <ListItem onClick={_ => switchTemplate(t)}>
        {selectedTemplates.includes(t) ?
            <ListIcon as={FaRegCheckCircle} color="green.500" /> :
            <ListIcon as={FaRegCircle} color="green.500" />
        }
        {t}
    </ListItem>)


    return <VStack spacing={5}>
        <FormControl id="name" isRequired>
            <FormLabel>Project Name</FormLabel>
            <Input type="text" onChange={e => setName(e.target.value)} value={name} />
        </FormControl>
        <FormControl id="templates">
            <FormLabel>Templates</FormLabel>
            <List spacing={3}>
                {templatesList}
            </List>
        </FormControl>

        <br />
        <ButtonGroup>
            <Button colorScheme="blue" onClick={createProject}
                isDisabled={name == null || name.length == 0}>
                Create
                </Button>
            <Button onClick={onClose}>Cancel</Button>
        </ButtonGroup>
    </VStack>
}
export default CreateProject