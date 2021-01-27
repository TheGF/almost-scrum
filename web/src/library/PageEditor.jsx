
import { React, useEffect, useState, useContext } from "react";
import ReactMde from "react-mde";
import ReactMarkdown from "react-markdown";
import "react-mde/lib/styles/css/react-mde-all.css";
import { Box, Button, Center } from "@chakra-ui/react";
import T from "../core/T";
import UserContext from '../UserContext';
import Server from '../server';
import MarkdownEditor from '../core/MarkdownEditor';

import {
    Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
} from "@chakra-ui/react"

function PageEditor(props) {
    const { project } = useContext(UserContext);
    const { page, setPage } = props
    const [content, setContent] = useState(null)
    const [height, setHeight] = useState(1000)
    const [selectedTab, setSelectedTab] = useState("write");

    function getFromServer() {
        if (page) {
            Server.downloadFromlibrary(project, `${page}/index.md`)
                .then(setContent)
                .then(_ => {
                    const h = document.getElementById('PageEditor') && document.getElementById('PageEditor').clientHeight-100
                    if (h) setHeight(h-100);
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
        Server.uploadFileToLibraryLater(project, page, new Blob([value]), 'index.md')
    }

    function Image(props) {
        const token = localStorage.token
        if (token) {
            const src = `${props.src}?token=${token}`
            return <img {...props} style={{ maxWidth: '20%', maxHeight: '20%' }} src={src}/>
        } else {
            return <img {...props} style={{ maxWidth: '20%', maxHeight: '20%' }} />
        }
    }

    return <Modal isOpen={content!=null} onClose={onClose} size="full">
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
                    imageFolder={page}
                    height={height}
                />
            </ModalBody>
            <ModalFooter>
            </ModalFooter>
        </ModalContent>
    </Modal>
}
export default PageEditor