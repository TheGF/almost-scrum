import {
    Button, Modal,
    ModalBody,
    ModalCloseButton, ModalContent,
    ModalFooter, ModalHeader
} from "@chakra-ui/react";
import { Tabs, TabList, TabPanels, Tab, TabPanel } from "@chakra-ui/react"
import { React } from "react";
import CreateProject from './CreateProject';
import ImportProject from './ImportProject';
import CloneFromGit from './CloneFromGit';

function AddProject(props) {
    const { isOpen, onCreate, onClose } = props

    function loadUserList() {
        // Server.authenticate(username, password)
        //     .then(r => onToken(r))
        //     .catch(r => setError(`Invalid Credentials: ${r.message}`));
    }

    return <Modal isOpen={isOpen} size="lg" >
        <ModalContent>
            <ModalHeader>Add Project</ModalHeader>
            <ModalBody>
                <Tabs>
                    <TabList>
                        <Tab>Create New</Tab>
                        <Tab>Import Folder</Tab>
                        <Tab>Git Clone</Tab>
                    </TabList>

                    <TabPanels>
                        <TabPanel>
                            <CreateProject onCreate={onCreate} onClose={onClose}/>
                        </TabPanel>
                        <TabPanel>
                            <ImportProject onCreate={onCreate} onClose={onClose}/>
                        </TabPanel>
                        <TabPanel>
                            <CloneFromGit onCreate={onCreate} onClose={onClose}/>
                        </TabPanel>
                    </TabPanels>
                </Tabs>
            </ModalBody>

        </ModalContent>
    </Modal >
}

export default AddProject