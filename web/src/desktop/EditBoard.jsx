
import {
    Button, ButtonGroup, FormControl, FormLabel, HStack, Input, List, ListItem, Modal,
    ModalBody, ModalCloseButton, ModalContent, ModalFooter, ModalHeader, Spacer, Switch, VStack
} from "@chakra-ui/react";
import { React, useContext, useEffect, useState } from "react";
import "react-mde/lib/styles/css/react-mde-all.css";
import Server from '../server';
import UserContext from '../UserContext';


function EditBoard(props) {
    const { project, info } = useContext(UserContext);
    const { name, onClose, onListBoards } = props
    const [isEmpty, setIsEmpty] = useState(false)
    const [newName, setNewName] = useState(null)
    const [updating, setUpdating] = useState(false)
    const [deleting, setDeleting] = useState(false)

    const boardTypes = info.config.boardTypes && info.config.boardTypes[name] || []
    const [types, setTypes] = useState(boardTypes)

    function checkIsEmpty() {
        if (name) {
            setNewName(name)
            Server.listTasks(project, name, '', 0, 1)
                .then(items => setIsEmpty(items.length == 0))

        }
    }
    useEffect(checkIsEmpty, [name])

    function update() {
        if (name != newName) {
            setUpdating(true)
            Server.renameBoard(project, name, newName)
                .then(onListBoards)
                .then(onClose)
                .then(_ => setUpdating(false))
        }
        if (types != boardTypes) {
            if ( info.config.boardTypes == null) {
                info.config.boardTypes = {}
            }
            info.config.boardTypes[name] = types
            Server.setProjectInfo(project, info)
        }
    }

    function switchType(type) {
        if (types.includes(type)) {
            setTypes(types.filter(t => t != type))
        } else {
            setTypes([type, ...types])
        }
    }

    function del() {
        setDeleting(true)
        Server.deleteBoard(project, name)
            .then(onListBoards)
            .then(onClose)
            .then(_ => setDeleting(false))
    }

    const typesList = info.models.map(m => <ListItem key={m.name}>
        <HStack w="80%">
            <label>{m.name}</label>
            <Spacer/>
            <Switch isChecked={types.includes(m.name)} onChange={_=>switchType(m.name)}/></HStack>
    </ListItem>)

    return <Modal isOpen={name != null} onClose={onClose} >
        <ModalContent >
            <ModalHeader>
                Edit Board
            </ModalHeader>
            <ModalCloseButton />
            <ModalBody id="EditBoard">
                <VStack spacing={5}>
                    <FormControl>
                        <FormLabel>Name</FormLabel>
                        <Input value={newName} onChange={e => setNewName(e.target.value)} />
                    </FormControl>
                    <FormControl>
                        <FormLabel>Types</FormLabel>
                        <List spacing={1}>
                            {typesList}
                        </List>
                    </FormControl>
                    <ButtonGroup>
                        <Button variant="primary" onClick={update}
                            loading={updating} disabled={name == newName && boardTypes == types}>
                            Update
                        </Button>
                        <Button variant="primary" onClick={del}
                            loading={deleting} disabled={!isEmpty}>
                            Delete
                        </Button>
                        <Button onClick={onClose}>Close</Button>
                    </ButtonGroup>
                </VStack>
            </ModalBody>
            <ModalFooter>
            </ModalFooter>
        </ModalContent>
    </Modal>
}
export default EditBoard