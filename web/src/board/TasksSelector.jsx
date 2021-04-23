import { React, useContext, useEffect, useRef, useState } from "react";
import {
    Button, Input, Modal, ModalBody, ModalCloseButton, ModalContent,
    ModalFooter, ModalHeader, ModalOverlay, Spacer, useDisclosure, HStack
} from "@chakra-ui/react";
import { BiEdit } from 'react-icons/bi';
import ReactTags from 'react-tag-autocomplete';
import Server from "../server";
import UserContext from '../UserContext';


function TasksSelector(props) {
    const { project } = useContext(UserContext);
    const { value, onChange, readOnly, maxSize } = props
    const { isOpen, onOpen, onClose } = useDisclosure()
    const [suggestions, setSuggestions] = useState([])
    const [tags, setTags] = useState(getTags())

    const ids = tags.map(t => t.id)
    const label = ids.join(',')
    
    function getTags() {
        const tasks = value && value.split(',') || []
        return tasks.map(t => ({ 
            id: t.split(/\.(.+)/)[0],
            name: t
        }))
    }

    function onClick() {
        Server.listTasks(project, '~', '')
            .then(items => setSuggestions(items.map(v => ({
                id: v.id,
                name: v.name,
            }))))
        onOpen()
    }
    function onLeave() {
        updateValue(tags)
        onClose()
    }

    function suggestionsTransform(query, suggestions) {
        const names = tags.map(t=>t.name)
        return suggestions.filter(s => s.name.toLowerCase().includes(query.toLowerCase()))
                        .filter(s=> !names.includes(s.name))
                        .sort((a,b) => b.id-a.id)
    }

    function updateValue(tags) {
        const value = tags.map(t => t.name).join(',')
        onChange(value)
    }

    function addTag(tag) {
        if (!readOnly && tags.length <= maxSize) {
            setTags([...tags, tag])
        }
    }

    function deleteTag(i) {
        if (!readOnly) {
            const tags_ = tags.slice(0)
            tags_.splice(i, 1)
            setTags(tags_)
        }
    }

    return <HStack>
        <Input value={label} readOnly size="xs" onClick={onClick} />
        <Spacer />
        <Button onClick={onClick} size="xs"><BiEdit /></Button>
        <Modal isOpen={isOpen} onClose={onLeave} size="3xl">
            <ModalOverlay />
            <ModalContent>
                <ModalHeader>Edit Tasks</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                <ReactTags
                        placeholderText="type and add tasks"
                        tags={tags}
                        minQueryLength={1}
                        maxSuggestionsLength={16}
                        suggestions={suggestions}
                        suggestionsTransform={suggestionsTransform}
                        onDelete={deleteTag}
                        onAddition={addTag}
                    />
                </ModalBody>
                <ModalFooter>
                    <Button colorScheme="blue" mr={3} onClick={onLeave}>
                        Close
                         </Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
    </HStack>
}

export default TasksSelector