import {
    Button,
    FormLabel, Input, Modal,




    ModalBody,
    ModalCloseButton, ModalContent,

    ModalFooter, ModalHeader, ModalOverlay, Select, useToast
} from "@chakra-ui/react";
import React, { useContext, useEffect, useState } from 'react';
import Server from '../server';
import UserContext from '../UserContext';



function MakePost(props) {
    const { project, info } = useContext(UserContext)
    const { value, onClose } = props
    const [boards, setBoards] = useState([])
    const [board, setBoard] = useState(null)
    const [type, setType] = useState(null)
    const [title, setTitle] = useState("")
    const toast = useToast()


    function listBoards() {
        Server.listBoards(project)
            .then(setBoards)
    }
    useEffect(listBoards, [])

    function makePost() {
        Server.postChatAttachmentAction(project, value.id, 0, 'make_post', board, title, type)
            .then(_ => {
                onClose()
                toast({
                    title: 'Created',
                    description: 'The post has been successfully created',
                    status: "success",
                    isClosable: true,
                })
            })
    }

    const boardsUI = boards.map(b => <option key={b} value={b}>{b}</option>)
    const typesUI = info.models.map(m => <option key={m.name} value={m.name}>{m.name}</option>)

    return <Modal isOpen={value != null} onClose={onClose} size="xl">
        <ModalOverlay />
        <ModalContent>
            <ModalHeader>Create a Post</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
                <FormLabel>Board</FormLabel>
                <Select placeholder="Choose" value={board} onChange={e => setBoard(e.target.value)}>{boardsUI}</Select>
                <FormLabel mt={2}>Type</FormLabel>
                <Select placeholder="Choose" value={type} onChange={e => setType(e.target.value)}>{typesUI}</Select>
                <FormLabel mt={2}>Title</FormLabel>
                <Input placeholder="Title" value={title} onChange={e => setTitle(e.target.value)} />
            </ModalBody>
            <ModalFooter>
                <Button colorScheme="blue" mr={3} disabled={!board || !type || !title.length}
                    onClick={makePost}>
                    Create
                </Button>
            </ModalFooter>
        </ModalContent>
    </Modal>
}

export default MakePost