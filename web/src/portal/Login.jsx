import {
    Button, FormControl, FormHelperText, FormLabel, Input, InputGroup,
    InputRightElement, Modal, ModalBody, ModalContent,
    ModalHeader, Spacer
} from "@chakra-ui/react";
import { React, useState } from "react";
import T from "../core/T";
import Server from '../server';

function Login(props) {
    const { systemUser, isOpen, onAuthenticated } = props
    const [showPassword, setShowPassword] = useState(false)
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [error, setError] = useState(null)

    function authenticate(e) {
        e && e.preventDefault()
        Server.authenticate(username, password)
            .then(r => onAuthenticated(r, password == 'changeme'))
            .catch(r => setError(`Invalid Credentials: ${r.message}`));
    }

    return <Modal isOpen={isOpen} >
        <ModalContent>
            <ModalHeader>Welcome to Almost Scrum</ModalHeader>
            <ModalBody>
                <form onSubmit={authenticate}>
                    <FormControl required>
                        <FormLabel><T>username</T></FormLabel>
                        <Input id="username" type="text" onChange={e => e && e.target
                            && setUsername(e.target.value)} />
                    </FormControl>
                    <FormControl required>
                        <FormLabel><T>password</T></FormLabel>
                        <InputGroup size="md">
                            <Input
                                id="password"
                                pr="4.5rem"
                                type={showPassword ? "text" : "password"}
                                onChange={e => e && e.target && setPassword(e.target.value)}
                            />
                            <InputRightElement width="4.5rem">
                                <Button h="1.75rem" size="sm"
                                    onClick={e => setShowPassword(!showPassword)}>
                                    {showPassword ? "Hide" : "Show"}
                                </Button>
                            </InputRightElement>
                        </InputGroup>
                    </FormControl>
                    <FormHelperText>{error ?
                        error :
                        `Use ${systemUser}:changeme the first time`
                    }</FormHelperText>
                    <Spacer minH="2em" />
                    <Button type="submit" colorScheme="blue" isFullWidth>
                        Login
                    </Button>
                </form>
            </ModalBody>
        </ModalContent>
    </Modal >
}

export default Login