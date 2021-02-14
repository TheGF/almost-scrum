import {
    Badge,
    Button, Heading, HStack, IconButton, Input, List, ListItem, Spacer, Text, useToast,
    VStack
} from "@chakra-ui/react";
import { React, useEffect, useState } from "react";
import { BsTrash, CgUserlane, BiReset } from 'react-icons/all';
import T from "../core/T";
import Server from '../server';

function LocalUsers(props) {
    const [users, setUsers] = useState([])
    const [login, setLogin] = useState(null)
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState(null)
    const [showChangePassword, setShowChangePassword] = useState(false)
    const [isChanging, setIsChanging] = useState(false)
    const toast = useToast()

    function loadUserInfo() {
        Server.getLocalUsers()
            .then(setUsers)
        Server.getLoginUser()
            .then(setLogin)
    }
    useEffect(loadUserInfo, []);

    function setCredentials(username, password, title, description) {
        setIsChanging(true)
        Server.postLocalUserCredentials(username, password)
            .then(_ => {
                loadUserInfo()
                toast({
                    title: title,
                    description: description,
                    status: "success",
                    duration: 9000,
                    isClosable: true,
                })
                setUsername('')
                setIsChanging(false)
            })
    }

    function changePassword() {
        setCredentials(login, password, `Password changed`)
    }

    function resetUser(username) {
        setCredentials(username, 'changeme', `User '${username}' reset`,
            <span>Password is <i>changeme</i></span>)
    }

    function delUser(username) {
        setCredentials(username, '', `User '${username}' removed`)
    }

    function addUser() {
        setCredentials(username, 'changeme', `User '${username}' added`,
            <span>Password is <i>changeme</i></span>)
    }

    const userList = users.filter(u => u.includes(username)).sort().map(u => {
        const isMe = u == login

        return <ListItem borderWidth={1} padding={1} >
            <HStack>
                <Text>
                    {isMe ? <HStack><b>{u}</b><CgUserlane /></HStack> : u}
                </Text>
                <Spacer />

                {!isMe ? <>
                    <IconButton icon={<BiReset />} title="Reset User's Password"
                        onClick={_ => resetUser(u)} size="sm" />
                    <IconButton icon={<BsTrash />} title="Remove User"
                        onClick={_ => delUser(u)} size="sm" />
                </> :
                    null
                }

            </HStack>
        </ListItem>
    })

    function validUser() {
        return username != null && username != '' && !users.includes(username)
            && username.charAt(0).match(/[a-zA-Z]/)
    }

    const changePasswordUI = showChangePassword ?
        <HStack>
            <Input type="password" w="100%" value={password} isLoading={isChanging}
                onChange={e => setPassword(e.target.value)} />
            <Button minW="10em" onClick={changePassword} disabled={!password}>
                <T>change</T>
            </Button>
        </HStack> :
        <Button onClick={_ => setShowChangePassword(true)}>
            Change My Password
        </Button>

    return <VStack align="left">
        {changePasswordUI}
        <Spacer minH="1em" />
        <Heading size="sm">Add/Filter users</Heading>
        <HStack>
            <Input type="text" w="100%" value={username} description="User login"
                onChange={e => setUsername(e.target.value)} />
            <Button minW="10em" onClick={addUser} disabled={!validUser()} isLoading={isChanging}>
                <T>Add User</T>
            </Button>
        </HStack>
        <List spacing={1} maxH="20em" style={{ overflow: 'auto' }}>
            {userList}
        </List>
    </VStack >
}
export default LocalUsers