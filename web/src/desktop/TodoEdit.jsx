import { MenuGroup } from "@chakra-ui/menu";
import { Button, FormControl, FormLabel, HStack, IconButton, Input, MenuItem, Spacer, Switch, Tab, Table, Td, Text, Tr, useDisclosure } from '@chakra-ui/react';
import {
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
} from "@chakra-ui/react"
import { Portal } from "@chakra-ui/react"
import { React, useEffect, useState, useContext } from "react";
import { GrFormAdd } from "react-icons/gr";
import Server from "../server";
import UserContext from '../UserContext';
import ReactDatePicker from 'react-datepicker';


function TodoEdit(props) {
    const {value, onConfirm, onDelete} = props
    const [item, setItem] = useState(value)

    useEffect(() => setItem(value), [value])

    function updateAction(e) {
        setItem({...item, action: e.target.value})
    }
    function updateEta(v) {
        setItem({...item, eta: v})
    }

    function onClose() {
        onConfirm(item)
    }

    function onDelete_() {
        onDelete(item)
    }

    const time = item && new Date(item.eta)

    return <Modal isOpen={value != null} onClose={onClose} size="2xl">
        <ModalOverlay />
        <ModalContent>
            <ModalHeader>Edit Action</ModalHeader>
            <ModalBody>
                <form>
                    <FormControl>
                        <FormLabel>Action</FormLabel>
                        <Input size="sm" value={item && item.action} autofocus onChange={updateAction} />
                    </FormControl>
                    <FormControl mt={4}>
                        <FormLabel>ETA</FormLabel>
                        <ReactDatePicker
                            selected={time}
                            onChange={updateEta}
                            showTimeSelect
                            showTimeSelectOnly
                            timeIntervals={30}
                            timeCaption="Time"
                            dateFormat="h:mm aa"
                        />
                    </FormControl>

                </form>
            </ModalBody>
            <ModalFooter>
                <Button colorScheme="blue" mr={3} onClick={onClose}>
                    Save
                </Button>
                <Button onClick={onDelete_}>Delete</Button>
            </ModalFooter>
        </ModalContent>
    </Modal>
}

export default TodoEdit