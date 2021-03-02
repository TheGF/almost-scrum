import {
    Button, ButtonGroup, Link, Modal, ModalBody, ModalCloseButton, ModalContent,
    ModalFooter, ModalHeader, ModalOverlay, Table, Tbody, Td, Th, Thead, Tr
} from '@chakra-ui/react';
import { React, useContext, useEffect, useState } from 'react';
import Utils from '../core/utils';
import Server from '../server';
import UserContext from '../UserContext';

function Versions(props) {
    const { project } = useContext(UserContext);
    const { path, file, onClose } = props
    const [versions, setVersions] = useState([]);
    const isLocalhost = Utils.isLocalhost()

    function getVersions() {
        if (!file) return
        const p = `${path}/${file.name}`
        Server.getVersions(project, p)
            .then(setVersions)
    }
    useEffect(getVersions, [file])

    function deleteFile(file) {
        Server.deleteFromLibrary(project, `${path}/${file.name}`, file.name.endsWith('.pg'))
            .then(listFolder)
    }

    function onFileClick(file) {
        const p = [path, file.name].join('/')
        if (file.dir || file.page) {
            changePath(p)
        } else {
            Server.openFromlibrary(project, p, true)
        }
    }

    function openFile(file) {
        const p = `${path}/${file.name}`
        Server.localOpenFromLibrary(project, p, true)
    }

    function getOpenButton(file) {
        return isLocalhost ?
            <Button onClick={_ => openFile(file)}>Open</Button> :
            null
    }

    const rows = versions && versions.map(file => {
        const idx = file.name.lastIndexOf('.')
        const match = file.name.match(/(.*?)((\d+\.)+\d+)?(\.\w*)?$/)
        if (match.length < 5) {
            return ''
        }

        const version = match[2] || ''
        return <Tr key={file.name}>
            <Td>
                <Link href="#" onClick={_ => onFileClick(file)} >
                    { file.name }
                </Link>
            </Td>
            <Td>{Utils.getFriendlyDate(file.modTime)}</Td>
            <Td>{file.owner}</Td>
            <Td>{version}</Td>
            <Td><span title={file.size}>{Utils.getFriendlySize(file.size)}</span></Td>
            <Td>
                <ButtonGroup size="sm" spacing={2}>
                    {getOpenButton(file)}
                </ButtonGroup>
            </Td>
        </Tr>
    })


    return <Modal isOpen={file != null} onClose={onClose} size="full">
        <ModalOverlay />
        <ModalContent maxW="80%" maxH="80%">
            <ModalHeader>Previous Versions</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
                <Table w="100%">
                    <Thead>
                        <Tr>
                            <Th>Name</Th>
                            <Th>Modified</Th>
                            <Th>Owner</Th>
                            <Th>Version</Th>
                            <Th>Size</Th>
                            <Th>Actions</Th>
                        </Tr>
                    </Thead>
                    <Tbody  >
                        {rows}
                    </Tbody>
                </Table>
            </ModalBody>

            <ModalFooter>
                <Button colorScheme="blue" mr={3} onClick={onClose}>
                    Close
        </Button>
            </ModalFooter>
        </ModalContent>
    </Modal>
}

export default Versions