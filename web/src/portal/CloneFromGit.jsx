import {
    Button, ButtonGroup, Divider, FormControl, FormLabel,
    Input, List, ListIcon, ListItem, VStack,
} from "@chakra-ui/react";
import { Checkbox, CheckboxGroup } from "@chakra-ui/react"
import { React, useEffect, useState } from "react";
import { FaRegCheckCircle, FaRegCircle } from 'react-icons/all';
import Server from '../server';


function CloneFromGit(props) {
    const { onClose, onCreate } = props
    const [templates, setTemplates] = useState([])
    const [selectedTemplates, setSelectedTemplates] = useState([])
    const [url, setUrl] = useState(null)
    const [injectProject, setInjectProject] = useState(null)
    const [cloneInProgress, setCloneInProgress] = useState(false)

    function loadTemplates() {
        Server.getTemplatesList()
            .then(setTemplates)
    }
    useEffect(loadTemplates, [])

    function cloneFromGit() {
        const m = /http.*\/([^\/]*?)(.git)?$/g.exec(url)
        if (m && m.length > 1) {
            const name = m[1]
            setCloneInProgress(true)
            Server.cloneFromGit(url, injectProject, selectedTemplates)
                .then(_ => onCreate(name))
                .then(_ => setCloneInProgress(false))
        }
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
        <FormControl id="url" isRequired>
            <FormLabel>Git URL</FormLabel>
            <Input type="text" onChange={e => setUrl(e.target.value)} value={url} />
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
            <Button colorScheme="blue" onClick={cloneFromGit} isLoading={cloneInProgress}>
                Clone
            </Button>
            <Button onClick={onClose}>Cancel</Button>
        </ButtonGroup>
    </VStack>
}
export default CloneFromGit