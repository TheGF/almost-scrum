import {
    Button, Modal,
    ModalBody,
    ModalCloseButton, ModalContent,
    ModalFooter, ModalHeader
} from "@chakra-ui/react";
import { React } from "react";

function Users(props) {
    const { isOpen, onClose } = props

    function loadUserList() {
        // Server.authenticate(username, password)
        //     .then(r => onToken(r))
        //     .catch(r => setError(`Invalid Credentials: ${r.message}`));
    }

    return <Modal isOpen={isOpen} size="lg" >
        <ModalContent>
            <ModalHeader>Users</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
            </ModalBody>
            <ModalFooter>
                <Button colorScheme="blue" mr={3} onClick={onClose}>Close</Button>
            </ModalFooter>
        </ModalContent>
    </Modal >
}

export default Users