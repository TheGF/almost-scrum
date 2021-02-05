import {
    Button, ButtonGroup, Divider, FormControl, FormLabel,
    Input, List, ListIcon, ListItem, VStack,
} from "@chakra-ui/react";
import { Checkbox, CheckboxGroup } from "@chakra-ui/react"
import { React, useEffect, useState } from "react";
import { FaRegCheckCircle, FaRegCircle } from 'react-icons/all';
import Server from '../server';

function ImportProject(props) {
    const { onClose, onCreate } = props
    const [templates, setTemplates] = useState([])
    const [selectedTemplates, setSelectedTemplates] = useState([])
    const [folder, setFolder] = useState(null)
    const [injectProject, setInjectProject] = useState(null)
    const [importInProgress, setImportInProgress] = useState(false)

    function loadTemplates() {
        Server.getTemplatesList()
            .then(setTemplates)
    }
    useEffect(loadTemplates, [])

    function importProject() {
        setImportInProgress(true)
        Server.importProject(folder, injectProject, selectedTemplates)
            .then(_ => onCreate(folder))
            .then(_ => setImportInProgress(false))
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

    const templatesList = templates.map(t => <ListItem onClick={_ => switchTemplate(t)}>
        {selectedTemplates.includes(t) ?
            <ListIcon as={FaRegCheckCircle} color="green.500" /> :
            <ListIcon as={FaRegCircle} color="green.500" />
        }
        {t}
    </ListItem>)


    return <VStack spacing={5}>
        <FormControl id="folder" isRequired>
            <FormLabel>Folder</FormLabel>
            <Input type="text" onChange={e => setFolder(e.target.value)} value={folder} />
        </FormControl>
        <FormControl id="inject">
            <Checkbox isChecked={injectProject} onChange={_ => setInjectProject(!injectProject)} >
                Add a new Project to the folder
        </Checkbox>
        </FormControl>
        {injectProject ? <FormControl id="templates">
            <FormLabel>Templates</FormLabel>
            <List spacing={3}>
                {templatesList}
            </List>
        </FormControl> : null
        }

        <br />
        <ButtonGroup>
            <Button colorScheme="blue" onClick={importProject} isLoading={importInProgress}>
                Import
            </Button>
            <Button onClick={onClose}>Cancel</Button>
        </ButtonGroup>
    </VStack>
}
export default ImportProject