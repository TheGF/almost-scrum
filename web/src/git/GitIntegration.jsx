import { React, useEffect, useState, useContext } from "react";
import {
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    Button,
} from "@chakra-ui/react"
import { Tabs, TabList, TabPanels, Tab, TabPanel } from "@chakra-ui/react"
import GitFiles from './GitFiles'
import GitMessage from './GitMessage';
import GitCommit from "./GitCommit";

function GitIntegration({ isOpen, onClose }) {
    const [gitStatus, setGitStatus] = useState(null)
    const [gitMessage, setGitMessage] = useState({header: '', body: {}})

    return <Modal isOpen={isOpen} onClose={onClose} size="full" top
        scrollBehavior="inside" >
        <ModalOverlay />
        <ModalContent top>
            <ModalHeader>Git Integration</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
                <Tabs>
                    <TabList>
                        <Tab>Commit Files</Tab>
                        <Tab>Commit Message</Tab>
                        <Tab isDisabled={!gitStatus || !gitMessage.header}>Commit</Tab>
                    </TabList>

                    <TabPanels>
                        <TabPanel>
                            <GitFiles gitStatus={gitStatus} setGitStatus={setGitStatus} />
                        </TabPanel>
                        <TabPanel>
                            <GitMessage gitMessage={gitMessage} setGitMessage={setGitMessage} />                        
                        </TabPanel>
                        <TabPanel>
                            <GitCommit gitStatus={gitStatus} gitMessage={gitMessage} />
                        </TabPanel>
                    </TabPanels>
                </Tabs>
            </ModalBody>

            <ModalFooter>
                <Button colorScheme="blue" mr={3} onClick={onClose}>
                    Close
            </Button>
            </ModalFooter>
        </ModalContent>
    </Modal>
}
export default GitIntegration