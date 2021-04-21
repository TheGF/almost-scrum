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

let monitorInterval = null

function Federation(props) {
    const { project } = useContext(UserContext)
    const { isOpen, onOpen, onClose } = useDisclosure()
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
                setUpdates(stat['new'] + stat['update'] + stat['conflict'])
                toast({
                    title: `Update from ${log.header.user}@${log.header.hostname}`,
                    description: `${stat['new']} new files, ${stat['update']} updates and ${stat['conflict']} conflicts` +
                        '; click on federation button to synchronize',
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

    function exportSince(time) {
        let d = new Date()

        switch (time) {
            case 'today': d.setDate(d.getDate() - 1); break
            case 'week': d.setDate(d.getDate() - 7); break
            case 'month': d.setMonth(d.getMonth() - 1); break
            case 'all': d.setYear(0); break
            default: d = null; break
        }

        Server.postFedExport(project, d)
            .then(files => {
                Server.postFedPush(project)
                    .then(_ => {
                        if (files) {
                            toast({
                                title: `Successful Export`,
                                description: files.join(','),
                                status: "success",
                                duration: 9000,
                                isClosable: true,
                            })
                        }
                    })
            })
    }

    function onFedStatus(status) {
        if (!status || !status.exchanges) {
            setStrength(0)
            return
        }

        const connectedExchanges = Object.entries(status.exchanges)
                                            .filter(([n, c]) => c)
                                            .map(([n, c]) => n)
        setStrength(connectedExchanges.length)
        notifyChanges(status.inbox || [])
    }

    function closeModal() {
        exportSince()
        onClose()
    }

    function check() {
        Server.postFedPull(project)
            .then(_ => Server.getFedStatus(project).then(onFedStatus))
    }

    function startMonitoring() {
        check()
        setInterval(check, 5 * 60000)
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
                            <Tab isDisabled={strenght == 0}><T>updates</T></Tab>
                            <Tab><T>exchanges</T></Tab>
                            <Tab><T>join</T></Tab>
                            <Tab isDisabled={strenght == 0}><T>invite</T></Tab>
                        </TabList>

                        <TabPanels>
                            <TabPanel>
                                <Updates key={isOpen} onClose={closeModal} exportSince={exportSince} />
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