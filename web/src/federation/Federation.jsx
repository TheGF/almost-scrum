import {
    Button, Icon, IconButton, Modal, ModalBody, ModalCloseButton, ModalContent,
    ModalFooter, ModalHeader, ModalOverlay, Tab, TabList, TabPanel, TabPanels, Tabs, useDisclosure
} from "@chakra-ui/react";
import { useToast } from "@chakra-ui/react"
import { React, useEffect, useContext, useState } from "react";
import { BiTransfer } from "react-icons/bi";
import T from "../core/T";
import Server from "../server";
import Sync from "./Sync";
import UserContext from '../UserContext';

let monitorInterval = null

function Federation(props) {
    const { project } = useContext(UserContext)
    const { isOpen, onOpen, onClose } = useDisclosure()
    const [logs, setLogs] = useState([])
    const [updates, setUpdates] = useState(0)
    const toast = useToast()


    function notifyChanges(logs) {
        const lastChange = new Date(localStorage.getItem('ash-fed-newest-fed-log'))

        for (const log of logs) {
            const creationTime = new Date(log.creationTime)
            if ( creationTime < lastChange) {
                continue
            }
            
            const stat = {outdated: 0, update:0, new: 0, conflict: 0}
            for (const item of Object.values(log.items)) {
                stat[item.match] += 1
            }

            if (stat['new'] || stat['update'] || stat['conflict']) {
                toast({
                    title: `Update from ${log.header.user}@${log.header.hostname}`,
                    description: `${stat['new']} new files, ${stat['update']} updates and ${stat['conflict']} conflicts`+
                                 '; click on federation button to update',
                    status: "success",
                    duration: 9000,
                    isClosable: true,
                  })
            }
        }
    }
    function getLogs() {
        Server.getFedLogs(project)
            .then(logs => {
                setLogs(logs)
                notifyChanges(logs)
            })
    }

    function startMonitoring() {
        if (monitorInterval) clearInterval(monitorInterval);
        monitorInterval = setInterval(getLogs, 10*60000)
    }
    useEffect(startMonitoring, [])

    return <>
        <Button  onClick={onOpen}><BiTransfer />{updates ? updates : null}</Button>
        <Modal isOpen={isOpen} onClose={onClose} size="full" top
            scrollBehavior="inside" >
            <ModalOverlay />
            <ModalContent top maxW="70%">
                <ModalHeader>Federation</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    <Tabs isLazy>
                        <TabList>
                            <Tab><T>sync</T></Tab>
                            <Tab><T>hubs</T></Tab>
                            <Tab><T>users</T></Tab>
                            <Tab><T>share</T></Tab>
                        </TabList>

                        <TabPanels>
                            <TabPanel>
                                <Sync key={isOpen} onClose={onClose}/>
                            </TabPanel>
                            <TabPanel>
                            </TabPanel>
                            <TabPanel>
                            </TabPanel>
                            <TabPanel>
                            </TabPanel>
                        </TabPanels>
                    </Tabs>
                </ModalBody>
            </ModalContent>
        </Modal>
    </>
}
export default Federation