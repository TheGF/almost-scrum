import {
    Button, FormControl, FormLabel, Input, Modal, ModalBody,
    ModalCloseButton, ModalContent, ModalFooter, ModalHeader, ModalOverlay
} from "@chakra-ui/react";
import { React, useState, useEffect } from "react";
import T from "../core/T";


function AskBoardName(props) {
    const { isOpen, onClose, onCreate, boards } = props
    const [name, setName] = useState('')
    const [isInvalid, setIsInvalid] = useState(checkInvalid(name))

    function checkInvalid(name) {
        return name == '' || (boards & boards.includes(name))
    }
    useEffect(_ => setIsInvalid(checkInvalid(name)), [name])

    return <Modal isOpen={isOpen} onClose={onClose}>
        <ModalOverlay />
        <ModalContent>
            <ModalHeader><T>new board</T></ModalHeader>
            <ModalCloseButton />
            <ModalBody>
                <FormControl id="newBoard">
                    <FormLabel><T>name</T></FormLabel>
                    <Input value={name} onChange={e => setName(e.target.value)}
                        isInvalid={isInvalid}
                        errorBorderColor="crimson" />
                </FormControl>
            </ModalBody>
            <ModalFooter>
                <Button colorScheme="blue" mr={3} disabled={isInvalid}
                    onClick={_ => onCreate(name)}>
                    Confirm
              </Button>
                <Button variant="ghost" onClick={onClose}>Cancel</Button>
            </ModalFooter>
        </ModalContent>
    </Modal>

}

export default AskBoardName