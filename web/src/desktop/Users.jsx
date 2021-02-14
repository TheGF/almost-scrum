import {
    Button, HStack, List, ListItem, VStack, IconButton, Input, Spacer, FormLabel
} from "@chakra-ui/react";
import {
    Accordion,
    AccordionItem,
    AccordionButton,
    AccordionPanel,
    AccordionIcon,
    Box
} from "@chakra-ui/react"
import { React, useState, useEffect, useContext } from "react";
import T from "../core/T";
import { BsTrash, CgUserlane } from 'react-icons/all';
import Server from '../server';
import UserContext from '../UserContext';

function Users(props) {
    const { project, info } = useContext(UserContext)
    const [users, setUsers] = useState([])
    const [username, setUsername] = useState(null)

    function loadUserList() {
        Server.listUsers(project)
            .then(setUsers)
    }
    useEffect(loadUserList, []);

    function setUser(name, userInfo) {
        Server.setUser(project, name, userInfo)
    }

    function delUser(name) {
        Server.delUser(project, name)
            .then(loadUserList)
    }

    function addUser() {
        Server.setUser(project, username, {})
            .then(_ => {
                loadUserList()
                setUsername('')
            })
    }

    const userList = users.sort().map(u => {
        const isMe = u == info.loginUser

        return <AccordionItem>
            <AccordionButton >
                <Box flex="1" textAlign="left" >
                    {isMe ? <HStack><b>{u}</b><CgUserlane/></HStack> : u}
                </Box>
                <AccordionIcon />
            </AccordionButton>
            <AccordionPanel pb={4}>
                <VStack>
                    <HStack align="top">
                        <VStack w="100%">
                            <Input placeholder="Name" type="text" />
                            <Input placeholder="E-mail" type="email" />
                            <Input placeholder="Office Location" type="text" />
                        </VStack>
                        <Box minW="120px" h="120px" borderWidth={1}>
                        </Box>
                    </HStack>
                    <Button w="100%" onClick={_ => delUser(u)} disabled={isMe}>
                        Remove
                    </Button>
                </VStack>
            </AccordionPanel>
        </AccordionItem>
    })

    return <VStack align="left">
        <HStack>
            <Input type="text" w="100%" value={username}
                onChange={e => setUsername(e.target.value)} />
            <Button minW="10em" onClick={addUser}>
                <T>Add User</T>
            </Button>
        </HStack>
        <Accordion maxH="20em" allowToggle style={{ overflow: 'auto' }}>
            {userList}
        </Accordion>
        <a href="https://www.vecteezy.com/free-vector/avatar-icon">Avatar Icon Vectors by Vecteezy</a>

    </VStack >
}
export default Users