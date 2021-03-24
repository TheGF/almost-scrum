
import {
    Button, Modal, ModalBody, ModalCloseButton, ModalContent,
    ModalFooter, ModalHeader
} from "@chakra-ui/react";
import { React, useContext, useEffect, useState } from "react";
import "react-mde/lib/styles/css/react-mde-all.css";
import MarkdownEditor from '../core/MarkdownEditor';
import Server from '../server';
import UserContext from '../UserContext';


function PageEditor(props) {
    const { project } = useContext(UserContext);
    const { page, setPage } = props
    const [content, setContent] = useState(null)
    const [height, setHeight] = useState(1000)
    const [selectedTab, setSelectedTab] = useState("write");

    function getFromServer() {
        if (page) {
            Server.downloadFromlibrary(project, page)
                .then(setContent)
                .then(_ => {
                    const h = document.getElementById('PageEditor') &&
                        document.getElementById('PageEditor').clientHeight - 100
                    if (h) setHeight(h - 100);
                })
        }
    }
    useEffect(getFromServer, [page])

    function onClose() {
        setPage(null)
        setContent(null)
    }

    function onChange(value) {
        setContent(value)
        const lastSlash = page.lastIndexOf("/")
        const folder = page.substring(0, lastSlash)
        const name =  page.substring(lastSlash+1)
        Server.uploadFileToLibraryLater(project, folder, new Blob([value]), name)
    }

    function Image(props) {
        const token = localStorage.token
        if (token) {
            const src = `${props.src}?token=${token}`
            return <img {...props} style={{ maxWidth: '20%', maxHeight: '20%' }} src={src} />
        } else {
            return <img {...props} style={{ maxWidth: '20%', maxHeight: '20%' }} />
        }
    }

    return <Modal isOpen={content != null} onClose={onClose} size="full">
        <ModalContent >
            <ModalHeader>
                <Button colorScheme="blue" mr={3} onClick={onClose}>
                    Close
            </Button>
            I{page}
            </ModalHeader>
            <ModalCloseButton />
            <ModalBody id="PageEditor">
                <MarkdownEditor
                    value={content}
                    onChange={onChange}
                    imageFolder="/.inline-images"
                    height={height}
                />
            </ModalBody>
            <ModalFooter>
            </ModalFooter>
        </ModalContent>
    </Modal>
}
export default PageEditor