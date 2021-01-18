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
import GitPush from './GitPush';
import GitSettings from './GitSettings';
import T from "../core/T";

function GitIntegration({ isOpen, onClose }) {
    const [gitStatus, setGitStatus] = useState(null)
    const [gitMessage, setGitMessage] = useState({ header: '', body: {} })
    const [commitDone, setCommitDone] = useState(false)

    function onCommit() {
        setGitStatus(null)
        setGitMessage(null)
        setCommitDone(true)
    }

    return <Modal isOpen={isOpen} onClose={onClose} size="full" top
        scrollBehavior="inside" >
        <ModalOverlay />
        <ModalContent top>
            <ModalHeader>Git Integration</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
                <Tabs>
                    <TabList>
                        <Tab><T>files</T></Tab>
                        <Tab><T>message</T></Tab>
                        <Tab isDisabled={!gitStatus || !gitMessage.header}
                            onCommit={onCommit}><T>commit</T></Tab>
                        <Tab><T>push</T></Tab>
                        <Tab><T>settings</T></Tab>
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
                        <TabPanel>
                            <GitPush/>
                        </TabPanel>
                        <TabPanel>
                            <GitSettings/>
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