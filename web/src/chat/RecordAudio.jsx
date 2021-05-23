
import React, { useState, Component } from 'react';
import { Button, HStack, IconButton, Progress, Spacer } from "@chakra-ui/react";
import MicRecorder from 'mic-recorder-to-mp3';
import { BiTrash, BsFillPauseFill, BsFillPlayFill, BsFillStopFill, MdFiberManualRecord } from 'react-icons/all';
import Server from "../server";
import { render } from '@testing-library/react';


const recorder = new MicRecorder();

function RecordAudio(props) {
    const { value, onChange } = props
    const [file, setFile] = useState(value)

    function stop() {
        recorder.stop().getMp3()
        .then(([buffer, blob]) => {
            const file = new File(buffer, 'me-at-thevoice.mp3', {
                type: blob.type,
                lastModified: Date.now()
            })
            setFile(file)
            onChange && onChange(file)
        })
    }

    if (file == null) {
        recorder.start()
        return <HStack width="100%">
            <Progress size="xs" width="100%" isIndeterminate /> :
                <IconButton onClick={stop}><BsFillStopFill /></IconButton>
        </HStack>
    }

    const src = URL.createObjectURL(file)
    return <HStack width="100%">
        <audio controls style={{ height: '24px', width: '100%' }}>
            <source src={src} type="audio/mp3" />
        </audio>
    </HStack>
}

export default RecordAudio