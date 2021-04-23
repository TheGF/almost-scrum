
import {
    Button, Modal, ModalBody, ModalCloseButton, ModalContent,
    ModalFooter, ModalHeader
} from "@chakra-ui/react";
import { React, useContext, useEffect, useState, useRef } from "react";
import MarkdownEditor from '../core/MarkdownEditor';
import Utils from "../core/utils";
import Server from '../server';
import UserContext from '../UserContext';


function PageEditor(props) {
    const { project } = useContext(UserContext);
    const { page, setPage } = props
    const [content, setContent] = useState(null)
    const [height, setHeight] = useState(100)
    const [selectedTab, setSelectedTab] = useState("write");
    const refBody = useRef(null)

    function getFromServer() {
        if (page) {
            Server.downloadFromlibrary(project, page)
                .then(setContent)
        }
    }

    function setH(element) {
        if (!element) return
        const h = element.clientHeight - 100
        if (h != height) {
            setHeight(h)
        }
    }

    useEffect(getFromServer, [page])

    function onClose() {
        Server.sendPendingTasks()
        setPage(null)
        setContent(null)
    }

    function onChange(value) {
        setContent(value)
        const lastSlash = page.lastIndexOf("/")
        const folder = page.substring(0, lastSlash)
        const name =  page.substring(lastSlash+1)

        console.log('Write update', value)

        Server.uploadFileToLibrary(project, folder, new Blob([value]), name)
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

    return <Modal isOpen={content != null} onClose={onClose} size="6xl" >
        <ModalContent m={0} h="90%" background="#eaeaea">
            <ModalHeader >
                <Button colorScheme="blue" mr={3} onClick={onClose}>
                    Close
            </Button>
            I{page}
            </ModalHeader>
            <ModalCloseButton />
            <ModalBody ref={c=>Utils.autoResize(c, 100, setHeight)} >
                <MarkdownEditor
                    value={content}
                    onChange={onChange}
                    imageFolder="/.inline-images"
                    height={height+100}
                />
            </ModalBody>
            <ModalFooter>
            </ModalFooter>
        </ModalContent>
    </Modal>
}
export default PageEditor