import {
    Button, Icon, IconButton, Modal, ModalBody, ModalCloseButton, ModalContent,
    ModalFooter, ModalHeader, ModalOverlay, Switch, Tab, TabList, TabPanel, TabPanels, Tabs, useDisclosure
} from "@chakra-ui/react";
import { useToast } from "@chakra-ui/react"
import { React, useEffect, useContext, useState } from "react";
import { BiTransfer } from "react-icons/bi";
import T from "../core/T";
import Server from "../server";
import Sync from "./Sync";
import UserContext from '../UserContext';
import { MdSignalCellular0Bar, MdSignalCellular1Bar, MdSignalCellular2Bar, MdSignalCellular3Bar, MdSignalCellular4Bar } from "react-icons/md";
import Exchanges from "./Exchanges";
import Share from "./Share";

let monitorInterval = null

function Federation(props) {
    const { project } = useContext(UserContext)
    const { isOpen, onOpen, onClose } = useDisclosure()
//    const [diffs, setDiffs] = useState([])
    const [strenght, setStrength] = useState(0)
    const [updates, setUpdates] = useState(0)
    const toast = useToast()


    function notifyChanges(diffs) {
        const lastChange = new Date(localStorage.getItem('ash-fed-newest-fed-log'))

        for (const log of diffs) {
            const creationTime = new Date(log.creationTime)
            if (creationTime < lastChange) {
                continue
            }

            const stat = { outdated: 0, update: 0, new: 0, conflict: 0 }
            for (const item of Object.values(log.items)) {
                stat[item.match] += 1
            }

            if (stat['new'] || stat['update'] || stat['conflict']) {
                toast({
                    title: `Update from ${log.header.user}@${log.header.hostname}`,
                    description: `${stat['new']} new files, ${stat['update']} updates and ${stat['conflict']} conflicts` +
                        '; click on federation button to update',
                    status: "success",
                    duration: 9000,
                    isClosable: true,
                })
            }
        }
    }
    function getDiffs() {
        Server.getFedDiffs(project, true)
            .then(notifyChanges)
    }

    function sync() {
        Server.postFedSync(project)
    }

    function startMonitoring() {
        Server.getFedStatus(project)
            .then(status => {
                let connectedExchanges = []
                let exchangesNum = 0
                for (const [name, connected] of Object.entries(status.exchanges)) {
                    if (connected) connectedExchanges.push(name)
                    exchangesNum++
                }
                setStrength(connectedExchanges.length)

                if (exchangesNum) {
                    if (connectedExchanges.length) {
                        Server.postFedExport(project)
                        getDiffs()
                        monitorInterval = setInterval(getDiffs, 2 * 60000)
                        toast({
                            title: `Connected`,
                            description: 'Successfully connected to '+connectedExchanges,
                            status: "warning",
                            duration: 9000,
                            isClosable: true,
                        })
                        } else {
                        if (monitorInterval) clearInterval(monitorInterval);
                        setTimeout(startMonitoring, 2 * 60000)
                    }
                } else {
                    toast({
                        title: `No Exchanges`,
                        description: 'No exchanges are configured and federation is not available',
                        status: "warning",
                        duration: 9000,
                        isClosable: true,
                    })
                }
            })

    }
    useEffect(startMonitoring, [])

    let signalBar = 0 
    switch (strenght) {
        case 0: signalBar = <MdSignalCellular0Bar/>; break
        case 1: signalBar = <MdSignalCellular1Bar/>; break
        case 2: signalBar = <MdSignalCellular2Bar/>; break
        case 3: signalBar = <MdSignalCellular3Bar/>; break
        default: signalBar = <MdSignalCellular4Bar/>
    }

    return <>
        <Button onClick={onOpen}>
            {signalBar}
            <BiTransfer />
            {updates ? updates : null}
        </Button>
        <Modal isOpen={isOpen} onClose={onClose} size="6xl" top
            scrollBehavior="inside" >
            <ModalOverlay />
            <ModalContent top maxW="70%">
                <ModalHeader>Federation</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    <Tabs isLazy>
                        <TabList>
                            <Tab><T>sync</T></Tab>
                            <Tab><T>exchanges</T></Tab>
                            <Tab><T>invite</T></Tab>
                            <Tab><T>join</T></Tab>
                        </TabList>

                        <TabPanels>
                            <TabPanel>
                                <Sync key={isOpen} onClose={onClose} />
                            </TabPanel>
                            <TabPanel>
                                <Exchanges onClose={onClose}/>
                            </TabPanel>
                            <TabPanel>
                                <Share onClose={onClose}/>
                            </TabPanel>
                            <TabPanel>
                                <Share onClose={onClose}/>
                            </TabPanel>
                        </TabPanels>
                    </Tabs>
                </ModalBody>
            </ModalContent>
        </Modal>
    </>
}
export default Federation