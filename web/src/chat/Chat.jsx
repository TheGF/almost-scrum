import {
    Button, Drawer, DrawerBody, DrawerCloseButton, DrawerContent, DrawerFooter, DrawerHeader, DrawerOverlay,
    HStack, IconButton, Input, Progress, Spacer, useDisclosure, VStack
} from "@chakra-ui/react";
import MicRecorder from 'mic-recorder-to-mp3';
import React, { useContext, useState } from 'react';
import { BiTrash, BiUserVoice, BsFillMicFill, BsFillPauseFill, BsFillPlayFill, BsFillStopFill } from 'react-icons/all';
import MarkdownEditor from "../core/MarkdownEditor";
import Server from "../server";
import UserContext from "../UserContext";

const recorder = new MicRecorder();

function Chat(props) {
    const { project } = useContext(UserContext)
    const { isOpen, onOpen, onClose } = useDisclosure()
    const [recording, setRecording] = useState(false)
    const [playing, setPlaying] = useState(false)
    const [mp3, setMp3] = useState(null)
    const [audio, setAudio] = useState(null)
    const btnRef = React.useRef()

    function startRecording() {
        recorder.start().then(() => {
            setRecording(true)
        }).catch((e) => {
            console.error(e);
        });
    }

    function stopRecording() {
        recorder.stop().getMp3()
            .then(([buffer, blob]) => {
                setRecording(false)
                const file = new File(buffer, 'me-at-thevoice.mp3', {
                    type: blob.type,
                    lastModified: Date.now()
                })
                setMp3(file);
                setAudio(new Audio(URL.createObjectURL(file)))
            })
    }

    function playRecording() {
        audio.play()
        audio.onended = _ => setPlaying(false)
        setPlaying(true)
    }

    function pauseRecording() {
        audio.pause()
        setPlaying(false)
    }

    function deleteRecording() {
        audio.pause()
        setAudio(null)
        setPlaying(false)
        setMp3(null)
    }

    function sendRecording() {
        Server.postChatMsg(project, mp3)
            .then(deleteRecording)
    }

    function getInput() {
        if (recording) {
            return <HStack width="100%">
                <Progress size="xs" width="100%" isIndeterminate /> :
                <IconButton onClick={stopRecording}><BsFillStopFill /></IconButton>
            </HStack>
        }

        if (mp3) {
            return <HStack width="100%">
                {playing ? <IconButton onClick={pauseRecording}><BsFillPauseFill /></IconButton> :
                    <IconButton onClick={playRecording}><BsFillPlayFill /></IconButton>}
                <IconButton onClick={deleteRecording}><BiTrash /></IconButton>
                <Spacer />
                <Button onClick={sendRecording}>Send</Button>
            </HStack>
        }
        return <VStack width="100%">
            <MarkdownEditor
                height="200"
                value=""
                disablePreview={true}
                imageFolder="/.inline-images"
            />
            <HStack >
                <Button>Send</Button>
                <IconButton onClick={startRecording}><BsFillMicFill /></IconButton>
            </HStack>
        </VStack>
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
                <DrawerCloseButton />
                <DrawerHeader>Anything to say?</DrawerHeader>

                <DrawerBody>
                    <VStack>
                        {getInput()}
                    </VStack>
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