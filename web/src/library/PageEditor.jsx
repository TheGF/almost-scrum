
import {
    Button, Modal, ModalBody, ModalCloseButton, ModalContent,
    ModalFooter, ModalHeader
} from "@chakra-ui/react";
import { React, useContext, useEffect, useState, useRef } from "react";
import "react-mde/lib/styles/css/react-mde-all.css";
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
                // .then(_ => setInterval(_=>{ 
                //         const h = document.getElementById('PageEditor') &&
                //         document.getElementById('PageEditor').clientHeight - 100
                //     console.log('height', h)
                //     if (h) setHeight(h - 100);
                //     }, 1000))
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
        <ModalContent m={0} >
            <ModalHeader>
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
                    height={height}
                />
            </ModalBody>
            <ModalFooter>
            </ModalFooter>
        </ModalContent>
    </Modal>
}
export default PageEditor