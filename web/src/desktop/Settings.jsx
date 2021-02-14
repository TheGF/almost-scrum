
import {
    Modal, ModalBody, ModalCloseButton, ModalContent, ModalFooter, ModalHeader, 
    ModalOverlay, Button
} from "@chakra-ui/react";
import { React } from "react";
import Users from './Users';

function Settings(props) {
    const { isOpen, onClose } = props

    return (
        <>
            <Modal isOpen={isOpen} onClose={onClose} size="lg">
                <ModalOverlay />
                <ModalContent>
                    <ModalHeader>Users</ModalHeader>
                    <ModalCloseButton />
                    <ModalBody>
                        <Users />
                    </ModalBody>

                    <ModalFooter>
                        <Button colorScheme="blue" mr={3} onClick={onClose}>
                            Close
                  </Button>
                    </ModalFooter>
                </ModalContent>
            </Modal>
        </>
    )

}

export default Settings