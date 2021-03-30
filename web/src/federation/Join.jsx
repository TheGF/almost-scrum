import {
    Box, Button, ButtonGroup, FormControl, FormLabel, Input, Link, Modal,
    ModalBody, ModalCloseButton, ModalContent, ModalHeader, Spacer,
    Text, Textarea, useDisclosure, useToast, VStack
} from "@chakra-ui/react";
import { React, useEffect, useState } from "react";
import { BiWorld } from "react-icons/bi";
import T from "../core/T";
import Server from "../server";

function Join(props) {
    const [token, setToken] = useState(null)
    const [key, setKey] = useState(null)
    const toast = useToast()

    function lookForInvite() {
        const queryParams = new URLSearchParams(window.location.search);
        const invite = queryParams.get('invite');
        const token = queryParams.get('token');

        if (invite == "" && token) {
            setToken(token)
            onOpen()
        }
    }
    useEffect(lookForInvite, [])

    function claim() {
        Server.postFedClaim({
            token: token,
            key: key,
        }).then(_ => {
            toast({
                title: `Claim Success`,
                description: 'The invite has been successfully claimed',
                status: "success",
                duration: 9000,
                isClosable: true,
            })

        })
    }

    return <VStack>
        <FormControl isRequired>
            <FormLabel><T>decryption key</T></FormLabel>
            <Input value={key} onChange={e => setKey(e.target.value)} />
        </FormControl>
        <FormControl isRequired>
            <FormLabel>Token</FormLabel>
            <Textarea rows={12} value={token} onChange={e => setToken(e.target.value)} />
        </FormControl>
        <ButtonGroup>
            <Button colorScheme="blue" onClick={claim}>Claim</Button>
            <Button onClick={onClose}>Close</Button>
        </ButtonGroup>
    </VStack>
}

export default Join