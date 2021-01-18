
import { React, useEffect, useState, useContext } from "react";
import ReactMde from "react-mde";
import ReactMarkdown from "react-markdown";
import "react-mde/lib/styles/css/react-mde-all.css";
import { Box, Button, Center } from "@chakra-ui/react";
import T from "../core/T";
import UserContext from '../UserContext';
import Server from '../server';
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
    const [selectedTab, setSelectedTab] = useState("write");

    function getFromServer() {
        if (page) {
            Server.downloadFromlibrary(project, `${page}/index.md`)
                .then(setContent)
        }
    }
    useEffect(getFromServer, [page])

    const save = async function* (data) {
        const name = `${Date.now()}`
        await Server.uploadFileToLibrary(project, page, new Blob([data]), name)
        yield `/api/v1/projects/${project}/library${page}/${name}`;
        return true;
    };
    function onClose() {
        setPage(null)
        setContent(null)
    }

    function onChange(value) {
        setContent(value)
    }

    function Image(props) {
        return <img {...props} style={{ maxWidth: '20%', maxHeight: '20%'}} />
    }

    return <Modal isOpen={content} onClose={onClose} size="full" >
        <ModalContent>
            <ModalHeader>{page}</ModalHeader>
            <ModalCloseButton />
            <ModalBody>
                <ReactMde
                    value={content}
                    onChange={onChange}
                    selectedTab={selectedTab}
                    onTabChange={setSelectedTab}
                    generateMarkdownPreview={(markdown) =>
                        Promise.resolve(<ReactMarkdown
                            source={markdown}
                            renderers={{ image: Image }}
                        />)
                    }
                    paste={{
                        saveImage: save
                    }}
                />
            </ModalBody>
            <ModalFooter>
                <Button colorScheme="blue" mr={3} onClick={onClose}>
                    Close
            </Button>
            </ModalFooter>
        </ModalContent>
    </Modal>
}
export default PageEditor