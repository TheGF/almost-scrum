import {
    Button, ButtonGroup, Modal, ModalBody, useToast,
    ModalContent, ModalHeader, Text, VStack, Input
} from "@chakra-ui/react";
import { React, useState, useEffect } from "react";
import T from "../core/T";
import Server from '../server';

function ChangePassword(props) {
    const [login, setLogin] = useState(null)
    const { isOpen, onClose } = props
    const [password, setPassword] = useState(null)
    const [isChanging, setIsChanging] = useState(false)
    const toast = useToast()

    function loadUserInfo() {
        Server.getLoginUser()
            .then(setLogin)
    }
    useEffect(loadUserInfo, []);


    function changePassword() {
        setIsChanging(true)
        Server.postLocalUserCredentials(login, password)
            .then(_ => {
                toast({
                    title: 'Password changed',
                    status: "success",
                    duration: 9000,
                    isClosable: true,
                })
                setIsChanging(false)
                onClose()
            })
    }


    return <Modal isOpen={isOpen} size="lg" >
        <ModalContent>
            <ModalHeader>Change Password</ModalHeader>
            <ModalBody>
                <VStack>
                    <Text>Your password is not safe</Text>
                    <Input type="password" w="100%" value={password} isLoading={isChanging}
                        onChange={e => setPassword(e.target.value)} />
                    <ButtonGroup>
                        <Button minW="10em" onClick={changePassword} disabled={!password}>
                            <T>change</T>
                        </Button>
                        <Button minW="10em" onClick={onClose}>
                            <T>close</T>
                        </Button>
                    </ButtonGroup>
                </VStack>
            </ModalBody>
        </ModalContent>
    </Modal >
}

export default ChangePassword