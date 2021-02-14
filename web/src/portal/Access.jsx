import {
    Button, HStack, Modal,
    ModalBody,
    ModalCloseButton, ModalContent,
    ModalFooter, ModalHeader, Text
} from "@chakra-ui/react";
import { Tabs, TabList, TabPanels, Tab, TabPanel } from "@chakra-ui/react"
import { React } from "react";
import { GiCometSpark } from "react-icons/gi";
import LocalUsers from './LocalUsers';

function Access(props) {
    const { isOpen, onClose } = props

    return <Modal isOpen={isOpen} size="lg" >
        <ModalContent>
            <ModalHeader>Access</ModalHeader>
            <ModalBody>
                <Tabs isLazy>
                    <TabList>
                        <Tab>Local Users</Tab>
                        <Tab>LDAP</Tab>
                        <Tab>OAuth2</Tab>
                    </TabList>

                    <TabPanels>
                        <TabPanel>
                            <LocalUsers />
                        </TabPanel>
                        <TabPanel>
                            <HStack>
                                <GiCometSpark />
                                <Text>Coming Soon</Text>
                            </HStack>
                        </TabPanel>
                        <TabPanel>
                            <HStack>
                                <GiCometSpark />
                                <Text>Coming Soon</Text>
                            </HStack>
                        </TabPanel>
                    </TabPanels>
                </Tabs>
            </ModalBody>
            <ModalFooter>
                <Button colorScheme="blue" mr={3} onClick={onClose}>Close</Button>
            </ModalFooter>
        </ModalContent>
    </Modal >
}

export default Access