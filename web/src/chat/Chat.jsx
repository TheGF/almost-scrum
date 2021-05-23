import {
    Box,
    Button, Drawer, DrawerBody, DrawerContent, DrawerFooter, DrawerOverlay,

    HStack, IconButton, useDisclosure, VStack
} from "@chakra-ui/react";
import React, { useContext, useState } from 'react';
import { AiOutlineSend, BiUserVoice, BsFillMicFill, GoTextSize, GrAttachment } from 'react-icons/all';
import MarkdownEditor from "../core/MarkdownEditor";
import Server from "../server";
import UserContext from "../UserContext";
import Messages from './Messages';
import RecordAudio from './RecordAudio';
import AttachFiles from './AttachFiles';


function Chat(props) {
    const { project } = useContext(UserContext)
    const { isOpen, onOpen, onClose } = useDisclosure()
    const [text, setText] = useState('')
    const [reload, setReload] = useState(false)
    const btnRef = React.useRef()
    const [inputType, setInputType] = useState('text')
    const [files, setFiles] = useState([])
    let audioRef = null

    function reset() {
        setInputType('text')
        setFiles([])
        setText('')
        setReload(!reload)
    }

    function sendMessage() {
        Server.postChatMsg(project, text, files)
            .then(reset)
    }

    function addAudio(file) {
        setFiles([file])
    }

    function updateInputType(inputType) {
        setInputType(inputType)
        setFiles([])
        setText('')
    }


    function getInput() {
        const toolbarItems = ['heading', 'bold', 'italic', 'quote', 'ul', 'ol', 'image', 'codeblock']

        const inputs = {
            text: <MarkdownEditor key={reload} toolbarItems={toolbarItems}
                height="150" value={text} hideModeSwitch={true}
                disablePreview={true} onChange={setText} imageFolder="/.inline-images" />,
            audio: <HStack width="100%">
                <RecordAudio value={files[0]} onChange={addAudio} />
            </HStack>,
            attach: <AttachFiles value={{ text: text, files: files }}
                onChange={v => { setFiles(v.files); setText(v.text) }} />
        }

        const textButton = inputType == 'text' ?
            <IconButton colorScheme="blue" size="sm" onClick={sendMessage} disabled={!text.length}>
                <AiOutlineSend />
            </IconButton> :
            <IconButton size="sm" onClick={_ => updateInputType('text')}><GoTextSize /></IconButton> //
        const audioButton = inputType == 'audio' ?
            <IconButton colorScheme="blue" size="sm" onClick={sendMessage} disabled={files == null}>
                <AiOutlineSend />
            </IconButton> :
            <IconButton size="sm" onClick={_ => updateInputType('audio')}><BsFillMicFill /></IconButton>
        const attachButton = inputType == 'attach' ?
            <IconButton colorScheme="blue" size="sm" onClick={sendMessage} disabled={files == null}>
                <AiOutlineSend />
            </IconButton> :
            <IconButton size="sm" onClick={_ => updateInputType('attach')}><GrAttachment /></IconButton>


        return <HStack w="100%">
            <Box w="100%">
                {inputs[inputType]}
            </Box>
            <VStack spacing={1} >
                {textButton}
                {audioButton}
                {attachButton}
            </VStack>
        </HStack>
    }



    return <>
        <IconButton ref={btnRef} onClick={onOpen}><BiUserVoice /></IconButton>
        <Drawer
            isOpen={isOpen}
            placement="right"
            onClose={onClose}
            finalFocusRef={btnRef}
            size="md"
        >
            <DrawerOverlay />
            <DrawerContent>

                <DrawerBody>
                    <VStack mb={5} width="100%">
                        {getInput()}
                    </VStack>
                    <Messages key={reload} />
                </DrawerBody>

                <DrawerFooter>
                    <Button variant="outline" mr={3} onClick={onClose}>
                        Close
                    </Button>
                </DrawerFooter>
            </DrawerContent>
        </Drawer>
    </>

}

export default Chat