import {
    Accordion, AccordionButton, AccordionIcon, AccordionItem, AccordionPanel,
    Box, Button, ButtonGroup, HStack, Input, MenuItem, Modal,
    ModalBody, ModalCloseButton, ModalContent, ModalFooter, ModalHeader,
    ModalOverlay, VStack, useDisclosure
} from "@chakra-ui/react";
import { React, useContext, useEffect, useState } from "react";
import { CgUserlane } from 'react-icons/all';
import T from "../core/T";
import Server from '../server';
import UserContext from '../UserContext';


function Users(props) {
    const { project, info, reload } = useContext(UserContext)
    const [users, setUsers] = useState([])
    const [username, setUsername] = useState(null)
    const { isOpen, onClose } = props
    const [requireReload, setRequireReload] = useState(false)

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
            .then(_ => {
                loadUserList()
                setRequireReload(true)
            })
    }

    function addUser() {
        Server.setUser(project, username, {})
            .then(_ => {
                loadUserList()
                setUsername('')
                setRequireReload(true)
            })
    }

    function closeModal() {
        onClose()
        if (requireReload) {
            reload()
            setRequireReload(false)
        }
    }

    function UserProfile(props) {
        const { user, isExpanded } = props;
        const isMe = user == info.loginUser
        const [userInfo, setUserInfo] = useState({})

        function getUser() {
            if (isExpanded) {
                Server.getUser(project, user)
                    .then(setUserInfo)
            }
        }
        useEffect(getUser, [isExpanded]);

        return <>
            <AccordionButton >
                <Box flex="1" textAlign="left" >
                    {isMe ? <HStack><b>{user}</b><CgUserlane /></HStack> : user}
                </Box>
                <AccordionIcon />
            </AccordionButton>
            <AccordionPanel pb={4}>
                <VStack>
                    <HStack align="top">
                        <VStack w="100%">
                            <Input placeholder="Name" type="text" value={userInfo.name}
                                onChange={e => setUserInfo({ ...userInfo, name: e.target.value })} />
                            <Input placeholder="E-mail" type="email" value={userInfo.email}
                                onChange={e => setUserInfo({ ...userInfo, email: e.target.value })} />
                            <Input placeholder="Office Location" value={userInfo.location}
                                onChange={e => setUserInfo({ ...userInfo, location: e.target.value })} />
                        </VStack>
                        <Box minW="120px" h="120px" borderWidth={1}>
                        </Box>
                    </HStack>
                    <ButtonGroup w="50%" >
                        <Button onClick={_ => setUser(user, userInfo)}>
                            Update
                        </Button>
                        <Button onClick={_ => delUser(user)} disabled={isMe}>
                            Remove
                        </Button>
                    </ButtonGroup>
                </VStack>
            </AccordionPanel>
        </>
    }

    const userList = users.sort().map(u => <AccordionItem key={u} >
        {({ isExpanded }) => <UserProfile key={u} user={u} isExpanded={isExpanded} />}
    </AccordionItem>)


    const body = <VStack align="left">
        <HStack>
            <Input type="text" w="100%" value={username}
                onChange={e => setUsername(e.target.value)} />
            <Button minW="10em" onClick={addUser}
                isDisabled={!username || !username.length} >
                <T>Add User</T>
            </Button>
        </HStack>
        <Accordion maxH="20em" allowToggle style={{ overflow: 'auto' }}>
            {userList}
        </Accordion>
        <a href="https://www.vecteezy.com/free-vector/avatar-icon">Avatar Icon Vectors by Vecteezy</a>
    </VStack >

    return <Modal isOpen={isOpen} onClose={closeModal} size="lg">
        <ModalOverlay />
        <ModalContent>
            <ModalHeader>Users</ModalHeader>
            <ModalBody>
                {body}
            </ModalBody>

            <ModalFooter>
                <Button colorScheme="blue" mr={3} onClick={closeModal}>
                    Close
                  </Button>
            </ModalFooter>
        </ModalContent>
    </Modal>

}
export default Users