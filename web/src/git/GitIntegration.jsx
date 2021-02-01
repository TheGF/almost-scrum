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
import GitPull from './GitPull';


function GitIntegration({ isOpen, onClose }) {
    const [gitMessage, setGitMessage] = useState({ header: '', body: {} })
    const [stagedFiles, setStagedFiles] = useState([])

    function reset() {
        setStagedFiles([])
        setGitMessage({ header: '', body: {} })
    }

    function close() {
        reset()
        onClose()
    }

    return <Modal isOpen={isOpen} onClose={onClose} size="full" top
        scrollBehavior="inside" >
        <ModalOverlay />
        <ModalContent top>
            <ModalHeader>Git Integration</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
                <Tabs defaultIndex={1}>
                    <TabList>
                        <Tab><T>pull</T></Tab>
                        <Tab><T>files</T></Tab>
                        <Tab><T>message</T></Tab>
                        <Tab isDisabled={!stagedFiles || !gitMessage.header}>
                            <T>commit</T>
                        </Tab>
                        <Tab><T>push</T></Tab>
                        <Tab><T>settings</T></Tab>
                    </TabList>

                    <TabPanels>
                        <TabPanel>
                            <GitPull />
                        </TabPanel>
                        <TabPanel>
                            <GitFiles stagedFiles={stagedFiles} setStagedFiles={setStagedFiles} />
                        </TabPanel>
                        <TabPanel>
                            <GitMessage gitMessage={gitMessage} setGitMessage={setGitMessage} />
                        </TabPanel>
                        <TabPanel>
                            <GitCommit stagedFiles={stagedFiles} gitMessage={gitMessage} onCommit={reset} />
                        </TabPanel>
                        <TabPanel>
                            <GitPush />
                        </TabPanel>
                        <TabPanel>
                            <GitSettings />
                        </TabPanel>
                    </TabPanels>
                </Tabs>
            </ModalBody>

            <ModalFooter>
                <Button colorScheme="blue" mr={3} onClick={close}>
                    Close
            </Button>
            </ModalFooter>
        </ModalContent>
    </Modal>
}
export default GitIntegration