import {
    Button, ButtonGroup, FormControl, FormLabel, Modal,
    ModalBody, ModalContent, ModalHeader, Select,
    useToast, VStack
} from "@chakra-ui/react";
import { React, useState } from "react";
import T from "../core/T";
import Join from "../federation/Join";
import Server from "../server";

function ClaimInvite(props) {
    const { projects, activeProject, token, setToken } = props
    const [selectedProject, setSelectedProject] = useState(activeProject || null)
    const [confirm, setConfirm] = useState(false)
    const toast = useToast()


    // function join() {
    //     Server.postFedClaim(project, key, token).then(_ => {
    //         toast({
    //             title: `Claim Success`,
    //             description: 'The invite has been successfully claimed',
    //             status: "success",
    //             duration: 9000,
    //             isClosable: true,
    //         })
    //         setInvite(null)
    //     })
    // }

    const projectsUI = projects.filter(p => p).map(p => <option key={p} value={p}>{p}</option>)

    const selectProjectUI = <VStack>
        <FormControl isRequired>
            <FormLabel><T>choose a project</T></FormLabel>
            <Select placeholder="Select option" value={selectedProject} onChange={e => setSelectedProject(e.target.value)}>
                {projectsUI}
            </Select>
        </FormControl>
        <ButtonGroup>
            <Button colorScheme="blue" onClick={_=>setConfirm(true)} isDisabled={!selectedProject}>Confirm</Button>
            <Button onClick={_ => setToken(null)}>Close</Button>
        </ButtonGroup>
    </VStack>


    return <Modal isOpen={token} size="6xl" >
        <ModalContent>
            <ModalHeader>Accept Invite</ModalHeader>
            <ModalBody>
                {confirm ? <Join project={selectedProject} token={token} onClose={_=>setToken(null)}/> : 
                selectProjectUI}
            </ModalBody>
        </ModalContent>
    </Modal >
}

export default ClaimInvite