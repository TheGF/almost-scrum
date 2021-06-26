import {
    Button, Icon, IconButton, Modal, ModalBody, ModalCloseButton, ModalContent,
    ModalFooter, ModalHeader, ModalOverlay, Switch, Tab, TabList, TabPanel, TabPanels, Tabs, useDisclosure
} from "@chakra-ui/react";
import { useToast } from "@chakra-ui/react"
import { React, useEffect, useContext, useState } from "react";
import { BiTransfer } from "react-icons/bi";
import T from "../core/T";
import Server from "../server";
import Updates from "./Updates";
import UserContext from '../UserContext';
import { MdSignalCellular0Bar, MdSignalCellular1Bar, MdSignalCellular2Bar, MdSignalCellular3Bar, MdSignalCellular4Bar } from "react-icons/md";
import Exchanges from "./Exchanges";
import Invite from "./Invite";
import Join from './Join';
import Parked from './Parked';

let monitorInterval = null

function Federation(props) {
    const { project, fedState, setFedState } = useContext(UserContext)
    const { isOpen, onOpen, onClose } = useDisclosure()
    const [strenght, setStrength] = useState(0)
    const [updates, setUpdates] = useState(0)
    const toast = useToast()
    const [lastCheck, setLastCheck] = useState(
        localStorage.getItem(`stg-${project}-lastCheck`) || JSON.stringify(new Date(0))
    )

    function notifyChanges(state, time) {
        localStorage.setItem(`stg-${project}-lastCheck`, JSON.stringify(time))
        setLastCheck(JSON.stringify(time))
        const parked = state.parked && state.parked.length || 0
        const updates = state.updates && state.updates.length || 0
        const sent = state.sent && state.sent.length || 0

        if (updates == 0 && sent == 0) {
            return
        }

        const description = `${updates} updated, ${sent} sent and ${parked} parked files;` +
        ' click on federation for information and conflict resolution'

        toast({
            title: `Federation updates`,
            description: description,
            status: "success",
            duration: 9000,
            isClosable: true,
        })
    }

    function closeModal() {
        onClose()
    }

    function check() {
        Server.getFedExchangeList(project).then(
            ls => {
                ls && setStrength(Object.keys(ls).filter(k => !ls[k]).length)
                Server.getFedState(project, lastCheck)
                    .then(s => notifyChanges(s, new Date()))
            })       
    }

    function startMonitoring() {
        check()
        setInterval(check, 30 * 1000)
    }
    useEffect(startMonitoring, [])

    let signalBar = 0
    switch (strenght) {
        case 0: signalBar = <MdSignalCellular0Bar />; break
        case 1: signalBar = <MdSignalCellular1Bar />; break
        case 2: signalBar = <MdSignalCellular2Bar />; break
        case 3: signalBar = <MdSignalCellular3Bar />; break
        default: signalBar = <MdSignalCellular4Bar />
    }

    return <>
        <Button onClick={onOpen}>
            {signalBar}
            <BiTransfer />
            {updates ? updates : null}
        </Button>
        <Modal isOpen={isOpen} onClose={closeModal} size="6xl" top
            scrollBehavior="inside" isLazy >
            <ModalOverlay />
            <ModalContent top maxW="70%">
                <ModalHeader>Federation</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    <Tabs isLazy>
                        <TabList>
                            <Tab isDisabled={strenght == 0}><T>parked</T></Tab>
                            <Tab isDisabled={strenght == 0}><T>updates</T></Tab>
                            <Tab><T>exchanges</T></Tab>
                            <Tab><T>join</T></Tab>
                            <Tab isDisabled={strenght == 0}><T>invite</T></Tab>
                        </TabList>

                        <TabPanels>
                            <TabPanel>
                                <Parked key={isOpen} onClose={closeModal} />
                            </TabPanel>
                            <TabPanel>
                                <Updates key={isOpen} onClose={closeModal} />
                            </TabPanel>
                            <TabPanel>
                                <Exchanges onClose={closeModal}/>
                            </TabPanel>
                            <TabPanel>
                                <Join onClose={closeModal} project={project} />
                            </TabPanel>
                            <TabPanel>
                                <Invite onClose={closeModal} />
                            </TabPanel>
                        </TabPanels>
                    </Tabs>
                </ModalBody>
            </ModalContent>
        </Modal>
    </>
}
export default Federation