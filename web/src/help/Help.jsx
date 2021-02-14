import {
    Button, Modal, ModalBody, ModalCloseButton, ModalContent,
    ModalFooter, ModalHeader, ModalOverlay, Tab, TabList, TabPanel, TabPanels, Tabs
} from "@chakra-ui/react";
import { React } from "react";
import T from "../core/T";
import HelpTab from './HelpTab';
import overview from "./overview.md";
import portal from "./portal.md";
import boards from "./boards.md";
import library from "./library.md"
import hacking from "./hacking.md"


function Help(props) {
    const { isOpen, onClose } = props

    const tabs = [
        [overview, 'overview'],
        [portal, 'portal'],
        [boards, 'boards'],
        [library, 'library'],
        [hacking, 'hacking'],
    ]

    const tabList = tabs.map(t => <Tab key={t[1]}>
        <T>{t[1]}</T>
    </Tab>)

    const tabsPanels = tabs.map(t => <TabPanel key={t[1]}>
        <HelpTab file={t[0]} />
    </TabPanel>)

    return <Modal isOpen={isOpen} onClose={onClose} size="full" top
        scrollBehavior="inside" >
        <ModalOverlay />
        <ModalContent top>
            <ModalBody>
                <Tabs isLazy >
                    <TabList>{tabList}</TabList>
                    <TabPanels>{tabsPanels}</TabPanels>
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
export default Help