import {
    AlertDialog, AlertDialogBody, AlertDialogContent,
    AlertDialogFooter, AlertDialogHeader, AlertDialogOverlay,
    Button
} from "@chakra-ui/react";
import { React, useRef } from "react";



function ConfirmChangeOwner({ owner, candidateOwner , setCandidateOwner, onConfirm }) {
    const onClose = () => setCandidateOwner(null)
    const cancelRef = useRef()

    return (
        <AlertDialog
            isOpen={!!candidateOwner}
            leastDestructiveRef={cancelRef}
            onClose={onClose}
        >
            <AlertDialogOverlay>
                <AlertDialogContent>
                    <AlertDialogHeader fontSize="lg" fontWeight="bold">
                        Confirm New Owner
              </AlertDialogHeader>

                    <AlertDialogBody>
                        Please ask <b>{owner}</b> to perform the change.<br/><br/>
                        You can force the change but this may corrupt the task. Be careful!
              </AlertDialogBody>

                    <AlertDialogFooter>
                        <Button ref={cancelRef} onClick={onClose}>
                            Cancel
                </Button>
                        <Button colorScheme="red" onClick={onConfirm} ml={3}>
                            Force
                </Button>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialogOverlay>
        </AlertDialog>
    )
}

export default ConfirmChangeOwner