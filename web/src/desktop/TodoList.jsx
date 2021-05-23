
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
import { React, useEffect, useState, useContext } from "react";
import { GrFormAdd } from "react-icons/gr";
import Server from "../server";
import UserContext from '../UserContext';
import ReactDatePicker from 'react-datepicker';



function TodoList(props) {
    const {value, onChange} = props


    function TodoItem(props) {
        const item = props.item
        const [done, setDone] = useState(item.done)

        const time = item.eta && new Date(item.eta).toLocaleTimeString()

        return <MenuItem onClick={_=>onChange(value, item)}  isFocusable={false}>
            <HStack w="100%">
                <Text isTruncated maxW="20em" title={`${item.action} - ${time}`}>{item.action}</Text>
                <Spacer />
                <Switch isChecked={done} onChange={_=>setDone(!done)} onClick={e=>e.stopPropagation()}/>
            </HStack>
        </MenuItem>
    }
    

    function addItem(e) {
        const item = { id: value.length, action: 'Todo', eta: new Date(), done: false }
        onChange([...value, item], item)
    }

    function changeTodo() {
    }

    const items = value && value.map(item => <TodoItem key={item.action} item={item}/>)

    return <>
        <MenuGroup title={
            <HStack>
                <span>TODO</span>
                <Spacer />
                <IconButton size="sm" onClick={addItem}><GrFormAdd /></IconButton>
            </HStack>} >

            {items}
        </MenuGroup>
    </>

}

export default TodoList