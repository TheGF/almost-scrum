import {
    Button, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader, Text, VStack,
} from "@chakra-ui/react";
import { React } from "react";
import { GiStopSign } from "react-icons/gi";



function NoAccess(props) {
    const { data, onExit } = props;

    const show = data != null
    const usersList = data && data.users.join(' ')
    return <>
    <Modal isOpen={show} >
        <ModalContent>
            <ModalHeader>No Access</ModalHeader>
            <ModalBody>
                <VStack>
                    <Text color="red">You have no access to the project</Text>
                    <Text>You can ask access to one of those users: {usersList}</Text>
                </VStack>
            </ModalBody>

            <ModalFooter>
                {onExit ? <Button colorScheme="blue" mr={3} onClick={onExit}>Leave Project</Button> : null}
            </ModalFooter>
        </ModalContent>
    </Modal >
    {show ? <GiStopSign size="30%" color="red"/>: null}
    </>
}

export default NoAccess